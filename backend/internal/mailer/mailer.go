package mailer

import (
	"fmt"
	"html/template"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/wneessen/go-mail"
)

type Mailer struct {
	client *mail.Client
	from   string
}

func New(cfg config.Config) (*Mailer, error) {
	client, err := mail.NewClient(
		cfg.SMTP.Host,
		mail.WithPort(cfg.SMTP.Port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(cfg.SMTP.Username),
		mail.WithPassword(cfg.SMTP.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &Mailer{
		client: client,
		from:   cfg.SMTP.From,
	}, nil
}

func (m *Mailer) SendHTML(to, subject string, tmpl *template.Template, data interface{}) error {
	msg := mail.NewMsg()
	if err := msg.From(m.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(subject)
	if err := msg.SetBodyHTMLTemplate(tmpl, data); err != nil {
		return fmt.Errorf("failed to set email body: %w", err)
	}

	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
