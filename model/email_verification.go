package model

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"time"

	"golang.org/x/crypto/blake2s"
	"gorm.io/gorm"

	"sealchat/utils"
)

type EmailVerificationScene string

const (
	EmailSceneSignup        EmailVerificationScene = "signup"
	EmailScenePasswordReset EmailVerificationScene = "password_reset"
	EmailSceneBind          EmailVerificationScene = "bind"
)

type EmailVerificationCodeModel struct {
	StringPKBaseModel
	Scene        EmailVerificationScene `gorm:"size:32;not null;index:idx_email_scene" json:"scene"`
	Email        string                 `gorm:"size:254;not null;index:idx_email_scene" json:"email"`
	UserID       *string                `gorm:"size:32;index" json:"userId,omitempty"`
	CodeHash     string                 `gorm:"size:64;not null" json:"-"`
	CodeSalt     string                 `gorm:"size:32;not null" json:"-"`
	ExpiresAt    time.Time              `gorm:"not null" json:"expiresAt"`
	ConsumedAt   *time.Time             `json:"consumedAt,omitempty"`
	AttemptCount int                    `gorm:"default:0" json:"attemptCount"`
	MaxAttempts  int                    `gorm:"default:5" json:"maxAttempts"`
	SendIP       string                 `gorm:"size:45" json:"-"`
	UserAgent    string                 `gorm:"size:512" json:"-"`
}

func (*EmailVerificationCodeModel) TableName() string {
	return "email_verification_codes"
}

var (
	ErrCodeExpired      = errors.New("验证码已过期")
	ErrCodeConsumed     = errors.New("验证码已使用")
	ErrCodeInvalid      = errors.New("验证码无效")
	ErrCodeMaxAttempts  = errors.New("验证码尝试次数过多")
	ErrCodeNotFound     = errors.New("验证码不存在")
	ErrEmailRateLimited = errors.New("发送频率过高，请稍后再试")
	ErrCodeUserMismatch = errors.New("验证码用户不匹配")
)

func generateNumericCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

func hashVerificationCode(code, salt string) string {
	hashBytes := blake2s.Sum256([]byte(code + salt))
	return base64.RawStdEncoding.EncodeToString(hashBytes[:])
}

func EmailVerificationCreate(scene EmailVerificationScene, email string, userID *string, sendIP, userAgent string) (*EmailVerificationCodeModel, string, error) {
	cfg := utils.GetConfig().EmailAuth

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, "", err
	}
	saltStr := base64.RawStdEncoding.EncodeToString(salt)

	code, err := generateNumericCode(cfg.CodeLength)
	if err != nil {
		return nil, "", err
	}
	codeHash := hashVerificationCode(code, saltStr)

	model := &EmailVerificationCodeModel{
		Scene:       scene,
		Email:       email,
		UserID:      userID,
		CodeHash:    codeHash,
		CodeSalt:    saltStr,
		ExpiresAt:   time.Now().Add(time.Duration(cfg.CodeTTLSeconds) * time.Second),
		MaxAttempts: cfg.MaxAttempts,
		SendIP:      sendIP,
		UserAgent:   userAgent,
	}
	model.ID = utils.NewID()

	if err := db.Create(model).Error; err != nil {
		return nil, "", err
	}

	return model, code, nil
}

func EmailVerificationVerify(scene EmailVerificationScene, email, code string) (*EmailVerificationCodeModel, error) {
	return EmailVerificationVerifyWithUserID(scene, email, code, nil)
}

func EmailVerificationVerifyWithUserID(scene EmailVerificationScene, email, code string, requiredUserID *string) (*EmailVerificationCodeModel, error) {
	var model EmailVerificationCodeModel
	err := db.Where("scene = ? AND email = ? AND consumed_at IS NULL", scene, email).
		Order("created_at DESC").
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCodeNotFound
		}
		return nil, err
	}

	if model.AttemptCount >= model.MaxAttempts {
		return nil, ErrCodeMaxAttempts
	}

	if err := db.Model(&model).Update("attempt_count", model.AttemptCount+1).Error; err != nil {
		return nil, err
	}

	if time.Now().After(model.ExpiresAt) {
		return nil, ErrCodeExpired
	}

	// 验证 userID 匹配（用于绑定场景）
	if requiredUserID != nil && model.UserID != nil && *requiredUserID != *model.UserID {
		return nil, ErrCodeUserMismatch
	}

	expectedHash := hashVerificationCode(code, model.CodeSalt)
	if expectedHash != model.CodeHash {
		return nil, ErrCodeInvalid
	}

	return &model, nil
}

func EmailVerificationConsume(id string) error {
	now := time.Now()
	return db.Model(&EmailVerificationCodeModel{}).
		Where("id = ?", id).
		Update("consumed_at", &now).Error
}

func EmailVerificationCheckRateLimit(email, sendIP string, limitPerIP int) error {
	if limitPerIP <= 0 {
		return nil // 0 或负数表示不限制
	}

	var count int64
	oneHourAgo := time.Now().Add(-time.Hour)

	if err := db.Model(&EmailVerificationCodeModel{}).
		Where("(email = ? OR send_ip = ?) AND created_at > ?", email, sendIP, oneHourAgo).
		Count(&count).Error; err != nil {
		return err
	}

	if count >= int64(limitPerIP) {
		return ErrEmailRateLimited
	}
	return nil
}

func EmailVerificationCleanExpired() error {
	expiryThreshold := time.Now().Add(-24 * time.Hour)
	return db.Where("expires_at < ?", expiryThreshold).Delete(&EmailVerificationCodeModel{}).Error
}
