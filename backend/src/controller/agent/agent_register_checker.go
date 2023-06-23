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

func (c *Controller) RegisterCheckerV1(ctx context.Context, req *apiagent.RegisterCheckerV1Request) (*apiagent.Empty, error) {
	reSession, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{
			AgentUID: null.StringFrom(reSession.AgentUID),
			Type:     null.StringFrom(req.Checker.Type),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checkers: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgCheckers) == 0 {
		// Not exists yet, create new checker
		pgChecker := &model.PgCheckerV1{}

		pgChecker.UID, err = utils.UUID()
		if err != nil {
			c.log.Errorf("Error generating checker uid: %s", err)

			return nil, fmt.Errorf("internal server error")
		}
		pgChecker.AgentUID = reSession.AgentUID
		pgChecker.Name = req.Checker.Name
		pgChecker.Type = req.Checker.Type
		pgChecker.Version = req.Checker.Version
		pgChecker.Capabilities, err = c.mapAPICheckerV1ToPgCheckerV1Capabilities(req.Checker)
		if err != nil {
			c.log.Errorf("Error mapping checker: %s", err)

			return nil, fmt.Errorf("bad request")
		}
		pgChecker.CustomChecks = req.Checker.CustomChecks

		err = c.postgresService.AddChecker(
			ctx,
			pgChecker,
		)
		if err != nil {
			c.log.Errorf("Error adding new checker: %s", err)

			return nil, fmt.Errorf("internal server error")
		}

		reSystemEventPayload := &model.ReSystemEventV1CheckerAddedPayload{}
		reSystemEventPayload.CheckerUID = pgChecker.UID

		reSystemEvent := &model.ReSystemEventV1{}
		reSystemEvent.Type = model.ReSystemEventV1TypeCheckerAdded
		reSystemEvent.Payload = reSystemEventPayload

		err = c.cacheService.PublishSystemEvent(reSystemEvent)
		if err != nil {
			c.log.Errorf("Error publishing system event: %s", err)

			return nil, fmt.Errorf("internal server error")
		}
	} else {
		// Checker already registered, update it
		pgChecker := pgCheckers[0]

		pgChecker.Name = req.Checker.Name
		pgChecker.Type = req.Checker.Type
		pgChecker.Version = req.Checker.Version
		pgChecker.Capabilities, err = c.mapAPICheckerV1ToPgCheckerV1Capabilities(req.Checker)
		if err != nil {
			c.log.Errorf("Error mapping checker: %s", err)

			return nil, fmt.Errorf("bad request")
		}
		pgChecker.CustomChecks = req.Checker.CustomChecks

		err = c.postgresService.UpdateChecker(
			ctx,
			pgChecker.UID,
			pgChecker,
		)
		if err != nil {
			c.log.Errorf("Error updating checker: %s", err)

			return nil, fmt.Errorf("internal server error")
		}

		reSystemEventPayload := &model.ReSystemEventV1CheckerUpdatedPayload{}
		reSystemEventPayload.CheckerUID = pgChecker.UID

		reSystemEvent := &model.ReSystemEventV1{}
		reSystemEvent.Type = model.ReSystemEventV1TypeCheckerUpdated
		reSystemEvent.Payload = reSystemEventPayload

		err = c.cacheService.PublishSystemEvent(reSystemEvent)
		if err != nil {
			c.log.Errorf("Error publishing system event: %s", err)

			return nil, fmt.Errorf("internal server error")
		}

		return &apiagent.Empty{}, nil
	}

	return &apiagent.Empty{}, nil
}
