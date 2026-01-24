package service

import (
	"errors"
	"fmt"
	"strings"

	"sealchat/model"
	"sealchat/utils"
)

var (
	ErrEmailAuthDisabled = errors.New("邮箱认证功能未启用")
	ErrSMTPNotConfigured = errors.New("SMTP 未配置")
	ErrEmailAlreadyUsed  = errors.New("该邮箱已被使用")
	ErrUserNotFound      = errors.New("用户不存在")
)

type EmailAuthService struct {
	cfg       utils.EmailAuthConfig
	emailSvc  *EmailService
	fromEmail utils.SMTPConfig
}

func NewEmailAuthService() *EmailAuthService {
	cfg := utils.GetConfig()
	if cfg == nil {
		return nil
	}

	// 邮箱验证码功能复用 emailNotification.smtp 配置
	smtpCfg := cfg.EmailNotification.SMTP

	return &EmailAuthService{
		cfg:       cfg.EmailAuth,
		emailSvc:  NewEmailService(smtpCfg),
		fromEmail: smtpCfg,
	}
}

func (s *EmailAuthService) IsEnabled() bool {
	return s != nil && s.cfg.Enabled && s.emailSvc != nil && s.emailSvc.IsConfigured()
}

func (s *EmailAuthService) SendSignupCode(email, sendIP, userAgent string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	exists, err := model.UserExistsByEmail(email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyUsed
	}

	if err := model.EmailVerificationCheckRateLimit(email, sendIP, s.cfg.RateLimitPerIP); err != nil {
		return err
	}

	_, code, err := model.EmailVerificationCreate(model.EmailSceneSignup, email, nil, sendIP, userAgent)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(email, code, "注册验证")
}

func (s *EmailAuthService) SendPasswordResetCode(account, sendIP, userAgent string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	user, err := model.UserGetByEmailOrUsername(account)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	email := ""
	if user.Email != nil && *user.Email != "" {
		email = *user.Email
	} else {
		return nil
	}

	if err := model.EmailVerificationCheckRateLimit(email, sendIP, s.cfg.RateLimitPerIP); err != nil {
		return err
	}

	_, code, err := model.EmailVerificationCreate(model.EmailScenePasswordReset, email, &user.ID, sendIP, userAgent)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(email, code, "密码重置")
}

func (s *EmailAuthService) SendBindCode(userID, email, sendIP, userAgent string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	exists, err := model.UserExistsByEmail(email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyUsed
	}

	if err := model.EmailVerificationCheckRateLimit(email, sendIP, s.cfg.RateLimitPerIP); err != nil {
		return err
	}

	_, code, err := model.EmailVerificationCreate(model.EmailSceneBind, email, &userID, sendIP, userAgent)
	if err != nil {
		return err
	}

	return s.sendVerificationEmail(email, code, "邮箱绑定")
}

func (s *EmailAuthService) VerifySignupCode(email, code string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	record, err := model.EmailVerificationVerify(model.EmailSceneSignup, email, code)
	if err != nil {
		return err
	}

	return model.EmailVerificationConsume(record.ID)
}

func (s *EmailAuthService) VerifyAndResetPassword(account, code, newPassword string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	user, err := model.UserGetByEmailOrUsername(account)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	if email == "" {
		email = account
	}

	record, err := model.EmailVerificationVerify(model.EmailScenePasswordReset, email, code)
	if err != nil {
		return err
	}

	if err := model.EmailVerificationConsume(record.ID); err != nil {
		return err
	}

	return model.UserUpdatePassword(user.ID, newPassword)
}

func (s *EmailAuthService) VerifyAndBindEmail(userID, email, code string) error {
	if !s.IsEnabled() {
		return ErrEmailAuthDisabled
	}

	exists, err := model.UserExistsByEmail(email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyUsed
	}

	// 使用带 userID 验证的函数，确保验证码是当前用户发起的
	record, err := model.EmailVerificationVerifyWithUserID(model.EmailSceneBind, email, code, &userID)
	if err != nil {
		return err
	}

	if err := model.EmailVerificationConsume(record.ID); err != nil {
		return err
	}

	if err := model.UserBindEmail(userID, email); err != nil {
		if errors.Is(err, model.ErrEmailAlreadyUsed) {
			return ErrEmailAlreadyUsed
		}
		return err
	}

	return nil
}

func (s *EmailAuthService) sendVerificationEmail(to, code, purpose string) error {
	subject := fmt.Sprintf("SealChat %s验证码", purpose)
	body := s.buildVerificationEmailHTML(code, purpose)
	return s.emailSvc.SendEmail(to, subject, body)
}

func (s *EmailAuthService) buildVerificationEmailHTML(code, purpose string) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; padding: 20px; }
.container { max-width: 500px; margin: 0 auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #fff; padding: 24px; text-align: center; }
.header h1 { margin: 0; font-size: 24px; }
.content { padding: 32px; text-align: center; }
.code { font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #667eea; padding: 16px 24px; background: #f8f9ff; border-radius: 8px; display: inline-block; margin: 16px 0; }
.footer { padding: 16px 24px; background: #f9f9f9; text-align: center; font-size: 12px; color: #999; }
.note { color: #666; font-size: 14px; margin-top: 16px; }
</style>
</head>
<body>
<div class="container">
<div class="header">
<h1>🔐 `)
	sb.WriteString(escapeHTML(purpose))
	sb.WriteString(`验证码</h1>
</div>
<div class="content">
<p>您正在进行`)
	sb.WriteString(escapeHTML(purpose))
	sb.WriteString(`操作，请使用以下验证码：</p>
<div class="code">`)
	sb.WriteString(escapeHTML(code))
	sb.WriteString(`</div>
<p class="note">验证码有效期为 5 分钟，请勿泄露给他人。</p>
</div>
<div class="footer">
此邮件由 SealChat 自动发送，请勿直接回复。
</div>
</div>
</body>
</html>`)

	return sb.String()
}
