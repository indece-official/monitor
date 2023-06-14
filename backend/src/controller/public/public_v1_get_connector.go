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
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqV1GetConnector(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	connectorUID, errResp := gousuchi.URLParamString(r, "connectorUID")
	if errResp != nil {
		return errResp
	}

	pgConnectors, err := c.postgresService.GetConnectors(
		r.Context(),
		&postgres.GetConnectorsFilter{
			ConnectorUID: null.StringFrom(connectorUID),
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading connectors: %s", err)
	}

	if len(pgConnectors) != 1 {
		return gousuchi.NotFound(r, "Connectors not found")
	}

	pgConnector := pgConnectors[0]

	reStatus, err := c.cacheService.GetConnectorStatus(pgConnector.UID)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading connector status: %s", err)
	}

	respData, err := c.mapPgConnectorV1ToAPIGetConnectorV1ResponseBody(pgConnector, reStatus)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded connector %s", pgConnector.UID)
}