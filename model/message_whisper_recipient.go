package model

type MessageWhisperRecipientModel struct {
	StringPKBaseModel
	MessageID string `json:"message_id" gorm:"size:100;uniqueIndex:idx_mwr_unique,priority:1;index:idx_mwr_message"`
	UserID    string `json:"user_id" gorm:"size:100;uniqueIndex:idx_mwr_unique,priority:2;index:idx_mwr_user"`
}

func (*MessageWhisperRecipientModel) TableName() string {
	return "message_whisper_recipients"
}

// GetWhisperRecipientIDs 获取消息的所有收件人ID
func GetWhisperRecipientIDs(messageID string) []string {
	if messageID == "" {
		return nil
	}
	var recipients []MessageWhisperRecipientModel
	GetDB().Where("message_id = ?", messageID).Find(&recipients)
	ids := make([]string, len(recipients))
	for i, r := range recipients {
		ids[i] = r.UserID
	}
	return ids
}

// GetWhisperRecipientIDsBatch 批量获取消息收件人ID
func GetWhisperRecipientIDsBatch(messageIDs []string) map[string][]string {
	if len(messageIDs) == 0 {
		return nil
	}
	var recipients []MessageWhisperRecipientModel
	GetDB().Where("message_id IN ?", messageIDs).Find(&recipients)
	result := make(map[string][]string)
	for _, r := range recipients {
		result[r.MessageID] = append(result[r.MessageID], r.UserID)
	}
	return result
}

// CreateWhisperRecipients 批量创建收件人记录
func CreateWhisperRecipients(messageID string, userIDs []string) error {
	if messageID == "" || len(userIDs) == 0 {
		return nil
	}
	recipients := make([]MessageWhisperRecipientModel, len(userIDs))
	for i, uid := range userIDs {
		recipients[i] = MessageWhisperRecipientModel{
			MessageID: messageID,
			UserID:    uid,
		}
		recipients[i].Init()
	}
	return GetDB().Create(&recipients).Error
}

// HasWhisperRecipient 判断用户是否为消息收件人
func HasWhisperRecipient(messageID, userID string) bool {
	if messageID == "" || userID == "" {
		return false
	}
	var count int64
	GetDB().Model(&MessageWhisperRecipientModel{}).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		Count(&count)
	return count > 0
}
