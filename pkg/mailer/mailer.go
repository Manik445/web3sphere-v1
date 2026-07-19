package mailer

import (
	"context"
	"fmt"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// Mailer defines the email sending interface.
type Mailer interface {
	Send(ctx context.Context, msg *Message) error
	SendTemplate(ctx context.Context, msg *TemplateMessage) error
	SendBulk(ctx context.Context, messages []*Message) error
}

// Message represents a simple email message.
type Message struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// TemplateMessage represents a templated email message.
type TemplateMessage struct {
	To           []string
	Subject      string
	TemplateName string
	Data         map[string]interface{}
}

// New creates a Mailer based on the configured provider.
func New(cfg *configs.Config, log *logger.Logger) (Mailer, error) {
	switch cfg.Mail.Provider {
	case "smtp":
		return NewSMTPMailer(&cfg.SMTP, log), nil
	case "mailgun":
		return NewMailgunMailer(&cfg.Mailgun, log), nil
	case "sendgrid":
		return NewSendGridMailer(&cfg.SendGrid, log), nil
	default:
		return nil, fmt.Errorf("unsupported mail provider: %s", cfg.Mail.Provider)
	}
}
