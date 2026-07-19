package mailer

import (
	"context"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// MailgunMailer implements the Mailer interface using Mailgun.
// Requires the mailgun-go SDK when API keys are provided.
type MailgunMailer struct {
	cfg *configs.MailgunConfig
	log *logger.Logger
}

// NewMailgunMailer creates a new Mailgun mailer.
func NewMailgunMailer(cfg *configs.MailgunConfig, log *logger.Logger) *MailgunMailer {
	return &MailgunMailer{cfg: cfg, log: log}
}

// Send sends an email via Mailgun.
func (m *MailgunMailer) Send(ctx context.Context, msg *Message) error {
	// Mailgun SDK integration placeholder.
	// When mailgun-go is added: mg := mailgun.NewMailgun(m.cfg.Domain, m.cfg.APIKey)
	m.log.Infof("[Mailgun] Would send email to: %v, subject: %s", msg.To, msg.Subject)
	return nil
}

// SendTemplate sends a templated email via Mailgun.
func (m *MailgunMailer) SendTemplate(ctx context.Context, msg *TemplateMessage) error {
	m.log.Infof("[Mailgun] Would send template '%s' to: %v", msg.TemplateName, msg.To)
	return nil
}

// SendBulk sends multiple emails via Mailgun.
func (m *MailgunMailer) SendBulk(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := m.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
