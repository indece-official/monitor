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

func (c *Controller) reqV1GetCheck(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	checkUID, errResp := gousuchi.URLParamString(r, "checkUID")
	if errResp != nil {
		return errResp
	}

	pgChecks, err := c.postgresService.GetChecks(
		r.Context(),
		&postgres.GetChecksFilter{
			CheckUID: null.StringFrom(checkUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading checks: %s", err)
	}

	if len(pgChecks) != 1 {
		return gousuchi.NotFound(r, "Checks not found")
	}

	pgCheck := pgChecks[0]

	respData, err := c.mapPgCheckV1ToAPIGetCheckV1ResponseBody(pgCheck)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded check %s", pgCheck.UID)
}
