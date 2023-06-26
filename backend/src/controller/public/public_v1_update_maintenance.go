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
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1UpdateMaintenance(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	maintenanceUID, errResp := gousuchi.URLParamString(r, "maintenanceUID")
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1UpdateMaintenanceJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	oldPgMaintenances, err := c.postgresService.GetMaintenances(
		r.Context(),
		&postgres.GetMaintenancesFilter{
			MaintenanceUID: null.StringFrom(maintenanceUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loadign existing maintenance: %s", err)
	}

	if len(oldPgMaintenances) != 1 {
		return gousuchi.NotFound(r, "Maintenance not found")
	}

	pgMaintenance, err := c.mapAPIUpdateMaintenanceV1RequestBodyToPgMaintenanceV1(requestBody, oldPgMaintenances[0])
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to config property: %s", err)
	}

	err = c.postgresService.UpdateMaintenance(r.Context(), pgMaintenance.UID, pgMaintenance)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error updating maintenance in postgres: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Updated maintenance '%s'", pgMaintenance.UID)
}
