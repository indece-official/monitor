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

	"github.com/indece-official/monitor/backend/src/generated/model/apiconnector"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) RegisterCheckV1(ctx context.Context, req *apiconnector.RegisterCheckV1Request) (*apiconnector.Empty, error) {
	reSession, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	pgConnectors, err := c.postgresService.GetConnectors(
		ctx,
		&postgres.GetConnectorsFilter{
			ConnectorUID: null.StringFrom(reSession.ConnectorUID),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading connector: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgConnectors) != 1 {
		c.log.Errorf("Own connector not found")

		return nil, fmt.Errorf("internal server error")
	}

	pgConnector := pgConnectors[0]

	pgCheckersByType := map[string]*model.PgCheckerV1{}
	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{
			ConnectorType: null.StringFrom(pgConnector.Type.String),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checkers: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	for _, pgChecker := range pgCheckers {
		pgCheckersByType[pgChecker.Type] = pgChecker
	}

	pgChecker := pgCheckersByType[req.Check.CheckerType]
	if pgChecker == nil {
		c.log.Warnf("Missing checker for check")

		return nil, fmt.Errorf("bad request")
	}

	pgExistingChecks, err := c.postgresService.GetChecks(
		ctx,
		&postgres.GetChecksFilter{
			CheckerUID: null.StringFrom(pgChecker.UID),
			HostUID:    null.StringFrom(pgConnector.HostUID),
			Type:       null.StringFrom(req.Check.Type),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checks: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgExistingChecks) == 1 {
		// Check already exists

		return &apiconnector.Empty{}, nil
	}

	pgCheck := &model.PgCheckV1{}

	pgCheck.UID, err = utils.UUID()
	if err != nil {
		c.log.Errorf("Error generating check uid: %s", err)

		return nil, fmt.Errorf("internal server error")
	}
	pgCheck.HostUID = pgConnector.HostUID
	pgCheck.CheckerUID = pgChecker.UID
	pgCheck.Name = req.Check.Name
	pgCheck.Type.Scan(req.Check.Type)
	// TODO: Schedule
	pgCheck.Config = &model.PgCheckV1Config{}
	pgCheck.Config.Params = []*model.PgCheckV1Param{}

	for _, reqCheckParam := range req.Check.Params {
		pgCheckParam := &model.PgCheckV1Param{}

		pgCheckParam.Name = reqCheckParam.Name
		pgCheckParam.Value = reqCheckParam.Value

		pgCheck.Config.Params = append(pgCheck.Config.Params, pgCheckParam)
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

	return &apiconnector.Empty{}, nil
}
