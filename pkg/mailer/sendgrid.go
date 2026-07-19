package mailer

import (
	"context"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// SendGridMailer implements the Mailer interface using SendGrid.
type SendGridMailer struct {
	cfg *configs.SendGridConfig
	log *logger.Logger
}

// NewSendGridMailer creates a new SendGrid mailer.
func NewSendGridMailer(cfg *configs.SendGridConfig, log *logger.Logger) *SendGridMailer {
	return &SendGridMailer{cfg: cfg, log: log}
}

// Send sends an email via SendGrid.
func (m *SendGridMailer) Send(ctx context.Context, msg *Message) error {
	// SendGrid SDK integration placeholder.
	m.log.Infof("[SendGrid] Would send email to: %v, subject: %s", msg.To, msg.Subject)
	return nil
}

// SendTemplate sends a templated email via SendGrid.
func (m *SendGridMailer) SendTemplate(ctx context.Context, msg *TemplateMessage) error {
	m.log.Infof("[SendGrid] Would send template '%s' to: %v", msg.TemplateName, msg.To)
	return nil
}

// SendBulk sends multiple emails via SendGrid.
func (m *SendGridMailer) SendBulk(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := m.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
