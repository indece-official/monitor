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

	"github.com/indece-official/monitor/backend/src/generated/model/apiagent"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) RegisterCheckV1(ctx context.Context, req *apiagent.RegisterCheckV1Request) (*apiagent.Empty, error) {
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

	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{
			AgentUID: null.StringFrom(pgAgent.UID),
			Type:     null.StringFrom(req.Check.CheckerType),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checkers: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgCheckers) == 0 {
		c.log.Warnf("Missing checker for check")

		return nil, fmt.Errorf("bad request")
	}

	pgChecker := pgCheckers[0]

	pgExistingChecks, err := c.postgresService.GetChecks(
		ctx,
		&postgres.GetChecksFilter{
			CheckerUID: null.StringFrom(pgChecker.UID),
			Type:       null.StringFrom(req.Check.Type),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checks: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgExistingChecks) == 1 {
		// Check already exists

		return &apiagent.Empty{}, nil
	}

	pgCheck := &model.PgCheckV1{}

	pgCheck.UID, err = utils.UUID()
	if err != nil {
		c.log.Errorf("Error generating check uid: %s", err)

		return nil, fmt.Errorf("internal server error")
	}
	pgCheck.CheckerUID = pgChecker.UID
	pgCheck.Name = req.Check.Name
	pgCheck.Type.Scan(req.Check.Type)
	pgCheck.Config = &model.PgCheckV1Config{}
	pgCheck.Config.Params = []*model.PgCheckV1Param{}

	for _, reqCheckParam := range req.Check.Params {
		pgCheckParam := &model.PgCheckV1Param{}

		pgCheckParam.Name = reqCheckParam.Name
		pgCheckParam.Value = reqCheckParam.Value

		pgCheck.Config.Params = append(pgCheck.Config.Params, pgCheckParam)
	}

	if req.Check.Schedule != "" {
		pgCheck.Schedule.Scan(req.Check.Schedule)
	} else {
		pgCheck.Schedule.Scan(nil)
	}

	if req.Check.Timeout != "" {
		pgCheck.Config.Timeout.Scan(req.Check.Timeout)
	} else {
		pgCheck.Config.Timeout.Scan(nil)
	}

	pgCheck.Custom = false

	if req.Check.Schedule != "" {
		pgCheck.Schedule.Scan(req.Check.Schedule)
	}

	err = c.postgresService.AddCheck(
		ctx,
		pgCheck,
	)
	if err != nil {
		c.log.Errorf("Error adding new checker: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	reSystemEventPayload := &model.ReSystemEventV1CheckAddedPayload{}
	reSystemEventPayload.CheckUID = pgCheck.UID

	reSystemEvent := &model.ReSystemEventV1{}
	reSystemEvent.Type = model.ReSystemEventV1TypeCheckAdded
	reSystemEvent.Payload = reSystemEventPayload

	err = c.cacheService.PublishSystemEvent(reSystemEvent)
	if err != nil {
		c.log.Errorf("Error publishing system event: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	return &apiagent.Empty{}, nil
}
