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

package initializer

import (
	"context"
	"fmt"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
)

func (c *Controller) generateCerts(ctx context.Context) error {
	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		ctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyAgentHost,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error loading config properties: %s", err)
	}

	caPEM, err := c.certService.GenerateCACert()
	if err != nil {
		return fmt.Errorf("error generating ca cert: %s", err)
	}

	err = c.postgresService.UpsertConfigProperty(
		ctx,
		&model.PgConfigPropertyV1{
			Key:   model.PgConfigPropertyV1KeyTLSCaCrt,
			Value: string(caPEM.Crt),
		},
	)
	if err != nil {
		return fmt.Errorf("error adding tls_ca_crt config property: %s", err)
	}

	err = c.postgresService.UpsertConfigProperty(
		ctx,
		&model.PgConfigPropertyV1{
			Key:   model.PgConfigPropertyV1KeyTLSCaKey,
			Value: string(caPEM.Key),
		},
	)
	if err != nil {
		return fmt.Errorf("error adding tls_ca_key config property: %s", err)
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyAgentHost] != nil &&
		pgConfigProperties[model.PgConfigPropertyV1KeyAgentHost].Value != "" {
		serverPEM, err := c.certService.GenerateServerCert(
			pgConfigProperties[model.PgConfigPropertyV1KeyAgentHost].Value,
			caPEM,
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
