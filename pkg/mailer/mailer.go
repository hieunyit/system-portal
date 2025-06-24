package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"system-portal/internal/shared/config"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type SMTPMailer struct {
	cfg config.SMTPConfig
}

func NewSMTPMailer(cfg config.SMTPConfig) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	if m.cfg.Host == "" || m.cfg.Port == 0 {
		return fmt.Errorf("smtp not configured")
	}
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	msg := []byte("From: " + m.cfg.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if m.cfg.TLS {
		c, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer c.Close()
		if err = c.StartTLS(&tls.Config{ServerName: m.cfg.Host}); err != nil {
			return err
		}
		if err = c.Auth(auth); err != nil {
			return err
		}
		if err = c.Mail(m.cfg.From); err != nil {
			return err
		}
		if err = c.Rcpt(to); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err = w.Write(msg); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}
		return c.Quit()
	}
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}
