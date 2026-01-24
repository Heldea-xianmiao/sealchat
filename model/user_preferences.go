package model

import (
	"gorm.io/gorm/clause"

	"sealchat/utils"
)

// UserPreferenceModel 用户偏好表
type UserPreferenceModel struct {
	StringPKBaseModel
	UserID    string `json:"userId" gorm:"index:idx_user_pref_user_key,unique;not null"`
	PrefKey   string `json:"key" gorm:"column:pref_key;size:64;index:idx_user_pref_user_key,unique;not null"`
	PrefValue string `json:"value" gorm:"column:pref_value;size:4096;not null"`
}

func (*UserPreferenceModel) TableName() string {
	return "user_preferences"
}

// UserPreferenceGet 获取用户偏好
func UserPreferenceGet(userID, key string) (*UserPreferenceModel, error) {
	var record UserPreferenceModel
	err := db.Where("user_id = ? AND pref_key = ?", userID, key).Limit(1).Find(&record).Error
	if err != nil {
		return nil, err
	}
	if record.ID == "" {
		return nil, nil
	}
	return &record, nil
}

// UserPreferenceUpsert 创建或更新用户偏好
func UserPreferenceUpsert(userID, key, value string) (*UserPreferenceModel, error) {
	record := &UserPreferenceModel{
		StringPKBaseModel: StringPKBaseModel{ID: utils.NewID()},
		UserID:            userID,
		PrefKey:           key,
		PrefValue:         value,
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "pref_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"pref_value", "updated_at"}),
	}).Create(record).Error
	if err != nil {
		return nil, err
	}
	// 直接返回请求值，避免读偏差
	record.PrefKey = key
	record.PrefValue = value
	return record, nil
}
