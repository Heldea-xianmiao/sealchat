package model

import (
	"strings"

	"sealchat/utils"
)

type MessageExternalRefModel struct {
	StringPKBaseModel
	ChannelID       string `json:"channelId" gorm:"size:100;uniqueIndex:uniq_msg_external_ref,priority:1;index"`
	Source          string `json:"source" gorm:"size:32;uniqueIndex:uniq_msg_external_ref,priority:2;index"`
	ExternalID      string `json:"externalId" gorm:"size:128;uniqueIndex:uniq_msg_external_ref,priority:3;index"`
	MessageID       string `json:"messageId" gorm:"size:100;index"`
	IntegrationID   string `json:"integrationId" gorm:"size:100;index"`
	ExternalActorID string `json:"externalActorId" gorm:"size:128;index"`
}

func (*MessageExternalRefModel) TableName() string {
	return "message_external_refs"
}

func MessageExternalRefGet(channelID, source, externalID string) (*MessageExternalRefModel, error) {
	channelID = strings.TrimSpace(channelID)
	source = strings.TrimSpace(source)
	externalID = strings.TrimSpace(externalID)
	if channelID == "" || source == "" || externalID == "" {
		return nil, nil
	}
	var item MessageExternalRefModel
	if err := db.Where("channel_id = ? AND source = ? AND external_id = ?", channelID, source, externalID).Limit(1).Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, nil
	}
	return &item, nil
}

func MessageExternalRefGetFirstByMessageID(messageID string) (*MessageExternalRefModel, error) {
	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return nil, nil
	}
	var item MessageExternalRefModel
	if err := db.Where("message_id = ?", messageID).Order("created_at ASC").Limit(1).Find(&item).Error; err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, nil
	}
	return &item, nil
}

func MessageExternalRefUpsert(channelID, source, externalID, messageID, integrationID, externalActorID string) (*MessageExternalRefModel, error) {
	channelID = strings.TrimSpace(channelID)
	source = strings.TrimSpace(source)
	externalID = strings.TrimSpace(externalID)
	messageID = strings.TrimSpace(messageID)
	integrationID = strings.TrimSpace(integrationID)
	externalActorID = strings.TrimSpace(externalActorID)
	if channelID == "" || source == "" || externalID == "" || messageID == "" {
		return nil, nil
	}
	existing, err := MessageExternalRefGet(channelID, source, externalID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		updates := map[string]any{
			"message_id":        messageID,
			"integration_id":    integrationID,
			"external_actor_id": externalActorID,
		}
		if err := db.Model(&MessageExternalRefModel{}).Where("id = ?", existing.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
		existing.MessageID = messageID
		existing.IntegrationID = integrationID
		existing.ExternalActorID = externalActorID
		return existing, nil
	}
	item := &MessageExternalRefModel{
		StringPKBaseModel: StringPKBaseModel{
			ID: utils.NewID(),
		},
		ChannelID:       channelID,
		Source:          source,
		ExternalID:      externalID,
		MessageID:       messageID,
		IntegrationID:   integrationID,
		ExternalActorID: externalActorID,
	}
	if err := db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
