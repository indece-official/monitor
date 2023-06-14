package model

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

type PgNotifierV1Type string

const (
	PgNotifierV1TypeEmailSmtp PgNotifierV1Type = "email_smtp"
)

type PgNotifierV1ConfigFilter struct {
	TagUIDs     []string      `json:"tag_uids"`
	Critical    bool          `json:"critical"`
	Warning     bool          `json:"warning"`
	Unknown     bool          `json:"unknown"`
	Decline     bool          `json:"decline"`
	MinDuration time.Duration `json:"min_duration"`
}

type PgNotifierV1ConfigParamsEmailSmtp struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

func (p *PgNotifierV1ConfigParamsEmailSmtp) Validate() error {
	if p.Port <= 0 {
		return fmt.Errorf("invalid port: %d", p.Port)
	}

	// TODO

	return nil
}

type PgNotifierV1ConfigParams struct {
	EmailSmtp *PgNotifierV1ConfigParamsEmailSmtp `json:"email_smtp"`
}

func (p *PgNotifierV1ConfigParams) Validate(notifierType PgNotifierV1Type) error {
	switch notifierType {
	case PgNotifierV1TypeEmailSmtp:
		if p.EmailSmtp == nil {
			return fmt.Errorf("missing config params for email_smtp")
		}

		err := p.EmailSmtp.Validate()
		if err != nil {
			return fmt.Errorf("error validating config params for email_smtp: %s", err)
		}
	}

	return nil
}

type PgNotifierV1Config struct {
	Filters []*PgNotifierV1ConfigFilter `json:"filters"`
	Params  *PgNotifierV1ConfigParams   `json:"params"`
}

func (p *PgNotifierV1Config) Validate(notifierType PgNotifierV1Type) error {
	if p.Params == nil {
		return fmt.Errorf("missing params")
	}

	err := p.Params.Validate(notifierType)
	if err != nil {
		return fmt.Errorf("error validating config: %s", err)
	}

	return nil
}

type PgNotifierV1 struct {
	UID              string
	Name             string
	Type             PgNotifierV1Type
	Config           *PgNotifierV1Config
	DatetimeDisabled null.Time
}

func (p *PgNotifierV1) Validate() error {
	if p.Config == nil {
		return fmt.Errorf("missing config")
	}

	err := p.Config.Validate(p.Type)
	if err != nil {
		return fmt.Errorf("error validating config: %s", err)
	}

	return nil
}
