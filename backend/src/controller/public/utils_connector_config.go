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
	"context"
	"encoding/base64"
	"fmt"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
)

func (c *Controller) generateConnectorConfigFile(ctx context.Context, pgConnector *model.PgConnectorV1) (string, error) {
	pgConfigProperies, err := c.postgresService.GetConfigProperties(
		ctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyTLSCaCrt,
				model.PgConfigPropertyV1KeyConnectorHost,
				model.PgConfigPropertyV1KeyConnectorPort,
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error loading config properties: %s", err)
	}

	if pgConfigProperies[model.PgConfigPropertyV1KeyTLSCaCrt] == nil ||
		pgConfigProperies[model.PgConfigPropertyV1KeyTLSCaCrt].Value == "" ||
		pgConfigProperies[model.PgConfigPropertyV1KeyConnectorHost] == nil ||
		pgConfigProperies[model.PgConfigPropertyV1KeyConnectorHost].Value == "" ||
		pgConfigProperies[model.PgConfigPropertyV1KeyConnectorPort] == nil ||
		pgConfigProperies[model.PgConfigPropertyV1KeyConnectorPort].Value == "" {
		return "", fmt.Errorf("missing config properties")
	}

	data := fmt.Sprintf("server_host=%s\n", pgConfigProperies[model.PgConfigPropertyV1KeyConnectorHost].Value)
	data += fmt.Sprintf("server_port=%s\n", pgConfigProperies[model.PgConfigPropertyV1KeyConnectorPort].Value)

	caCrtBase64 := base64.StdEncoding.EncodeToString([]byte(pgConfigProperies[model.PgConfigPropertyV1KeyTLSCaCrt].Value))
	clientCrtBase64 := base64.StdEncoding.EncodeToString([]byte(pgConnector.TLSClientCrt))
	clientKeyBase64 := base64.StdEncoding.EncodeToString([]byte(pgConnector.TLSClientKey))

	data += fmt.Sprintf("ca_crt=%s\n", caCrtBase64)
	data += fmt.Sprintf("client_crt=%s\n", clientCrtBase64)
	data += fmt.Sprintf("client_key=%s\n", clientKeyBase64)

	return data, nil
}
