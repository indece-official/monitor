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
	"time"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgConnectorV1ToAPIConnectorV1(pgConnector *model.PgConnectorV1, reStatus *model.ReConnectorStatusV1) (*apipublic.ConnectorV1, error) {
	apiConnector := &apipublic.ConnectorV1{}

	apiConnector.Uid = pgConnector.UID
	apiConnector.HostUid = pgConnector.HostUID
	apiConnector.Type = pgConnector.Type.Ptr()
	apiConnector.Version = pgConnector.Version.Ptr()

	if !pgConnector.DatetimeRegistered.Valid || (reStatus != nil && reStatus.Status == model.ReConnectorStatusV1StatusUnregistered) {
		apiConnector.Status = apipublic.ConnectorV1StatusUNREGISTERED
	} else if reStatus == nil || reStatus.Status == model.ReConnectorStatusV1StatusReady {
		apiConnector.Status = apipublic.ConnectorV1StatusREADY
	} else if reStatus != nil && reStatus.Status == model.ReConnectorStatusV1StatusError {
		apiConnector.Status = apipublic.ConnectorV1StatusERROR
	} else {
		apiConnector.Status = apipublic.ConnectorV1StatusUNKNOWN
	}

	if reStatus != nil {
		apiConnector.Connected = reStatus.DatetimeLastPing.Valid && time.Since(reStatus.DatetimeLastPing.Time) < 30*time.Second
		apiConnector.LastPing = reStatus.DatetimeLastPing.Ptr()
		apiConnector.Error = reStatus.Error.Ptr()
	} else {
		apiConnector.Connected = false
	}

	return apiConnector, nil
}

func (c *Controller) mapAPIAddConnectorV1RequestBodyToPgConnectorV1(requestBody *apipublic.V1AddConnectorJSONRequestBody) (*model.PgConnectorV1, error) {
	pgConnector := &model.PgConnectorV1{}

	pgConnector.HostUID = requestBody.HostUid

	return pgConnector, nil
}

func (c *Controller) mapPgConnectorV1ToAPIGetConnectorsV1ResponseBody(pgConnectors []*model.PgConnectorV1, reStatuses map[string]*model.ReConnectorStatusV1) (*apipublic.V1GetConnectorsJSONResponseBody, error) {
	resp := &apipublic.V1GetConnectorsJSONResponseBody{}

	resp.Connectors = []apipublic.ConnectorV1{}

	for _, pgConnector := range pgConnectors {
		apiConnector, err := c.mapPgConnectorV1ToAPIConnectorV1(pgConnector, reStatuses[pgConnector.UID])
		if err != nil {
			return nil, err
		}

		resp.Connectors = append(resp.Connectors, *apiConnector)
	}

	return resp, nil
}

func (c *Controller) mapPgConnectorV1ToAPIGetConnectorV1ResponseBody(pgConnector *model.PgConnectorV1, reStatus *model.ReConnectorStatusV1) (*apipublic.V1GetConnectorJSONResponseBody, error) {
	resp := &apipublic.V1GetConnectorJSONResponseBody{}

	apiConnector, err := c.mapPgConnectorV1ToAPIConnectorV1(pgConnector, reStatus)
	if err != nil {
		return nil, err
	}

	resp.Connector = *apiConnector

	return resp, nil
}

func (c *Controller) mapPgConnectorV1ToAPIAddConnectorV1ResponseBody(pgConnector *model.PgConnectorV1, configFile string) (*apipublic.V1AddConnectorJSONResponseBody, error) {
	resp := &apipublic.V1AddConnectorJSONResponseBody{}

	resp.ConnectorUid = pgConnector.UID
	resp.ConfigFile = configFile

	return resp, nil
}
