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
	"github.com/indece-official/monitor/backend/src/utils"
)

func (c *Controller) reqV1AddHost(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddHostJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgHost, err := c.mapAPIAddHostV1RequestBodyToPgHostV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to host: %s", err)
	}

	pgHost.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating host uid: %s", err)
	}

	err = c.postgresService.AddHost(r.Context(), pgHost)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding host to postgres: %s", err)
	}

	reSystemEventPayload := &model.ReSystemEventV1HostAddedPayload{}
	reSystemEventPayload.HostUID = pgHost.UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeHostAdded
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error publishing system event: %s", err)
	}

	respData, err := c.mapPgHostV1ToAPIAddHostV1ResponseBody(pgHost)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added host '%s'", pgHost.UID)
}
