// indece Monitor
// Copyright (C) 2023 indece UG (haftungsbeschränkt)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License or any
// later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package public

import (
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgNotifierV1TypeToAPINotifierV1Type(pgNotifierType model.PgNotifierV1Type) (apipublic.NotifierV1Type, error) {
	switch pgNotifierType {
	case model.PgNotifierV1TypeEmailSmtp:
		return apipublic.EMAILSMTP, nil
	default:
		return "", fmt.Errorf("invalid notifier type: %s", pgNotifierType)
	}
}

func (c *Controller) mapAPINotifierV1TypeToPgNotifierV1Type(apiNotifierType apipublic.NotifierV1Type) (model.PgNotifierV1Type, error) {
	switch apiNotifierType {
	case apipublic.EMAILSMTP:
		return model.PgNotifierV1TypeEmailSmtp, nil
	default:
		return "", fmt.Errorf("invalid notifier type: %s", apiNotifierType)
	}
}

func (c *Controller) mapPgNotifierV1ConfigFilterToAPINotifierV1Filter(pgNotifierFilter *model.PgNotifierV1ConfigFilter) (*apipublic.NotifierV1Filter, error) {
	apiNotifierFilter := &apipublic.NotifierV1Filter{}

	apiNotifierFilter.TagUids = pgNotifierFilter.TagUIDs
	apiNotifierFilter.Critical = pgNotifierFilter.Critical
	apiNotifierFilter.Warning = pgNotifierFilter.Warning
	apiNotifierFilter.Unknown = pgNotifierFilter.Unknown
	apiNotifierFilter.Decline = pgNotifierFilter.Decline
	apiNotifierFilter.MinDuration = pgNotifierFilter.MinDuration.String()

	return apiNotifierFilter, nil
}

func (c *Controller) mapAPINotifierV1FilterToPgNotifierV1ConfigFilter(apiNotifierFilter *apipublic.NotifierV1Filter) (*model.PgNotifierV1ConfigFilter, error) {
	var err error

	pgNotifierFilter := &model.PgNotifierV1ConfigFilter{}

	pgNotifierFilter.TagUIDs = apiNotifierFilter.TagUids
	pgNotifierFilter.Critical = apiNotifierFilter.Critical
	pgNotifierFilter.Warning = apiNotifierFilter.Warning
	pgNotifierFilter.Unknown = apiNotifierFilter.Unknown
	pgNotifierFilter.Decline = apiNotifierFilter.Decline
	pgNotifierFilter.MinDuration, err = time.ParseDuration(apiNotifierFilter.MinDuration)
	if err != nil {
		return nil, fmt.Errorf("error parsing min_duration: %s", err)
	}

	return pgNotifierFilter, nil
}

func (c *Controller) mapPgNotifierV1ConfigParamsEmailSmtpToAPINotifierV1ConfigParamsEmailSmtp(pgNotifierConfigParamsEmailSmtp *model.PgNotifierV1ConfigParamsEmailSmtp) (*apipublic.NotifierV1ConfigParamsEmailSmtp, error) {
	apiNotifierConfigParamsEmailSmtp := &apipublic.NotifierV1ConfigParamsEmailSmtp{}

	apiNotifierConfigParamsEmailSmtp.Host = pgNotifierConfigParamsEmailSmtp.Host
	apiNotifierConfigParamsEmailSmtp.Port = pgNotifierConfigParamsEmailSmtp.Port
	apiNotifierConfigParamsEmailSmtp.User = pgNotifierConfigParamsEmailSmtp.User
	apiNotifierConfigParamsEmailSmtp.Password = pgNotifierConfigParamsEmailSmtp.Password
	apiNotifierConfigParamsEmailSmtp.From = pgNotifierConfigParamsEmailSmtp.From
	apiNotifierConfigParamsEmailSmtp.To = pgNotifierConfigParamsEmailSmtp.To

	return apiNotifierConfigParamsEmailSmtp, nil
}

func (c *Controller) mapAPINotifierV1ConfigParamsEmailSmtpToPgNotifierV1ConfigParamsEmailSmtp(apiNotifierConfigParamsEmailSmtp *apipublic.NotifierV1ConfigParamsEmailSmtp) (*model.PgNotifierV1ConfigParamsEmailSmtp, error) {
	pgNotifierConfigParamsEmailSmtp := &model.PgNotifierV1ConfigParamsEmailSmtp{}

	pgNotifierConfigParamsEmailSmtp.Host = apiNotifierConfigParamsEmailSmtp.Host
	pgNotifierConfigParamsEmailSmtp.Port = apiNotifierConfigParamsEmailSmtp.Port
	pgNotifierConfigParamsEmailSmtp.User = apiNotifierConfigParamsEmailSmtp.User
	pgNotifierConfigParamsEmailSmtp.Password = apiNotifierConfigParamsEmailSmtp.Password
	pgNotifierConfigParamsEmailSmtp.From = apiNotifierConfigParamsEmailSmtp.From
	pgNotifierConfigParamsEmailSmtp.To = apiNotifierConfigParamsEmailSmtp.To

	return pgNotifierConfigParamsEmailSmtp, nil
}

func (c *Controller) mapPgNotifierV1ConfigParamsToAPINotifierV1ConfigParams(pgNotifierConfigParams *model.PgNotifierV1ConfigParams) (*apipublic.NotifierV1ConfigParams, error) {
	var err error

	apiNotifierConfigParams := &apipublic.NotifierV1ConfigParams{}

	if pgNotifierConfigParams.EmailSmtp != nil {
		apiNotifierConfigParams.EmailSmtp, err = c.mapPgNotifierV1ConfigParamsEmailSmtpToAPINotifierV1ConfigParamsEmailSmtp(pgNotifierConfigParams.EmailSmtp)
		if err != nil {
			return nil, fmt.Errorf("error mapping email_smtp params: %s", err)
		}
	}

	return apiNotifierConfigParams, nil
}

func (c *Controller) mapAPINotifierV1ConfigParamsToPgNotifierV1ConfigParams(apiNotifierConfigParams *apipublic.NotifierV1ConfigParams) (*model.PgNotifierV1ConfigParams, error) {
	var err error

	pgNotifierConfigParams := &model.PgNotifierV1ConfigParams{}

	if apiNotifierConfigParams.EmailSmtp != nil {
		pgNotifierConfigParams.EmailSmtp, err = c.mapAPINotifierV1ConfigParamsEmailSmtpToPgNotifierV1ConfigParamsEmailSmtp(apiNotifierConfigParams.EmailSmtp)
		if err != nil {
			return nil, fmt.Errorf("error mapping email_smtp params: %s", err)
		}
	}

	return pgNotifierConfigParams, nil
}

func (c *Controller) mapPgNotifierV1ConfigToAPINotifierV1Config(pgNotifierConfig *model.PgNotifierV1Config) (*apipublic.NotifierV1Config, error) {
	var err error

	apiNotifierConfig := &apipublic.NotifierV1Config{}

	apiNotifierConfig.Filters = []apipublic.NotifierV1Filter{}
	for _, pgNotifierFilter := range pgNotifierConfig.Filters {
		apiNotifierFilter, err := c.mapPgNotifierV1ConfigFilterToAPINotifierV1Filter(pgNotifierFilter)
		if err != nil {
			return nil, fmt.Errorf("error mapping filter: %s", err)
		}

		apiNotifierConfig.Filters = append(apiNotifierConfig.Filters, *apiNotifierFilter)
	}

	apiNotifierConfigParams, err := c.mapPgNotifierV1ConfigParamsToAPINotifierV1ConfigParams(pgNotifierConfig.Params)
	if err != nil {
		return nil, fmt.Errorf("error mapping params: %s", err)
	}

	apiNotifierConfig.Params = *apiNotifierConfigParams

	return apiNotifierConfig, nil
}

func (c *Controller) mapAPINotifierV1ConfigToPgNotifierV1Config(apiNotifierConfig *apipublic.NotifierV1Config) (*model.PgNotifierV1Config, error) {
	var err error

	pgNotifierConfig := &model.PgNotifierV1Config{}

	pgNotifierConfig.Filters = []*model.PgNotifierV1ConfigFilter{}
	for _, apiNotifierFilter := range apiNotifierConfig.Filters {
		apiNotifierFilterCpy := apiNotifierFilter
		pgNotifierFilter, err := c.mapAPINotifierV1FilterToPgNotifierV1ConfigFilter(&apiNotifierFilterCpy)
		if err != nil {
			return nil, fmt.Errorf("error mapping filter: %s", err)
		}

		pgNotifierConfig.Filters = append(pgNotifierConfig.Filters, pgNotifierFilter)
	}

	pgNotifierConfig.Params, err = c.mapAPINotifierV1ConfigParamsToPgNotifierV1ConfigParams(&apiNotifierConfig.Params)
	if err != nil {
		return nil, fmt.Errorf("error mapping params: %s", err)
	}

	return pgNotifierConfig, nil
}

func (c *Controller) mapPgNotifierV1ToAPINotifierV1(pgNotifier *model.PgNotifierV1, addConfig bool) (*apipublic.NotifierV1, error) {
	var err error

	apiNotifier := &apipublic.NotifierV1{}

	apiNotifier.Uid = pgNotifier.UID
	apiNotifier.Name = pgNotifier.Name

	apiNotifier.Type, err = c.mapPgNotifierV1TypeToAPINotifierV1Type(pgNotifier.Type)
	if err != nil {
		return nil, err
	}

	if addConfig {
		apiNotifier.Config, err = c.mapPgNotifierV1ConfigToAPINotifierV1Config(pgNotifier.Config)
		if err != nil {
			return nil, err
		}
	}

	return apiNotifier, nil
}

func (c *Controller) mapPgNotifierV1ToAPIGetNotifiersV1ResponseBody(pgNotifiers []*model.PgNotifierV1, addConfig bool) (*apipublic.V1GetNotifiersJSONResponseBody, error) {
	resp := &apipublic.V1GetNotifiersJSONResponseBody{}

	resp.Notifiers = []apipublic.NotifierV1{}

	for _, pgNotifier := range pgNotifiers {
		apiNotifier, err := c.mapPgNotifierV1ToAPINotifierV1(pgNotifier, addConfig)
		if err != nil {
			return nil, err
		}

		resp.Notifiers = append(resp.Notifiers, *apiNotifier)
	}

	return resp, nil
}

func (c *Controller) mapPgNotifierV1ToAPIGetNotifierV1ResponseBody(pgNotifier *model.PgNotifierV1, addConfig bool) (*apipublic.V1GetNotifierJSONResponseBody, error) {
	resp := &apipublic.V1GetNotifierJSONResponseBody{}

	apiNotifier, err := c.mapPgNotifierV1ToAPINotifierV1(pgNotifier, addConfig)
	if err != nil {
		return nil, err
	}

	resp.Notifier = *apiNotifier

	return resp, nil
}

func (c *Controller) mapAPIAddNotifierV1RequestBodyToPgNotifierV1(requestBody *apipublic.V1AddNotifierJSONRequestBody) (*model.PgNotifierV1, error) {
	var err error

	pgNotifier := &model.PgNotifierV1{}

	pgNotifier.Name = requestBody.Name

	pgNotifier.Type, err = c.mapAPINotifierV1TypeToPgNotifierV1Type(requestBody.Type)
	if err != nil {
		return nil, fmt.Errorf("error mapping notifier type: %s", err)
	}

	pgNotifier.Config, err = c.mapAPINotifierV1ConfigToPgNotifierV1Config(&requestBody.Config)
	if err != nil {
		return nil, fmt.Errorf("error mapping notifier config: %s", err)
	}

	return pgNotifier, nil
}

func (c *Controller) mapAPIUpdateNotifierV1RequestBodyToPgNotifierV1(requestBody *apipublic.V1UpdateNotifierJSONRequestBody, oldPgNotifier *model.PgNotifierV1) (*model.PgNotifierV1, error) {
	var err error

	tmp := *oldPgNotifier
	pgNotifier := tmp

	pgNotifier.Name = requestBody.Name

	pgNotifier.Config, err = c.mapAPINotifierV1ConfigToPgNotifierV1Config(&requestBody.Config)
	if err != nil {
		return nil, fmt.Errorf("error mapping notifier config: %s", err)
	}

	return &pgNotifier, nil
}

func (c *Controller) mapPgNotifierV1ToAPIAddNotifierV1ResponseBody(pgNotifier *model.PgNotifierV1) (*apipublic.V1AddNotifierJSONResponseBody, error) {
	resp := &apipublic.V1AddNotifierJSONResponseBody{}

	resp.NotifierUid = pgNotifier.UID

	return resp, nil
}
