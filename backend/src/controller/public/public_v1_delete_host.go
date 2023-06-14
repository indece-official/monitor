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
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1DeleteHost(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r, model.PgUserV1RoleAdmin)
	if errResp != nil {
		return errResp
	}

	hostUID, errResp := gousuchi.URLParamString(r, "hostUID")
	if errResp != nil {
		return errResp
	}

	oldPgHosts, err := c.postgresService.GetHosts(
		r.Context(),
		&postgres.GetHostsFilter{
			HostUID: null.StringFrom(hostUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loadign existing host: %s", err)
	}

	if len(oldPgHosts) != 1 {
		return gousuchi.NotFound(r, "Host not found")
	}

	err = c.postgresService.DeleteHost(r.Context(), hostUID)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error updating host in postgres: %s", err)
	}

	reSystemEventPayload := &model.ReSystemEventV1HostDeletedPayload{}
	reSystemEventPayload.HostUID = oldPgHosts[0].UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeHostDeleted
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error publishing system event: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithDetailedMessage("Deleted host '%s'", hostUID)
}
