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

func (c *Controller) reqV1GetOwnUser(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	reSession, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	pgUsers, err := c.postgresService.GetUsers(
		r.Context(),
		&postgres.GetUsersFilter{
			UserUID: null.StringFrom(reSession.UserUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading users: %s", err)
	}

	if len(pgUsers) != 1 {
		return gousuchi.InternalServerError(r, "Couldn't find own user")
	}

	pgUser := pgUsers[0]

	respData, err := c.mapPgUserV1ToAPIGetOwnUserV1ResponseBody(pgUser)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded own user %s", pgUser.UID)
}
