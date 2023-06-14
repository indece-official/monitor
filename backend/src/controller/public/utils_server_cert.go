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
	"fmt"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/cert"
	"github.com/indece-official/monitor/backend/src/service/postgres"
)

func (c *Controller) generateServerCert(ctx context.Context) error {
	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		ctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyTLSCaCrt,
				model.PgConfigPropertyV1KeyTLSCaKey,
				model.PgConfigPropertyV1KeyConnectorHost,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error loading config properties: %s", err)
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyConnectorHost] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyConnectorHost].Value == "" {
		return fmt.Errorf("connector host not configured")
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value == "" ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey].Value == "" {
		return fmt.Errorf("root ca not configured")
	}

	serverPEM, err := c.certService.GenerateServerCert(
		pgConfigProperties[model.PgConfigPropertyV1KeyConnectorHost].Value,
		&cert.PEMCert{
			Crt: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value),
			Key: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey].Value),
		},
	)
	if err != nil {
		return fmt.Errorf("error generating server cert: %s", err)
	}

	err = c.postgresService.UpsertConfigProperty(
		ctx,
		&model.PgConfigPropertyV1{
			Key:   model.PgConfigPropertyV1KeyTLSServerCrt,
			Value: string(serverPEM.Crt),
		},
	)
	if err != nil {
		return fmt.Errorf("error adding tls_server_crt config property: %s", err)
	}

	err = c.postgresService.UpsertConfigProperty(
		ctx,
		&model.PgConfigPropertyV1{
			Key:   model.PgConfigPropertyV1KeyTLSServerKey,
			Value: string(serverPEM.Key),
		},
	)
	if err != nil {
		return fmt.Errorf("error adding tls_server_key config property: %s", err)
	}

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeCertsUpdated
	reSystemEvent.Payload = &model.ReSystemEventV1CertsUpdatedPayload{}

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return fmt.Errorf("error publishing system event: %s", err)
	}

	return nil
}
