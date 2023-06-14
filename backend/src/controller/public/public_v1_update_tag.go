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

func (c *Controller) reqV1UpdateTag(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	tagUID, errResp := gousuchi.URLParamString(r, "tagUID")
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1UpdateTagJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	oldPgTags, err := c.postgresService.GetTags(
		r.Context(),
		&postgres.GetTagsFilter{
			TagUID: null.StringFrom(tagUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loadign existing tag: %s", err)
	}

	if len(oldPgTags) != 1 {
		return gousuchi.NotFound(r, "Tag not found")
	}

	pgTag, err := c.mapAPIUpdateTagV1RequestBodyToPgTagV1(requestBody, oldPgTags[0])
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to config property: %s", err)
	}

	err = c.postgresService.UpdateTag(r.Context(), pgTag.UID, pgTag)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error updating tag in postgres: %s", err)
	}

	reSystemEventPayload := &model.ReSystemEventV1TagUpdatedPayload{}
	reSystemEventPayload.TagUID = oldPgTags[0].UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeTagUpdated
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error publishing system event: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Updated tag '%s'", pgTag.UID)
}
