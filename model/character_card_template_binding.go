package model

import "strings"

const (
	CharacterCardTemplateModeManaged  = "managed"
	CharacterCardTemplateModeDetached = "detached"
)

type CharacterCardTemplateBindingModel struct {
	StringPKBaseModel
	UserID           string `json:"userId" gorm:"size:100;index:idx_cc_template_binding_unique,priority:1"`
	ChannelID        string `json:"channelId" gorm:"size:100;index:idx_cc_template_binding_unique,priority:2;index"`
	ExternalCardID   string `json:"externalCardId" gorm:"size:100;index:idx_cc_template_binding_unique,priority:3"`
	CardName         string `json:"cardName" gorm:"size:64"`
	SheetType        string `json:"sheetType" gorm:"size:32;index"`
	Mode             string `json:"mode" gorm:"size:16;index"`
	TemplateID       string `json:"templateId" gorm:"size:100;index"`
	TemplateSnapshot string `json:"templateSnapshot" gorm:"type:text"`
}

func (*CharacterCardTemplateBindingModel) TableName() string {
	return "character_card_template_bindings"
}

func CharacterCardTemplateBindingList(userID string, channelID string) ([]*CharacterCardTemplateBindingModel, error) {
	var items []*CharacterCardTemplateBindingModel
	err := db.Where("user_id = ? AND channel_id = ?", userID, strings.TrimSpace(channelID)).
		Order("updated_at desc").
		Find(&items).Error
	return items, err
}

func CharacterCardTemplateBindingGet(userID string, channelID string, externalCardID string) (*CharacterCardTemplateBindingModel, error) {
	item := &CharacterCardTemplateBindingModel{}
	err := db.Where("user_id = ? AND channel_id = ? AND external_card_id = ?", userID, strings.TrimSpace(channelID), strings.TrimSpace(externalCardID)).
		Take(item).Error
	if err != nil {
		return nil, err
	}
	return item, nil
}

func CharacterCardTemplateBindingCreate(item *CharacterCardTemplateBindingModel) error {
	return db.Create(item).Error
}

func CharacterCardTemplateBindingUpdate(id string, values map[string]any) error {
	if len(values) == 0 {
		return nil
	}
	return db.Model(&CharacterCardTemplateBindingModel{}).Where("id = ?", id).Updates(values).Error
}

func CharacterCardTemplateBindingDetachByTemplateID(userID string, templateID string, snapshot string) error {
	updates := map[string]any{
		"mode":              CharacterCardTemplateModeDetached,
		"template_id":       "",
		"template_snapshot": snapshot,
	}
	return db.Model(&CharacterCardTemplateBindingModel{}).
		Where("user_id = ?", userID).
		Where("template_id = ?", strings.TrimSpace(templateID)).
		Where("mode = ?", CharacterCardTemplateModeManaged).
		Updates(updates).Error
}
