package model

import (
	"strings"

	"gorm.io/gorm"
)

// WorldFavoriteModel 记录用户收藏的世界。
type WorldFavoriteModel struct {
	StringPKBaseModel
	WorldID string `json:"worldId" gorm:"size:100;index:idx_world_favorite_user,priority:1"`
	UserID  string `json:"userId" gorm:"size:100;index:idx_world_favorite_user,priority:2"`
}

func (*WorldFavoriteModel) TableName() string {
	return "world_favorites"
}

func (m *WorldFavoriteModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.Init()
	}
	return nil
}

// ListWorldFavoriteIDs 返回用户收藏的世界 ID。
func ListWorldFavoriteIDs(userID string) ([]string, error) {
	if strings.TrimSpace(userID) == "" {
		return []string{}, nil
	}
	var favorites []WorldFavoriteModel
	err := GetDB().Where("user_id = ?", userID).Order("created_at ASC").Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(favorites))
	for _, fav := range favorites {
		if strings.TrimSpace(fav.WorldID) != "" {
			ids = append(ids, fav.WorldID)
		}
	}
	return ids, nil
}

// SetWorldFavorite 设置收藏状态。
func SetWorldFavorite(worldID, userID string, favorite bool) error {
	db := GetDB()
	if favorite {
		var fav WorldFavoriteModel
		if err := db.Where("world_id = ? AND user_id = ?", worldID, userID).Limit(1).Find(&fav).Error; err != nil {
			return err
		}
		if fav.ID != "" {
			return nil
		}
		fav = WorldFavoriteModel{WorldID: worldID, UserID: userID}
		return db.Create(&fav).Error
	}
	return db.Where("world_id = ? AND user_id = ?", worldID, userID).Delete(&WorldFavoriteModel{}).Error
}
