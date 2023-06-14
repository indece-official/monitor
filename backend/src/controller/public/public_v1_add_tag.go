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

func (c *Controller) reqV1AddTag(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	requestBody := &apipublic.V1AddTagJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	pgTag, err := c.mapAPIAddTagV1RequestBodyToPgTagV1(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping request to tag: %s", err)
	}

	pgTag.UID, err = utils.UUID()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error generating tag uid: %s", err)
	}

	err = c.postgresService.AddTag(r.Context(), pgTag)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error adding tag to postgres: %s", err)
	}

	respData, err := c.mapPgTagV1ToAPIAddTagV1ResponseBody(pgTag)
	if err != nil {
		return gousuchi.BadRequest(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Added tag '%s'", pgTag.UID)
}
