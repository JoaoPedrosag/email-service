package mailer

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
    APIKey string
    From   string
}

func New() *Mailer {
    return &Mailer{
        APIKey: os.Getenv("SENDGRID_API_KEY"),
        From:   os.Getenv("EMAIL_FROM"),
    }
}

func (m *Mailer) Send(to, subject, body string) error {
    from := mail.NewEmail("Email Service", m.From)
    toEmail := mail.NewEmail("", to)
    msg := mail.NewSingleEmail(from, subject, toEmail, body, body)

    client := sendgrid.NewSendClient(m.APIKey)
    res, err := client.Send(msg)
    if err != nil {
        return err
    }

    if res.StatusCode >= 400 {
        return fmt.Errorf("SendGrid error: status %d, body: %s", res.StatusCode, res.Body)
    }

    return nil
}
