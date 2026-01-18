package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	MaxConfigHistoryVersions = 10
	ConfigCurrentID          = "main"
)

// ConfigCurrentModel 当前配置（单行）
type ConfigCurrentModel struct {
	ID               string    `gorm:"primaryKey;size:32" json:"id"`
	CurrentVersionID string    `gorm:"size:32;index" json:"currentVersionId"`
	ConfigJSON       string    `gorm:"type:text" json:"-"`
	ConfigHash       string    `gorm:"size:64" json:"configHash"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (ConfigCurrentModel) TableName() string {
	return "config_current"
}

// ConfigHistoryModel 配置版本历史
type ConfigHistoryModel struct {
	StringPKBaseModel
	Version    int64  `gorm:"index" json:"version"`
	ConfigJSON string `gorm:"type:text" json:"-"`
	ConfigHash string `gorm:"size:64" json:"configHash"`
	Source     string `gorm:"size:32" json:"source"` // file/api/rollback/init
	Note       string `gorm:"size:255" json:"note"`
}

func (ConfigHistoryModel) TableName() string {
	return "config_history"
}

// ConfigHashJSON 计算配置 JSON 的 SHA256 哈希
func ConfigHashJSON(configJSON string) string {
	h := sha256.Sum256([]byte(configJSON))
	return hex.EncodeToString(h[:])
}

// GetCurrentConfig 获取当前配置
func GetCurrentConfig() (*ConfigCurrentModel, error) {
	var current ConfigCurrentModel
	err := db.Where("id = ?", ConfigCurrentID).First(&current).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &current, nil
}

// GetConfigHistoryList 获取配置历史列表（按版本号降序）
func GetConfigHistoryList() ([]ConfigHistoryModel, error) {
	var list []ConfigHistoryModel
	err := db.Order("version DESC").Find(&list).Error
	return list, err
}

// GetConfigHistoryByVersion 根据版本号获取配置历史
func GetConfigHistoryByVersion(version int64) (*ConfigHistoryModel, error) {
	var history ConfigHistoryModel
	err := db.Where("version = ?", version).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

// GetConfigHistoryByID 根据 ID 获取配置历史
func GetConfigHistoryByID(id string) (*ConfigHistoryModel, error) {
	var history ConfigHistoryModel
	err := db.Where("id = ?", id).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

// SaveConfigVersion 保存配置版本
// configJSON: 配置的 JSON 字符串
// source: 来源标识 (file/api/rollback/init)
// note: 变更说明
func SaveConfigVersion(configJSON string, source string, note string) error {
	configHash := ConfigHashJSON(configJSON)

	// 检查是否与当前配置相同（避免重复保存）
	current, err := GetCurrentConfig()
	if err != nil {
		return err
	}
	if current != nil && current.ConfigHash == configHash {
		return nil // 配置未变化，无需保存
	}

	// 获取下一个版本号
	var maxVersion int64
	if err := db.Model(&ConfigHistoryModel{}).Select("COALESCE(MAX(version), 0)").Scan(&maxVersion).Error; err != nil {
		return err
	}
	nextVersion := maxVersion + 1

	// 创建历史记录
	history := ConfigHistoryModel{
		Version:    nextVersion,
		ConfigJSON: configJSON,
		ConfigHash: configHash,
		Source:     source,
		Note:       note,
	}
	history.Init()

	if err := db.Create(&history).Error; err != nil {
		return err
	}

	// 更新或创建当前配置
	now := time.Now()
	if current == nil {
		current = &ConfigCurrentModel{
			ID:               ConfigCurrentID,
			CurrentVersionID: history.ID,
			ConfigJSON:       configJSON,
			ConfigHash:       configHash,
			UpdatedAt:        now,
		}
		if err := db.Create(current).Error; err != nil {
			return err
		}
	} else {
		if err := db.Model(current).Updates(map[string]interface{}{
			"current_version_id": history.ID,
			"config_json":        configJSON,
			"config_hash":        configHash,
			"updated_at":         now,
		}).Error; err != nil {
			return err
		}
	}

	// 清理超出限制的旧版本
	if err := cleanupOldConfigVersions(); err != nil {
		log.Printf("清理旧配置版本失败: %v", err)
	}

	return nil
}

// cleanupOldConfigVersions 清理超出限制的旧版本
func cleanupOldConfigVersions() error {
	var count int64
	if err := db.Model(&ConfigHistoryModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count <= MaxConfigHistoryVersions {
		return nil
	}

	// 获取需要删除的版本（最旧的）
	deleteCount := count - MaxConfigHistoryVersions
	var toDelete []ConfigHistoryModel
	if err := db.Order("version ASC").Limit(int(deleteCount)).Find(&toDelete).Error; err != nil {
		return err
	}

	for _, item := range toDelete {
		if err := db.Unscoped().Delete(&item).Error; err != nil {
			log.Printf("删除配置版本 %d 失败: %v", item.Version, err)
		}
	}

	return nil
}

// RollbackToVersion 回滚到指定版本
func RollbackToVersion(version int64) (string, error) {
	history, err := GetConfigHistoryByVersion(version)
	if err != nil {
		return "", err
	}
	if history == nil {
		return "", gorm.ErrRecordNotFound
	}

	// 保存为新版本（source=rollback）
	note := fmt.Sprintf("Rollback from version %d", version)
	if err := SaveConfigVersion(history.ConfigJSON, "rollback", note); err != nil {
		return "", err
	}

	return history.ConfigJSON, nil
}

// ConfigToJSON 将配置对象序列化为 JSON
func ConfigToJSON(config interface{}) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ConfigFromJSON 从 JSON 反序列化配置对象
func ConfigFromJSON(configJSON string, config interface{}) error {
	return json.Unmarshal([]byte(configJSON), config)
}
