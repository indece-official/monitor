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
	"github.com/indece-official/monitor/backend/src/utils"
)

func (c *Controller) reqV1AddMaintenance(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddMaintenanceJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgMaintenance, err := c.mapAPIAddMaintenanceV1RequestBodyToPgMaintenanceV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to maintenance: %s", err)
	}

	pgMaintenance.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating maintenance uid: %s", err)
	}

	err = c.postgresService.AddMaintenance(r.Context(), pgMaintenance)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding maintenance to postgres: %s", err)
	}

	respData, err := c.mapPgMaintenanceV1ToAPIAddMaintenanceV1ResponseBody(pgMaintenance)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added maintenance '%s'", pgMaintenance.UID)
}
