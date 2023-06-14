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
)

func (c *Controller) reqV1Logout(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	cookie, err := r.Cookie("dsusersession")
	if err != nil {
		return gousuchi.JSON(r, map[string]interface{}{}).
			WithDetailedMessage("Logged out user - error loading session cookie: %s", err)
	}

	if cookie == nil || cookie.Value == "" {
		return gousuchi.JSON(r, map[string]interface{}{}).
			WithDetailedMessage("Logged out user - empty session cookie")
	}

	sessionKey := cookie.Value

	reSession, err := c.cacheService.GetUserSession(sessionKey)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading user session: %s", err)
	}

	if reSession == nil {
		return gousuchi.JSON(r, map[string]interface{}{}).
			WithDetailedMessage("Logged out user - session does not exist")
	}

	err = c.cacheService.DeleteUserSession(sessionKey)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error deleting user session: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Logged out user %s", reSession.UserUID)
}
