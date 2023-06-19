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
	"github.com/indece-official/monitor/backend/src/utils"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1AddCheck(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddCheckJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgCheck, err := c.mapAPIAddCheckV1RequestBodyToPgCheckV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to check: %s", err)
	}

	pgCheck.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating check uid: %s", err)
	}

	pgCheckers, err := c.postgresService.GetCheckers(
		r.Context(),
		&postgres.GetCheckersFilter{
			CheckerUID: null.StringFrom(pgCheck.CheckerUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading checker from postgres: %s", err)
	}

	if len(pgCheckers) != 1 {
		return gousuchi.BadRequest(r, "No matching checker found")
	}

	pgChecker := pgCheckers[0]

	if !pgChecker.CustomChecks {
		return gousuchi.BadRequest(r, "No custom checks allowed for this checker")
	}

	err = c.postgresService.AddCheck(r.Context(), pgCheck)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding check to postgres: %s", err)
	}

	reSystemEventPayload := &model.ReSystemEventV1CheckAddedPayload{}
	reSystemEventPayload.CheckUID = pgCheck.UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeCheckAdded
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error publishing system event: %s", err)
	}

	respData, err := c.mapPgCheckV1ToAPIAddCheckV1ResponseBody(pgCheck)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added check '%s'", pgCheck.UID)
}
