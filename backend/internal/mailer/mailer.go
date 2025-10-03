package mailer

import (
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

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
		mail.WithTimeout(time.Duration(cfg.SMTP.Timeout)*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &Mailer{
		client: client,
		from:   cfg.SMTP.From,
	}, nil
}

func (m *Mailer) SendHTML(to, subject string, templateFile string, data any) error {
	msg := mail.NewMsg()
	if err := msg.From(m.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(subject)
	tmpl, err := template.ParseFS(templateFS, fmt.Sprintf("templates/%s", templateFile))
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	if err := msg.SetBodyHTMLTemplate(tmpl, data); err != nil {
		return fmt.Errorf("failed to set email body: %w", err)
	}

	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
