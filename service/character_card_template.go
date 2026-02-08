package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"sealchat/model"
)

const (
	CharacterCardTemplateDefaultScopeGlobal = "global"
	CharacterCardTemplateDefaultScopeSheet  = "sheet"
)

type CharacterCardTemplateInput struct {
	Name            string
	SheetType       string
	Content         string
	IsGlobalDefault bool
	IsSheetDefault  bool
}

type CharacterCardTemplateUpdateInput struct {
	Name            *string
	SheetType       *string
	Content         *string
	IsGlobalDefault *bool
	IsSheetDefault  *bool
}

type CharacterCardTemplateBindingInput struct {
	ChannelID        string
	ExternalCardID   string
	CardName         string
	SheetType        string
	Mode             string
	TemplateID       string
	TemplateSnapshot string
}

func normalizeCharacterCardTemplateInput(input *CharacterCardTemplateInput) error {
	if input == nil {
		return errors.New("参数错误")
	}
	input.Name = strings.TrimSpace(input.Name)
	input.SheetType = strings.TrimSpace(input.SheetType)
	input.Content = strings.TrimSpace(input.Content)
	if input.Name == "" {
		return errors.New("模板名称不能为空")
	}
	if len([]rune(input.Name)) > 100 {
		return errors.New("模板名称长度需在100个字符以内")
	}
	if input.SheetType != "" && len([]rune(input.SheetType)) > 32 {
		return errors.New("模板规制类型长度需在32个字符以内")
	}
	if input.Content == "" {
		return errors.New("模板内容不能为空")
	}
	return nil
}

func normalizeCharacterCardTemplateUpdateInput(input *CharacterCardTemplateUpdateInput) error {
	if input == nil {
		return errors.New("参数错误")
	}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return errors.New("模板名称不能为空")
		}
		if len([]rune(name)) > 100 {
			return errors.New("模板名称长度需在100个字符以内")
		}
		input.Name = &name
	}
	if input.SheetType != nil {
		sheetType := strings.TrimSpace(*input.SheetType)
		if sheetType != "" && len([]rune(sheetType)) > 32 {
			return errors.New("模板规制类型长度需在32个字符以内")
		}
		input.SheetType = &sheetType
	}
	if input.Content != nil {
		content := strings.TrimSpace(*input.Content)
		if content == "" {
			return errors.New("模板内容不能为空")
		}
		input.Content = &content
	}
	return nil
}

func CharacterCardTemplateList(userID string, sheetType string) ([]*model.CharacterCardTemplateModel, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("缺少用户ID")
	}
	return model.CharacterCardTemplateList(userID, sheetType)
}

func CharacterCardTemplateGet(userID string, templateID string) (*model.CharacterCardTemplateModel, error) {
	template, err := model.CharacterCardTemplateGetByID(strings.TrimSpace(templateID))
	if err != nil {
		return nil, err
	}
	if template.UserID != userID {
		return nil, errors.New("无权访问该模板")
	}
	return template, nil
}

