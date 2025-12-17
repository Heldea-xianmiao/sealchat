package model

import (
	"strings"

	"sealchat/utils"
)

type WebhookIdentityBindingModel struct {
	StringPKBaseModel
	ChannelID       string `json:"channelId" gorm:"size:100;index"`
	IntegrationID   string `json:"integrationId" gorm:"size:100;uniqueIndex:uniq_whk_identity,priority:1;index"`
	BotUserID       string `json:"botUserId" gorm:"size:100;index"`
	Source          string `json:"source" gorm:"size:32;uniqueIndex:uniq_whk_identity,priority:2;index"`
	ExternalActorID string `json:"externalActorId" gorm:"size:128;uniqueIndex:uniq_whk_identity,priority:3;index"`
	IdentityID      string `json:"identityId" gorm:"size:100;index"`
}

func (*WebhookIdentityBindingModel) TableName() string {
	return "webhook_identity_bindings"
}

func WebhookIdentityBindingGet(integrationID, source, externalActorID string) (*WebhookIdentityBindingModel, error) {
	integrationID = strings.TrimSpace(integrationID)
	source = strings.TrimSpace(source)
	externalActorID = strings.TrimSpace(externalActorID)
	if integrationID == "" || source == "" || externalActorID == "" {
		return nil, nil
	}
	var item WebhookIdentityBindingModel
	if err := db.Where("integration_id = ? AND source = ? AND external_actor_id = ?", integrationID, source, externalActorID).
		Limit(1).Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, nil
	}
	return &item, nil
}

func WebhookIdentityBindingUpsert(channelID, integrationID, botUserID, source, externalActorID, identityID string) (*WebhookIdentityBindingModel, error) {
	channelID = strings.TrimSpace(channelID)
	integrationID = strings.TrimSpace(integrationID)
	botUserID = strings.TrimSpace(botUserID)
	source = strings.TrimSpace(source)
	externalActorID = strings.TrimSpace(externalActorID)
	identityID = strings.TrimSpace(identityID)
	if channelID == "" || integrationID == "" || botUserID == "" || source == "" || externalActorID == "" || identityID == "" {
		return nil, nil
	}
	existing, err := WebhookIdentityBindingGet(integrationID, source, externalActorID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		updates := map[string]any{
			"channel_id":  channelID,
			"bot_user_id": botUserID,
			"identity_id": identityID,
		}
		if err := db.Model(&WebhookIdentityBindingModel{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
		existing.ChannelID = channelID
		existing.BotUserID = botUserID
		existing.IdentityID = identityID
		return existing, nil
	}
	item := &WebhookIdentityBindingModel{
		StringPKBaseModel: StringPKBaseModel{
			ID: utils.NewID(),
		},
		ChannelID:       channelID,
		IntegrationID:   integrationID,
		BotUserID:       botUserID,
		Source:          source,
		ExternalActorID: externalActorID,
		IdentityID:      identityID,
	}
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
