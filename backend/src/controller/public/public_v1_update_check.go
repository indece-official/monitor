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

func (c *Controller) reqV1UpdateCheck(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	checkUID, errResp := gousuchi.URLParamString(r, "checkUID")
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1UpdateCheckJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	oldPgChecks, err := c.postgresService.GetChecks(
		r.Context(),
		&postgres.GetChecksFilter{
			CheckUID: null.StringFrom(checkUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loadign existing check: %s", err)
	}

	if len(oldPgChecks) != 1 {
		return gousuchi.NotFound(r, "Check not found")
	}

	oldPgCheck := oldPgChecks[0]
	if !oldPgCheck.Custom {
		return gousuchi.NotFound(r, "Check is not a custom check")
	}

	pgCheck, err := c.mapAPIUpdateCheckV1RequestBodyToPgCheckV1(requestBody, oldPgCheck)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to config property: %s", err)
	}

	err = c.postgresService.UpdateCheck(r.Context(), pgCheck.UID, pgCheck)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error updating check in postgres: %s", err)
	}

	reSystemEventPayload := &model.ReSystemEventV1CheckUpdatedPayload{}
	reSystemEventPayload.CheckUID = oldPgChecks[0].UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeCheckUpdated
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error publishing system event: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Updated check '%s'", pgCheck.UID)
}
