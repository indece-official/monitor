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

package model

import "gopkg.in/guregu/null.v4"

type ReConnectorStatusV1Status string

const (
	ReConnectorStatusV1StatusUnknown      ReConnectorStatusV1Status = "unknown"
	ReConnectorStatusV1StatusUnregistered ReConnectorStatusV1Status = "unregistered"
	ReConnectorStatusV1StatusReady        ReConnectorStatusV1Status = "ready"
	ReConnectorStatusV1StatusError        ReConnectorStatusV1Status = "error"
)

type ReConnectorStatusV1 struct {
	ConnectorUID     string                    `json:"connector_uid"`
	Status           ReConnectorStatusV1Status `json:"status"`
	Error            null.String               `json:"error"`
	DatetimeLastPing null.Time                 `json:"datetime_last_ping"`
}
