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

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgConfigPropertyV1KeyToAPIConfigPropertyV1Key(pgConfigPropertyKey model.PgConfigPropertyV1Key) (apipublic.ConfigPropertyV1Key, error) {
	switch pgConfigPropertyKey {
	case model.PgConfigPropertyV1KeySetupFinished:
		return apipublic.SETUPFINISHED, nil
	case model.PgConfigPropertyV1KeyTLSCaCrt:
		return apipublic.TLSCACRT, nil
	case model.PgConfigPropertyV1KeyTLSServerCrt:
		return apipublic.TLSSERVERCRT, nil
	case model.PgConfigPropertyV1KeyAgentHost:
		return apipublic.AGENTHOST, nil
	case model.PgConfigPropertyV1KeyAgentPort:
		return apipublic.AGENTPORT, nil
	case model.PgConfigPropertyV1KeyHistoryMaxAge:
		return apipublic.HISTORYMAXAGE, nil
	default:
		return "", fmt.Errorf("invalid config property key: %s", pgConfigPropertyKey)
	}
}

func (c *Controller) mapAPIConfigPropertyV1KeyToPgConfigPropertyV1Key(apiConfigPropertyKey apipublic.ConfigPropertyV1Key) (model.PgConfigPropertyV1Key, error) {
	switch apiConfigPropertyKey {
	case apipublic.AGENTHOST:
		return model.PgConfigPropertyV1KeyAgentHost, nil
	case apipublic.AGENTPORT:
		return model.PgConfigPropertyV1KeyAgentPort, nil
	case apipublic.HISTORYMAXAGE:
		return model.PgConfigPropertyV1KeyHistoryMaxAge, nil
	default:
		return "", fmt.Errorf("invalid config property key: %s", apiConfigPropertyKey)
	}
}

func (c *Controller) mapAPISetConfigPropertyV1RequestBodyToPgConfigPropertyV1(requestBody *apipublic.V1SetConfigPropertyJSONRequestBody, key apipublic.ConfigPropertyV1Key) (*model.PgConfigPropertyV1, error) {
	var err error

	pgConfigProperty := &model.PgConfigPropertyV1{}

	pgConfigProperty.Key, err = c.mapAPIConfigPropertyV1KeyToPgConfigPropertyV1Key(key)
	if err != nil {
		return nil, fmt.Errorf("error mapping key: %s", err)
	}

	pgConfigProperty.Value = requestBody.Value

	return pgConfigProperty, nil
}

func (c *Controller) mapPgConfigPropertyV1ToAPIConfigPropertyV1(pgConfigProperty *model.PgConfigPropertyV1) (*apipublic.ConfigPropertyV1, error) {
	var err error

	apiConfigProperty := &apipublic.ConfigPropertyV1{}

	apiConfigProperty.Key, err = c.mapPgConfigPropertyV1KeyToAPIConfigPropertyV1Key(pgConfigProperty.Key)
	if err != nil {
		return nil, err
	}

	apiConfigProperty.Value = pgConfigProperty.Value
	apiConfigProperty.Editable = model.PgConfigPropertyV1Protections[pgConfigProperty.Key] == model.PgConfigPropertyV1ProtectionPublic

	return apiConfigProperty, nil
}

func (c *Controller) mapPgConfigPropertyV1ToAPIGetConfigV1ResponseBody(pgConfigProperties []*model.PgConfigPropertyV1) (*apipublic.V1GetConfigJSONResponseBody, error) {
	resp := &apipublic.V1GetConfigJSONResponseBody{}

	resp.Properties = []apipublic.ConfigPropertyV1{}

	for _, pgConfigProperty := range pgConfigProperties {
		apiConfigProperty, err := c.mapPgConfigPropertyV1ToAPIConfigPropertyV1(pgConfigProperty)
		if err != nil {
			return nil, err
		}

		resp.Properties = append(resp.Properties, *apiConfigProperty)
	}

	return resp, nil
}