func CharacterCardTemplateCreate(userID string, input *CharacterCardTemplateInput) (*model.CharacterCardTemplateModel, error) {
	if err := normalizeCharacterCardTemplateInput(input); err != nil {
		return nil, err
	}
	item := &model.CharacterCardTemplateModel{
		UserID:          userID,
		Name:            input.Name,
		SheetType:       input.SheetType,
		Content:         input.Content,
		IsGlobalDefault: input.IsGlobalDefault,
		IsSheetDefault:  input.IsSheetDefault,
	}
	err := model.GetDB().Transaction(func(tx *gorm.DB) error {
		if input.IsGlobalDefault {
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("is_global_default = ?", true).
				Update("is_global_default", false).Error; err != nil {
				return err
			}
		}
		if input.IsSheetDefault {
			if item.SheetType == "" {
				return errors.New("设置规制默认模板时，sheetType 不能为空")
			}
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("sheet_type = ?", item.SheetType).
				Where("is_sheet_default = ?", true).
				Update("is_sheet_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(item).Error
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func CharacterCardTemplateUpdate(userID string, templateID string, input *CharacterCardTemplateUpdateInput) (*model.CharacterCardTemplateModel, error) {
	if err := normalizeCharacterCardTemplateUpdateInput(input); err != nil {
		return nil, err
	}
	template, err := CharacterCardTemplateGet(userID, templateID)
	if err != nil {
		return nil, err
	}
	values := map[string]any{}
	nextSheetType := template.SheetType
	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.SheetType != nil {
		nextSheetType = *input.SheetType
		values["sheet_type"] = nextSheetType
	}
	if input.Content != nil {
		values["content"] = *input.Content
	}
	if input.IsGlobalDefault != nil {
		values["is_global_default"] = *input.IsGlobalDefault
	}
	if input.IsSheetDefault != nil {
		values["is_sheet_default"] = *input.IsSheetDefault
	}
	err = model.GetDB().Transaction(func(tx *gorm.DB) error {
		if input.IsGlobalDefault != nil && *input.IsGlobalDefault {
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("is_global_default = ?", true).
				Where("id <> ?", template.ID).
				Update("is_global_default", false).Error; err != nil {
				return err
			}
		}
		if input.IsSheetDefault != nil && *input.IsSheetDefault {
			if strings.TrimSpace(nextSheetType) == "" {
				return errors.New("设置规制默认模板时，sheetType 不能为空")
			}
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("sheet_type = ?", nextSheetType).
				Where("is_sheet_default = ?", true).
				Where("id <> ?", template.ID).
				Update("is_sheet_default", false).Error; err != nil {
				return err
			}
		}
		if input.SheetType != nil && template.IsSheetDefault && nextSheetType != "" {
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("id = ?", template.ID).
				Update("is_sheet_default", false).Error; err != nil {
				return err
			}
			values["is_sheet_default"] = true
		}
		return tx.Model(&model.CharacterCardTemplateModel{}).Where("id = ?", template.ID).Updates(values).Error
	})
	if err != nil {
		return nil, err
	}
	return model.CharacterCardTemplateGetByID(template.ID)
}

func CharacterCardTemplateDelete(userID string, templateID string) error {
	template, err := CharacterCardTemplateGet(userID, templateID)
	if err != nil {
		return err
	}
	return model.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.CharacterCardTemplateBindingModel{}).
			Where("user_id = ?", userID).
			Where("template_id = ?", template.ID).
			Where("mode = ?", model.CharacterCardTemplateModeManaged).
			Updates(map[string]any{
				"mode":              model.CharacterCardTemplateModeDetached,
				"template_id":       "",
				"template_snapshot": template.Content,
			}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", template.ID).Delete(&model.CharacterCardTemplateModel{}).Error
	})
}

func CharacterCardTemplateSetDefault(userID string, templateID string, scope string) (*model.CharacterCardTemplateModel, error) {
	template, err := CharacterCardTemplateGet(userID, templateID)
	if err != nil {
		return nil, err
	}
	scope = strings.TrimSpace(scope)
	err = model.GetDB().Transaction(func(tx *gorm.DB) error {
		switch scope {
		case CharacterCardTemplateDefaultScopeGlobal:
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("is_global_default = ?", true).
				Update("is_global_default", false).Error; err != nil {
				return err
			}
			return tx.Model(&model.CharacterCardTemplateModel{}).Where("id = ?", template.ID).Updates(map[string]any{
				"is_global_default": true,
			}).Error
		case CharacterCardTemplateDefaultScopeSheet:
			if strings.TrimSpace(template.SheetType) == "" {
				return errors.New("当前模板缺少 sheetType，无法设为规制默认模板")
			}
			if err := tx.Model(&model.CharacterCardTemplateModel{}).
				Where("user_id = ?", userID).
				Where("sheet_type = ?", template.SheetType).
				Where("is_sheet_default = ?", true).
				Update("is_sheet_default", false).Error; err != nil {
				return err
			}
			return tx.Model(&model.CharacterCardTemplateModel{}).Where("id = ?", template.ID).Updates(map[string]any{
				"is_sheet_default": true,
			}).Error
		default:
			return errors.New("默认模板作用域无效")
		}
	})
	if err != nil {
		return nil, err
	}
	return model.CharacterCardTemplateGetByID(template.ID)
}

func normalizeCharacterCardTemplateBindingInput(input *CharacterCardTemplateBindingInput) error {
	if input == nil {
		return errors.New("参数错误")
	}
	input.ChannelID = strings.TrimSpace(input.ChannelID)
	input.ExternalCardID = strings.TrimSpace(input.ExternalCardID)
	input.CardName = strings.TrimSpace(input.CardName)
	input.SheetType = strings.TrimSpace(input.SheetType)
	input.Mode = strings.TrimSpace(input.Mode)
	input.TemplateID = strings.TrimSpace(input.TemplateID)
	input.TemplateSnapshot = strings.TrimSpace(input.TemplateSnapshot)
	if input.ChannelID == "" {
		return errors.New("缺少频道ID")
	}
	if input.ExternalCardID == "" {
		return errors.New("缺少角色卡ID")
	}
	if input.Mode == "" {
		input.Mode = model.CharacterCardTemplateModeManaged
	}
	switch input.Mode {
	case model.CharacterCardTemplateModeManaged:
		if input.TemplateID == "" {
			return errors.New("模板库模式下 templateId 不能为空")
		}
	case model.CharacterCardTemplateModeDetached:
		if input.TemplateSnapshot == "" {
			return errors.New("自定义模式下模板内容不能为空")
		}
	default:
		return errors.New("模板绑定模式无效")
	}
	return nil
}

func CharacterCardTemplateBindingList(userID string, channelID string) ([]*model.CharacterCardTemplateBindingModel, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("缺少用户ID")
	}
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return nil, errors.New("缺少频道ID")
	}
	if err := ensureChannelMembership(userID, channelID); err != nil {
		return nil, err
	}
	return model.CharacterCardTemplateBindingList(userID, channelID)
}

func CharacterCardTemplateBindingUpsert(userID string, input *CharacterCardTemplateBindingInput) (*model.CharacterCardTemplateBindingModel, error) {
	if err := normalizeCharacterCardTemplateBindingInput(input); err != nil {
		return nil, err
	}
	if err := ensureChannelMembership(userID, input.ChannelID); err != nil {
		return nil, err
	}
	if input.Mode == model.CharacterCardTemplateModeManaged {
		template, err := model.CharacterCardTemplateGetByID(input.TemplateID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("模板不存在")
			}
			return nil, err
		}
		if template.UserID != userID {
			return nil, errors.New("无权绑定该模板")
		}
		if input.SheetType == "" {
			input.SheetType = template.SheetType
		}
		input.TemplateSnapshot = ""
	}

	existing, err := model.CharacterCardTemplateBindingGet(userID, input.ChannelID, input.ExternalCardID)
	if err == nil {
		updates := map[string]any{
			"card_name":         input.CardName,
			"sheet_type":        input.SheetType,
			"mode":              input.Mode,
			"template_id":       input.TemplateID,
			"template_snapshot": input.TemplateSnapshot,
		}
		if err := model.CharacterCardTemplateBindingUpdate(existing.ID, updates); err != nil {
			return nil, err
		}
		return model.CharacterCardTemplateBindingGet(userID, input.ChannelID, input.ExternalCardID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := &model.CharacterCardTemplateBindingModel{
		UserID:           userID,
		ChannelID:        input.ChannelID,
		ExternalCardID:   input.ExternalCardID,
		CardName:         input.CardName,
		SheetType:        input.SheetType,
		Mode:             input.Mode,
		TemplateID:       input.TemplateID,
		TemplateSnapshot: input.TemplateSnapshot,
	}
	if err := model.CharacterCardTemplateBindingCreate(item); err != nil {
		return nil, err
	}
	return item, nil
}
