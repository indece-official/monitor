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

func (c *Controller) reqV1GetAgents(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	pgAgents, err := c.postgresService.GetAgents(
		r.Context(),
		&postgres.GetAgentsFilter{},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading agents: %s", err)
	}

	reStatuses := map[string]*model.ReAgentStatusV1{}
	for _, pgAgent := range pgAgents {
		reStatus, err := c.cacheService.GetAgentStatus(pgAgent.UID)
		if err != nil {
			return gousuchi.InternalServerError(r, "Error loading agent status: %s", err)
		}

		if reStatus != nil {
			reStatuses[pgAgent.UID] = reStatus
		}
	}

	respData, err := c.mapPgAgentV1ToAPIGetAgentsV1ResponseBody(pgAgents, reStatuses)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d agents", len(pgAgents))
}
