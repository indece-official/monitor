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

func (c *Controller) reqV1GetConnectors(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	pgConnectors, err := c.postgresService.GetConnectors(
		r.Context(),
		&postgres.GetConnectorsFilter{},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading connectors: %s", err)
	}

	reStatuses := map[string]*model.ReConnectorStatusV1{}
	for _, pgConnector := range pgConnectors {
		reStatus, err := c.cacheService.GetConnectorStatus(pgConnector.UID)
		if err != nil {
			return gousuchi.InternalServerError(r, "Error loading connector status: %s", err)
		}

		if reStatus != nil {
			reStatuses[pgConnector.UID] = reStatus
		}
	}

	respData, err := c.mapPgConnectorV1ToAPIGetConnectorsV1ResponseBody(pgConnectors, reStatuses)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d connectors", len(pgConnectors))
}
