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
)

func (c *Controller) reqV1GetMaintenances(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	active, errResp := gousuchi.OptionalQueryParamBool(r, "active")
	if errResp != nil {
		return errResp
	}

	from, errResp := gousuchi.OptionalQueryParamInt64(r, "from")
	if errResp != nil {
		return errResp
	}

	if !from.Valid {
		from.Scan(0)
	}

	size, errResp := gousuchi.OptionalQueryParamInt64(r, "size")
	if errResp != nil {
		return errResp
	}

	if !size.Valid {
		size.Scan(30)
	}

	pgMaintenances, err := c.postgresService.GetMaintenances(
		r.Context(),
		&postgres.GetMaintenancesFilter{
			NowActive: active.Valid && active.Bool,
			From:      from,
			Size:      size,
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading maintenances: %s", err)
	}

	respData, err := c.mapPgMaintenanceV1ToAPIGetMaintenancesV1ResponseBody(pgMaintenances, false)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d maintenances", len(pgMaintenances))
}
