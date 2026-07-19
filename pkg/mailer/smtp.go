package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// SMTPMailer implements the Mailer interface using SMTP.
type SMTPMailer struct {
	cfg *configs.SMTPConfig
	log *logger.Logger
}

// NewSMTPMailer creates a new SMTP mailer.
func NewSMTPMailer(cfg *configs.SMTPConfig, log *logger.Logger) *SMTPMailer {
	return &SMTPMailer{cfg: cfg, log: log}
}

// Send sends a simple email via SMTP.
func (m *SMTPMailer) Send(ctx context.Context, msg *Message) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	from := fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.From)

	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	for _, to := range msg.To {
		body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=utf-8\r\n\r\n%s",
			from, to, msg.Subject, contentType, msg.Body)

		var auth smtp.Auth
		if m.cfg.User != "" {
			auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
		}

		if m.cfg.Encryption == "tls" {
			err := m.sendTLS(addr, auth, m.cfg.From, to, []byte(body))
			if err != nil {
				m.log.Errorf("Failed to send email to %s: %v", to, err)
				return err
			}
		} else {
			err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, []byte(body))
			if err != nil {
				m.log.Errorf("Failed to send email to %s: %v", to, err)
				return err
			}
		}
	}

	m.log.Infof("Email sent to: %s", strings.Join(msg.To, ", "))
	return nil
}

// SendTemplate sends a templated email. For SMTP, it renders the template into an HTML body.
func (m *SMTPMailer) SendTemplate(ctx context.Context, msg *TemplateMessage) error {
	// For SMTP, render template data into a simple HTML body.
	body := renderTemplate(msg.TemplateName, msg.Data)
	return m.Send(ctx, &Message{
		To:      msg.To,
		Subject: msg.Subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendBulk sends multiple emails.
func (m *SMTPMailer) SendBulk(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := m.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (m *SMTPMailer) sendTLS(addr string, auth smtp.Auth, from, to string, body []byte) error {
	host := strings.Split(addr, ":")[0]
	tlsConfig := &tls.Config{ServerName: host}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

// renderTemplate is a basic template renderer for SMTP provider.
func renderTemplate(name string, data map[string]interface{}) string {
	switch name {
	case "verify_email":
		return fmt.Sprintf(`
			<h2>Welcome to Web3Sphere!</h2>
			<p>Your verification code is: <strong>%v</strong></p>
			<p>This code expires in %v minutes.</p>
		`, data["otp"], data["expiry_minutes"])
	case "reset_password":
		return fmt.Sprintf(`
			<h2>Password Reset</h2>
			<p>Your password reset code is: <strong>%v</strong></p>
			<p>This code expires in %v minutes.</p>
		`, data["otp"], data["expiry_minutes"])
	case "welcome":
		return fmt.Sprintf(`
			<h2>Welcome to Web3Sphere, %v!</h2>
			<p>Your account has been successfully verified. Start building the future of Web3!</p>
		`, data["name"])
	default:
		return fmt.Sprintf("Template: %s, Data: %v", name, data)
	}
}
