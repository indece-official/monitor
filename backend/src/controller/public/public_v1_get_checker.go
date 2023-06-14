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

func (c *Controller) reqV1GetChecker(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	checkerUID, errResp := gousuchi.URLParamString(r, "checkerUID")
	if errResp != nil {
		return errResp
	}

	pgCheckers, err := c.postgresService.GetCheckers(
		r.Context(),
		&postgres.GetCheckersFilter{
			CheckerUID: null.StringFrom(checkerUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading checkers: %s", err)
	}

	if len(pgCheckers) != 1 {
		return gousuchi.NotFound(r, "Checkers not found")
	}

	pgChecker := pgCheckers[0]

	respData, err := c.mapPgCheckerV1ToAPIGetCheckerV1ResponseBody(pgChecker)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded checker %s", pgChecker.UID)
}
