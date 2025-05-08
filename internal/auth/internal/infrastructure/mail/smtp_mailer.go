package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/istiak-004/myFolio-microservices/auth/internal/domain/ports"
)

type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPMailer(host string, port int, username, password, from string) ports.Mailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *SMTPMailer) SendVerificationEmail(email, name, verificationURL string) error {
	// Load email template
	tmpl, err := template.New("verification").Parse(`
    <html>
    <body>
        <h1>Welcome, {{.Name}}!</h1>
        <p>Please verify your email by clicking the link below:</p>
        <a href="{{.URL}}">Verify Email</a>
        <p>This link will expire in 24 hours.</p>
    </body>
    </html>
    `)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  verificationURL,
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	to := []string{email}
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: Verify Your Email\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s",
		email, m.from, body.String()))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", m.host, m.port),
		auth,
		m.from,
		to,
		msg,
	)
}
