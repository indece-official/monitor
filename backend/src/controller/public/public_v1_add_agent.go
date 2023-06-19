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

func (c *Controller) reqV1AddAgent(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddAgentJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgAgent, err := c.mapAPIAddAgentV1RequestBodyToPgAgentV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to agent: %s", err)
	}

	pgAgent.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating agent uid: %s", err)
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
		pgAgent.UID,
		&cert.PEMCert{
			Crt: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value),
			Key: []byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaKey].Value),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating client cert: %s", err)
	}

	pgAgent.Certs = &model.PgAgentV1Certs{}
	pgAgent.Certs.Certs = []*model.PgAgentV1Cert{}

	pgAgentCert := &model.PgAgentV1Cert{}
	pgAgentCert.TLSClientCrt = string(clientPEM.Crt)
	pgAgentCert.TLSClientKey = string(clientPEM.Key)
	pgAgentCert.CreateAt = clientPEM.CreatedAt.Time
	pgAgentCert.ValidUntil = clientPEM.ValidUntil.Time

	pgAgent.Certs.Certs = append(pgAgent.Certs.Certs, pgAgentCert)

	err = c.postgresService.AddAgent(r.Context(), pgAgent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding agent to postgres: %s", err)
	}

	configFile, err := c.generateAgentConfigFile(r.Context(), pgAgent, clientPEM)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating agent config file: %s", err)
	}

	respData, err := c.mapPgAgentV1ToAPIAddAgentV1ResponseBody(pgAgent, configFile)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added agent '%s'", pgAgent.UID)
}
