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

package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/generated/model/apiconnector"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) PingV1(ctx context.Context, req *apiconnector.PingV1Request) (*apiconnector.Empty, error) {
	reSession, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	reStatus, err := c.cacheService.GetConnectorStatus(reSession.ConnectorUID)
	if err != nil {
		c.log.Errorf("Error getting connector status: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if reStatus == nil {
		reStatus = &model.ReConnectorStatusV1{}

		reStatus.ConnectorUID = reSession.ConnectorUID
	}

	reStatus.DatetimeLastPing.Scan(time.Now())

	err = c.cacheService.SetConnectorStatus(
		reSession.ConnectorUID,
		reStatus,
	)
	if err != nil {
		c.log.Errorf("Error setting connector status: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	return &apiconnector.Empty{}, nil
}
