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
	case model.PgNotifierV1TypeHttp:
		return apipublic.HTTP, nil
	case model.PgNotifierV1TypeMicrosoftTeams:
		return apipublic.MICROSOFTTEAMS, nil
	default:
		return "", fmt.Errorf("invalid notifier type: %s", pgNotifierType)
	}
}

func (c *Controller) mapAPINotifierV1TypeToPgNotifierV1Type(apiNotifierType apipublic.NotifierV1Type) (model.PgNotifierV1Type, error) {
	switch apiNotifierType {
	case apipublic.EMAILSMTP:
		return model.PgNotifierV1TypeEmailSmtp, nil
	case apipublic.HTTP:
		return model.PgNotifierV1TypeHttp, nil
	case apipublic.MICROSOFTTEAMS:
		return model.PgNotifierV1TypeMicrosoftTeams, nil
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

func (c *Controller) mapPgNotifierV1ConfigParamsHttpMethodToAPINotifierV1ConfigParamsHttpMethod(pgHttpMethod model.PgNotifierV1ConfigParamsHttpMethod) (apipublic.NotifierV1ConfigParamsHttpMethod, error) {
	switch pgHttpMethod {
	case model.PgNotifierV1ConfigParamsHttpMethodGet:
		return apipublic.GET, nil
	case model.PgNotifierV1ConfigParamsHttpMethodPost:
		return apipublic.POST, nil
	case model.PgNotifierV1ConfigParamsHttpMethodPut:
		return apipublic.PUT, nil
	default:
		return "", fmt.Errorf("invalid http method: %s", pgHttpMethod)
	}
}

func (c *Controller) mapAPINotifierV1ConfigParamsHttpMethodToPgNotifierV1ConfigParamsHttpMethod(apiHttpMethod apipublic.NotifierV1ConfigParamsHttpMethod) (model.PgNotifierV1ConfigParamsHttpMethod, error) {
	switch apiHttpMethod {
	case apipublic.GET:
		return model.PgNotifierV1ConfigParamsHttpMethodGet, nil
	case apipublic.POST:
		return model.PgNotifierV1ConfigParamsHttpMethodPost, nil
	case apipublic.PUT:
		return model.PgNotifierV1ConfigParamsHttpMethodPut, nil
	default:
		return "", fmt.Errorf("invalid http method: %s", apiHttpMethod)
	}
}

func (c *Controller) mapPgNotifierV1ConfigParamsHttpHeaderToAPINotifierV1ConfigParamsHttpHeader(pgNotifierConfigParamsHttpHeader *model.PgNotifierV1ConfigParamsHttpHeader) (*apipublic.NotifierV1ConfigParamsHttpHeader, error) {
	apiNotifierConfigParamsHttpHeader := &apipublic.NotifierV1ConfigParamsHttpHeader{}

	apiNotifierConfigParamsHttpHeader.Name = pgNotifierConfigParamsHttpHeader.Name
	apiNotifierConfigParamsHttpHeader.Value = pgNotifierConfigParamsHttpHeader.Value

	return apiNotifierConfigParamsHttpHeader, nil
}

func (c *Controller) mapAPINotifierV1ConfigParamsHttpHeaderToPgNotifierV1ConfigParamsHttpHeader(apiNotifierConfigParamsHttpHeader *apipublic.NotifierV1ConfigParamsHttpHeader) (*model.PgNotifierV1ConfigParamsHttpHeader, error) {
	pgNotifierConfigParamsHttpHeader := &model.PgNotifierV1ConfigParamsHttpHeader{}

	pgNotifierConfigParamsHttpHeader.Name = apiNotifierConfigParamsHttpHeader.Name
	pgNotifierConfigParamsHttpHeader.Value = apiNotifierConfigParamsHttpHeader.Value

	return pgNotifierConfigParamsHttpHeader, nil
}

func (c *Controller) mapPgNotifierV1ConfigParamsHttpToAPINotifierV1ConfigParamsHttp(pgNotifierConfigParamsHttp *model.PgNotifierV1ConfigParamsHttp) (*apipublic.NotifierV1ConfigParamsHttp, error) {
	var err error

	apiNotifierConfigParamsHttp := &apipublic.NotifierV1ConfigParamsHttp{}

	apiNotifierConfigParamsHttp.Url = pgNotifierConfigParamsHttp.URL
	apiNotifierConfigParamsHttp.Method, err = c.mapPgNotifierV1ConfigParamsHttpMethodToAPINotifierV1ConfigParamsHttpMethod(pgNotifierConfigParamsHttp.Method)
	if err != nil {
		return nil, fmt.Errorf("error mapping http method: %s", err)
	}

	apiNotifierConfigParamsHttp.Headers = []apipublic.NotifierV1ConfigParamsHttpHeader{}
	for _, pgHeader := range pgNotifierConfigParamsHttp.Headers {
		apiHeader, err := c.mapPgNotifierV1ConfigParamsHttpHeaderToAPINotifierV1ConfigParamsHttpHeader(pgHeader)
		if err != nil {
			return nil, fmt.Errorf("error mapping http header: %s", err)
		}

		apiNotifierConfigParamsHttp.Headers = append(apiNotifierConfigParamsHttp.Headers, *apiHeader)
	}

	apiNotifierConfigParamsHttp.Body = pgNotifierConfigParamsHttp.Body.Ptr()

	return apiNotifierConfigParamsHttp, nil
}

func (c *Controller) mapAPINotifierV1ConfigParamsHttpToPgNotifierV1ConfigParamsHttp(apiNotifierConfigParamsHttp *apipublic.NotifierV1ConfigParamsHttp) (*model.PgNotifierV1ConfigParamsHttp, error) {
	var err error

	pgNotifierConfigParamsHttp := &model.PgNotifierV1ConfigParamsHttp{}

	pgNotifierConfigParamsHttp.URL = apiNotifierConfigParamsHttp.Url
	pgNotifierConfigParamsHttp.Method, err = c.mapAPINotifierV1ConfigParamsHttpMethodToPgNotifierV1ConfigParamsHttpMethod(apiNotifierConfigParamsHttp.Method)
	if err != nil {
		return nil, fmt.Errorf("error mapping http method: %s", err)
	}

	pgNotifierConfigParamsHttp.Headers = []*model.PgNotifierV1ConfigParamsHttpHeader{}
	for _, apiHeader := range apiNotifierConfigParamsHttp.Headers {
		apiHeaderCpy := apiHeader

		pgHeader, err := c.mapAPINotifierV1ConfigParamsHttpHeaderToPgNotifierV1ConfigParamsHttpHeader(&apiHeaderCpy)
		if err != nil {
			return nil, fmt.Errorf("error mapping http header: %s", err)
		}

		pgNotifierConfigParamsHttp.Headers = append(pgNotifierConfigParamsHttp.Headers, pgHeader)
	}

	if apiNotifierConfigParamsHttp.Body != nil && *apiNotifierConfigParamsHttp.Body != "" {
		pgNotifierConfigParamsHttp.Body.Scan(*apiNotifierConfigParamsHttp.Body)
	}

	return pgNotifierConfigParamsHttp, nil
}

func (c *Controller) mapPgNotifierV1ConfigParamsMicrosoftTeamsToAPINotifierV1ConfigParamsMicrosoftTeams(pgNotifierConfigParamsMicrosoftTeams *model.PgNotifierV1ConfigParamsMicrosoftTeams) (*apipublic.NotifierV1ConfigParamsMicrosoftTeams, error) {
	apiNotifierConfigParamsMicrosoftTeams := &apipublic.NotifierV1ConfigParamsMicrosoftTeams{}

	apiNotifierConfigParamsMicrosoftTeams.WebhookUrl = pgNotifierConfigParamsMicrosoftTeams.WebhookURL

	return apiNotifierConfigParamsMicrosoftTeams, nil
}

func (c *Controller) mapAPINotifierV1ConfigParamsMicrosoftTeamsToPgNotifierV1ConfigParamsMicrosoftTeams(apiNotifierConfigParamsMicrosoftTeams *apipublic.NotifierV1ConfigParamsMicrosoftTeams) (*model.PgNotifierV1ConfigParamsMicrosoftTeams, error) {
	pgNotifierConfigParamsMicrosoftTeams := &model.PgNotifierV1ConfigParamsMicrosoftTeams{}

	pgNotifierConfigParamsMicrosoftTeams.WebhookURL = apiNotifierConfigParamsMicrosoftTeams.WebhookUrl

	return pgNotifierConfigParamsMicrosoftTeams, nil
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

	if pgNotifierConfigParams.Http != nil {
		apiNotifierConfigParams.Http, err = c.mapPgNotifierV1ConfigParamsHttpToAPINotifierV1ConfigParamsHttp(pgNotifierConfigParams.Http)
		if err != nil {
			return nil, fmt.Errorf("error mapping http params: %s", err)
		}
	}

	if pgNotifierConfigParams.MicrosoftTeams != nil {
		apiNotifierConfigParams.MicrosoftTeams, err = c.mapPgNotifierV1ConfigParamsMicrosoftTeamsToAPINotifierV1ConfigParamsMicrosoftTeams(pgNotifierConfigParams.MicrosoftTeams)
		if err != nil {
			return nil, fmt.Errorf("error mapping microsoft_teams params: %s", err)
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

	if apiNotifierConfigParams.Http != nil {
		pgNotifierConfigParams.Http, err = c.mapAPINotifierV1ConfigParamsHttpToPgNotifierV1ConfigParamsHttp(apiNotifierConfigParams.Http)
		if err != nil {
			return nil, fmt.Errorf("error mapping http params: %s", err)
		}
	}

	if apiNotifierConfigParams.MicrosoftTeams != nil {
		pgNotifierConfigParams.MicrosoftTeams, err = c.mapAPINotifierV1ConfigParamsMicrosoftTeamsToPgNotifierV1ConfigParamsMicrosoftTeams(apiNotifierConfigParams.MicrosoftTeams)
		if err != nil {
			return nil, fmt.Errorf("error mapping microsoft_teams params: %s", err)
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
