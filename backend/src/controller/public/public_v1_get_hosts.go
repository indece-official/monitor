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
)

func (c *Controller) reqV1GetHosts(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	pgHosts, err := c.postgresService.GetHosts(
		r.Context(),
		&postgres.GetHostsFilter{},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading hosts: %s", err)
	}

	reCheckStatuses, err := c.cacheService.GetAllCheckStatuses()
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading check status: %s", err)
	}

	reCheckStatusesByHost := map[string][]*model.ReCheckStatusV1{}
	for _, reCheckStatus := range reCheckStatuses {
		if reCheckStatusesByHost[reCheckStatus.HostUID] == nil {
			reCheckStatusesByHost[reCheckStatus.HostUID] = []*model.ReCheckStatusV1{}
		}

		reCheckStatusesByHost[reCheckStatus.HostUID] = append(reCheckStatusesByHost[reCheckStatus.HostUID], reCheckStatus)
	}

	respData, err := c.mapPgHostV1ToAPIGetHostsV1ResponseBody(pgHosts, reCheckStatusesByHost)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d hosts", len(pgHosts))
}
