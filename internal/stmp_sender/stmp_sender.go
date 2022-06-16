package stmp_sender

import (
	"fmt"
	"net/smtp"

	"github.com/homework3/notification/internal/config"
	"github.com/rs/zerolog/log"
)

type MailSender struct {
	address string
	from    string
}

func NewSender(cfg *config.Smtp) *MailSender {
	return &MailSender{
		address: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		from:    "company@example.com",
	}
}

func (s *MailSender) SendMail(email, subject, text string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", s.from, email, subject, text)

	err := smtp.SendMail(s.address, nil, s.from, []string{email}, []byte(msg))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send email")

		return err
	}

	return nil
}
