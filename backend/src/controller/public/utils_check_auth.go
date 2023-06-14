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
	"time"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) checkAuth(r *http.Request, roles ...model.PgUserV1Role) (*model.ReUserSessionV1, gousuchi.IResponse) {
	setupTokenHeader := r.Header.Get("X-Setup-Token")
	if setupTokenHeader != "" {
		setupToken, err := c.cacheService.GetSetupToken()
		if err != nil {
			return nil, gousuchi.InternalServerError(r, "Error loading setup token: %s", err)
		}

		if setupToken.Valid && setupToken.String == setupTokenHeader {
			reSession := &model.ReUserSessionV1{}

			reSession.UserUID = "setup"
			reSession.SessionKey = ""
			reSession.Roles = []model.PgUserV1Role{
				model.PgUserV1RoleAdmin,
				model.PgUserV1RoleSetup,
			}
			reSession.DatetimeCreated = time.Now()

			return reSession, nil
		}
	}

	cookie, err := r.Cookie("dsusersession")
	if err != nil {
		return nil, gousuchi.Unauthorized(r, "Error loading session cookie: %s", err)
	}

	if cookie == nil || cookie.Value == "" {
		return nil, gousuchi.Unauthorized(r, "Empty session cookie")
	}

	sessionKey := cookie.Value

	reSession, err := c.cacheService.GetUserSession(sessionKey)
	if err != nil {
		return nil, gousuchi.InternalServerError(r, "Error loading user session: %s", err)
	}

	if reSession == nil {
		return nil, gousuchi.Unauthorized(r, "Session does not exist")
	}

	for _, role := range roles {
		found := false

		for _, sessionRole := range reSession.Roles {
			if sessionRole == role {
				found = true
				break
			}
		}

		if !found {
			return nil, gousuchi.Forbidden(r, "Missing role %s", role)
		}
	}

	return reSession, nil
}
