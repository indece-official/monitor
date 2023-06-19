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
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) RegisterAgentV1(ctx context.Context, req *apiagent.RegisterAgentV1Request) (*apiagent.Empty, error) {
	reSession, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	pgAgents, err := c.postgresService.GetAgents(
		ctx,
		&postgres.GetAgentsFilter{
			AgentUID: null.StringFrom(reSession.AgentUID),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading agent: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgAgents) != 1 {
		c.log.Errorf("Own agent not found")

		return nil, fmt.Errorf("internal server error")
	}

	pgAgent := pgAgents[0]
	pgAgent.Version.Scan(req.Version)
	pgAgent.Type.Scan(req.Type)

	pgAgent.DatetimeRegistered.Scan(time.Now())

	err = c.postgresService.UpdateAgent(
		ctx,
		pgAgent.UID,
		pgAgent,
	)
	if err != nil {
		c.log.Errorf("Error updating agent: %s", err)

		return nil, fmt.Errorf("internal server error")
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

	reStatus.Status = model.ReAgentStatusV1StatusReady

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
