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

package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/generated/model/apiagent"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) PingV1(ctx context.Context, req *apiagent.PingV1Request) (*apiagent.Empty, error) {
	reSession, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	reStatus, err := c.cacheService.GetAgentStatus(reSession.AgentUID)
	if err != nil {
		c.log.Errorf("Error getting agent status: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if reStatus == nil {
		reStatus = &model.ReAgentStatusV1{}

		reStatus.AgentUID = reSession.AgentUID
	}

	reStatus.DatetimeLastPing.Scan(time.Now())

	err = c.cacheService.SetAgentStatus(
		reSession.AgentUID,
		reStatus,
	)
	if err != nil {
		c.log.Errorf("Error setting agent status: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	return &apiagent.Empty{}, nil
}
