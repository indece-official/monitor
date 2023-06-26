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
	"net/http"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1GetMaintenance(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	maintenanceUID, errResp := gousuchi.URLParamString(r, "maintenanceUID")
	if errResp != nil {
		return errResp
	}

	pgMaintenances, err := c.postgresService.GetMaintenances(
		r.Context(),
		&postgres.GetMaintenancesFilter{
			MaintenanceUID: null.StringFrom(maintenanceUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading maintenances: %s", err)
	}

	if len(pgMaintenances) != 1 {
		return gousuchi.NotFound(r, "Maintenance not found")
	}

	pgMaintenance := pgMaintenances[0]

	respData, err := c.mapPgMaintenanceV1ToAPIGetMaintenanceV1ResponseBody(pgMaintenance, true)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded maintenance %s", pgMaintenance.UID)
}
