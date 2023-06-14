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
	"encoding/json"
	"net/http"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/cert"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
)

func (c *Controller) reqV1AddConnector(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddConnectorJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgConnector, err := c.mapAPIAddConnectorV1RequestBodyToPgConnectorV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to connector: %s", err)
	}

	pgConnector.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating connector uid: %s", err)
	}

	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		r.Context(),
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyTLSCaCrt,
				model.PgConfigPropertyV1KeyTLSCaKey,
			},
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading config properties: %s", err)
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value == "" ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey].Value == "" {
		return gousuchi.InternalServerError(r, "No clients cert available")
	}

	clientPEM, err := c.certService.GenerateClientCert(
		pgConnector.UID,
		&cert.PEMCert{
			Crt: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value),
			Key: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey].Value),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating client cert: %s", err)
	}

	pgConnector.TLSClientCrt = string(clientPEM.Crt)
	pgConnector.TLSClientKey = string(clientPEM.Key)

	err = c.postgresService.AddConnector(r.Context(), pgConnector)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding connector to postgres: %s", err)
	}

	configFile, err := c.generateConnectorConfigFile(r.Context(), pgConnector)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating connector config file: %s", err)
	}

	respData, err := c.mapPgConnectorV1ToAPIAddConnectorV1ResponseBody(pgConnector, configFile)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added connector '%s'", pgConnector.UID)
}
