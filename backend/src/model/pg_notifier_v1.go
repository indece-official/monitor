package model

import (
	"fmt"
	"net/url"
	"time"

	"gopkg.in/guregu/null.v4"
)

type PgNotifierV1Type string

const (
	PgNotifierV1TypeEmailSmtp      PgNotifierV1Type = "email_smtp"
	PgNotifierV1TypeHttp           PgNotifierV1Type = "http"
	PgNotifierV1TypeMicrosoftTeams PgNotifierV1Type = "microsoft_teams"
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

type PgNotifierV1ConfigParamsHttpMethod string

const (
	PgNotifierV1ConfigParamsHttpMethodGet  PgNotifierV1ConfigParamsHttpMethod = "GET"
	PgNotifierV1ConfigParamsHttpMethodPost PgNotifierV1ConfigParamsHttpMethod = "POST"
	PgNotifierV1ConfigParamsHttpMethodPut  PgNotifierV1ConfigParamsHttpMethod = "PUT"
)

type PgNotifierV1ConfigParamsHttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PgNotifierV1ConfigParamsHttp struct {
	URL     string                                `json:"url"`
	Method  PgNotifierV1ConfigParamsHttpMethod    `json:"method"`
	Headers []*PgNotifierV1ConfigParamsHttpHeader `json:"headers"`
	Body    null.String                           `json:"body"`
}

func (p *PgNotifierV1ConfigParamsHttp) Validate() error {
	_, err := url.Parse(p.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %s", err)
	}

	// TODO

	return nil
}

type PgNotifierV1ConfigParamsMicrosoftTeams struct {
	WebhookURL string `json:"webhook_url"`
}

func (p *PgNotifierV1ConfigParamsMicrosoftTeams) Validate() error {
	_, err := url.Parse(p.WebhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook_url: %s", err)
	}

	return nil
}

type PgNotifierV1ConfigParams struct {
	EmailSmtp      *PgNotifierV1ConfigParamsEmailSmtp      `json:"email_smtp"`
	Http           *PgNotifierV1ConfigParamsHttp           `json:"http"`
	MicrosoftTeams *PgNotifierV1ConfigParamsMicrosoftTeams `json:"microsoft_teams"`
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
	case PgNotifierV1TypeHttp:
		if p.Http == nil {
			return fmt.Errorf("missing config params for http")
		}

		err := p.Http.Validate()
		if err != nil {
			return fmt.Errorf("error validating config params for http: %s", err)
		}
	case PgNotifierV1TypeMicrosoftTeams:
		if p.Http == nil {
			return fmt.Errorf("missing config params for microsoft_teams")
		}

		err := p.MicrosoftTeams.Validate()
		if err != nil {
			return fmt.Errorf("error validating config params for microsoft_teams: %s", err)
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
