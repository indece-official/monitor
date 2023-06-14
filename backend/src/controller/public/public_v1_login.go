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
	"time"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1Login(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	requestBody := &apipublic.V1LoginJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgUsers, err := c.postgresService.GetUsers(
		r.Context(),
		&postgres.GetUsersFilter{
			Username: null.StringFrom(requestBody.Username),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading user from postgres: %s", err)
	}

	if len(pgUsers) != 1 {
		return gousuchi.BadRequest(r, "User not found")
	}

	pgUser := pgUsers[0]

	if pgUser.DatetimeLocked.Valid {
		return gousuchi.BadRequest(r, "User is locked")
	}

	if pgUser.Source != model.PgUserV1SourceLocal {
		return gousuchi.BadRequest(r, "User is not from local source")
	}

	if !pgUser.PasswordHash.Valid {
		return gousuchi.BadRequest(r, "User has no password set")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(pgUser.PasswordHash.String),
		[]byte(requestBody.Password),
	)
	if err != nil {
		return gousuchi.BadRequest(r, "Wrong password: %s", err)
	}

	sessionKey, err := utils.RandString(128)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating session key: %s", err)
	}

	reSession := &model.ReUserSessionV1{}

	reSession.SessionKey = sessionKey
	reSession.UserUID = pgUser.UID
	reSession.Roles = pgUser.LocalRoles
	reSession.DatetimeCreated = time.Now()

	err = c.cacheService.SetUserSession(
		sessionKey,
		reSession,
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error storing session in redis: %s", err)
	}

	expiration := time.Now().Add(24 * 356 * time.Hour)
	cookie := &http.Cookie{
		Name:     "dsusersession",
		Value:    sessionKey,
		Expires:  expiration,
		Path:     "/",
		Secure:   *cookieSecure,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Logged in user %s", pgUser.UID)
}
