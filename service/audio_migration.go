package service

import (
	"sort"

	"gorm.io/gorm"

	"sealchat/model"
)

type AudioFolderMigrationStats struct {
	FolderTotal       int64 `json:"folderTotal"`
	WorldFolderTotal  int64 `json:"worldFolderTotal"`
	WorldAssetTotal   int64 `json:"worldAssetTotal"`
	DuplicatePaths    int64 `json:"duplicatePaths"`
	DuplicateFolders  int64 `json:"duplicateFolders"`
	MergedFolders     int64 `json:"mergedFolders"`
	UpdatedAssets     int64 `json:"updatedAssets"`
	UpdatedScenes     int64 `json:"updatedScenes"`
	UpdatedChildren   int64 `json:"updatedChildren"`
}

type AudioFolderMergeResult struct {
	Path            string   `json:"path"`
	KeptID          string   `json:"keptId"`
	MergedIDs       []string `json:"mergedIds"`
	AssetsUpdated   int64    `json:"assetsUpdated"`
	ScenesUpdated   int64    `json:"scenesUpdated"`
	ChildrenUpdated int64    `json:"childrenUpdated"`
}

func GetAudioFolderMigrationPreview() (*AudioFolderMigrationStats, error) {
	db := model.GetDB()
	stats := &AudioFolderMigrationStats{}

	if err := db.Model(&model.AudioFolder{}).Count(&stats.FolderTotal).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&model.AudioFolder{}).
		Where("scope = ?", model.AudioScopeWorld).
		Count(&stats.WorldFolderTotal).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&model.AudioAsset{}).
		Where("scope = ?", model.AudioScopeWorld).
		Count(&stats.WorldAssetTotal).Error; err != nil {
		return nil, err
	}

	type pathCount struct {
		Path string
		Cnt  int64
	}
	var duplicates []pathCount
	if err := db.Model(&model.AudioFolder{}).
		Select("path, COUNT(*) as cnt").
		Group("path").
		Having("COUNT(*) > 1").
		Scan(&duplicates).Error; err != nil {
		return nil, err
	}
	stats.DuplicatePaths = int64(len(duplicates))
	for _, item := range duplicates {
		if item.Cnt > 1 {
			stats.DuplicateFolders += item.Cnt - 1
		}
	}
	return stats, nil
}

func MigrateAudioFoldersToCommon(dryRun bool) (*AudioFolderMigrationStats, []AudioFolderMergeResult, error) {
	stats, err := GetAudioFolderMigrationPreview()
	if err != nil {
		return nil, nil, err
	}

	db := model.GetDB()
	var folders []model.AudioFolder
	if err := db.Find(&folders).Error; err != nil {
		return nil, nil, err
	}

	grouped := map[string][]model.AudioFolder{}
	for _, folder := range folders {
		grouped[folder.Path] = append(grouped[folder.Path], folder)
	}

	results := make([]AudioFolderMergeResult, 0)
	err = db.Transaction(func(tx *gorm.DB) error {
		for path, group := range grouped {
			if len(group) <= 1 {
				continue
			}
			sort.Slice(group, func(i, j int) bool {
				if group[i].CreatedAt.Equal(group[j].CreatedAt) {
					return group[i].ID < group[j].ID
				}
				return group[i].CreatedAt.Before(group[j].CreatedAt)
			})
			kept := group[0]
			mergeResult := AudioFolderMergeResult{
				Path:   path,
				KeptID: kept.ID,
			}
			for _, dup := range group[1:] {
				mergeResult.MergedIDs = append(mergeResult.MergedIDs, dup.ID)
				var assetsUpdated int64
				var scenesUpdated int64
				var childrenUpdated int64
				if dryRun {
					if err := tx.Model(&model.AudioAsset{}).Where("folder_id = ?", dup.ID).Count(&assetsUpdated).Error; err != nil {
						return err
					}
					if err := tx.Model(&model.AudioScene{}).Where("folder_id = ?", dup.ID).Count(&scenesUpdated).Error; err != nil {
						return err
					}
					if err := tx.Model(&model.AudioFolder{}).Where("parent_id = ?", dup.ID).Count(&childrenUpdated).Error; err != nil {
						return err
					}
				} else {
					res := tx.Model(&model.AudioAsset{}).Where("folder_id = ?", dup.ID).Update("folder_id", kept.ID)
					if res.Error != nil {
						return res.Error
					}
					assetsUpdated = res.RowsAffected
					res = tx.Model(&model.AudioScene{}).Where("folder_id = ?", dup.ID).Update("folder_id", kept.ID)
					if res.Error != nil {
						return res.Error
					}
					scenesUpdated = res.RowsAffected
					res = tx.Model(&model.AudioFolder{}).Where("parent_id = ?", dup.ID).Update("parent_id", kept.ID)
					if res.Error != nil {
						return res.Error
					}
					childrenUpdated = res.RowsAffected
					if err := tx.Delete(&model.AudioFolder{}, "id = ?", dup.ID).Error; err != nil {
						return err
					}
				}
				mergeResult.AssetsUpdated += assetsUpdated
				mergeResult.ScenesUpdated += scenesUpdated
				mergeResult.ChildrenUpdated += childrenUpdated
				stats.MergedFolders += 1
				stats.UpdatedAssets += assetsUpdated
				stats.UpdatedScenes += scenesUpdated
				stats.UpdatedChildren += childrenUpdated
			}
			results = append(results, mergeResult)
		}

		if dryRun {
			return nil
		}

		if err := tx.Model(&model.AudioFolder{}).
			Updates(map[string]interface{}{"scope": model.AudioScopeCommon, "world_id": nil}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.AudioAsset{}).
			Where("scope = ?", model.AudioScopeWorld).
			Updates(map[string]interface{}{"scope": model.AudioScopeCommon, "world_id": nil}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return stats, results, nil
}
