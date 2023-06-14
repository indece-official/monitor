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
)

func (c *Controller) reqV1SetConfigProperty(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	key, errResp := gousuchi.URLParamString(r, "key")
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1SetConfigPropertyJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgConfigProperty, err := c.mapAPISetConfigPropertyV1RequestBodyToPgConfigPropertyV1(requestBody, apipublic.ConfigPropertyV1Key(key))
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to config property: %s", err)
	}

	if model.PgConfigPropertyV1Protections[pgConfigProperty.Key] != model.PgConfigPropertyV1ProtectionPublic {
		return gousuchi.BadRequest(r, "Property does not have access mode 'public'")
	}

	err = c.postgresService.UpsertConfigProperty(r.Context(), pgConfigProperty)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error upserting config property in postgres: %s", err)
	}

	if pgConfigProperty.Key == model.PgConfigPropertyV1KeyConnectorHost {
		err = c.generateServerCert(r.Context())
		if err != nil {
			return gousuchi.InternalServerError(r, "Error generating new server certificates: %s", err)
		}
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Set config property '%s'", pgConfigProperty.Key)
}
