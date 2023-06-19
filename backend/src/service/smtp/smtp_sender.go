package smtp

import (
	"bytes"
	"sync"
	"time"

	"github.com/go-mail/mail"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/model"
)

type Sender struct {
	log         *logger.Log
	dialer      *mail.Dialer
	closer      mail.SendCloser
	config      *model.PgNotifierV1ConfigParamsEmailSmtp
	lastSend    *time.Time
	mutexCloser sync.Mutex
}

func (s *Sender) autoclose() {
	if s.lastSend == nil || time.Since(*s.lastSend) < 40*time.Second {
		return
	}

	s.Close()
}

func (s *Sender) Close() {
	s.mutexCloser.Lock()
	defer s.mutexCloser.Unlock()

	if s.closer == nil {
		return
	}

	err := s.closer.Close()
	if err != nil {
		s.log.Warnf("Can't close smtp connection: %s", err)
	}

	s.closer = nil
}

func (s *Sender) SendEmail(m *Email) error {
	var err error

	msg := mail.NewMessage()

	msg.SetHeader("From", s.config.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)

	if m.BodyPlain != "" {
		msg.SetBody("text/plain", m.BodyPlain)
	}

	if m.BodyHTML != "" {
		msg.SetBody("text/html", m.BodyHTML)
	}

	if m.Attachements != nil {
		for i := range m.Attachements {
			attachement := m.Attachements[i]
			reader := bytes.NewReader(attachement.Content)

			if attachement.Embedded {
				msg.EmbedReader(attachement.Filename, reader)
			} else {
				msg.AttachReader(attachement.Filename, reader)
			}
		}
	}

	s.mutexCloser.Lock()
	defer s.mutexCloser.Unlock()

	if s.closer == nil {
		closer, err := s.dialer.Dial()
		if err != nil {
			return err
		}

		now := time.Now()
		s.lastSend = &now
		s.closer = closer
	}

	err = mail.Send(s.closer, msg)
	if err != nil {
		if s.closer != nil {
			err = s.closer.Close()
			if err != nil {
				s.log.Errorf("Error lcosing connection after error: %s", err)
			}

			s.closer = nil
		}

		return err
	}

	return nil
}

func (s *Service) Open(config *model.PgNotifierV1ConfigParamsEmailSmtp) *Sender {
	dialer := mail.NewDialer(
		config.Host,
		config.Port,
		config.User,
		config.Password,
	)

	dialer.Timeout = 35 * time.Second
	dialer.RetryFailure = true

	sender := &Sender{
		dialer: dialer,
		config: config,
	}

	s.runningFuncs.Add(1)
	go func() {
		defer s.runningFuncs.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		stop, subStop := s.stopBroadcaster.Subscribe()
		defer subStop.Unsubscribe()

		for {
			select {
			case <-stop:
				sender.Close()
				return
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-ticker.C:
				sender.autoclose()
			}
		}
	}()

	return sender
}
