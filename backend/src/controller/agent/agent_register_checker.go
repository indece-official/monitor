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

	pgCheckersByType := map[string]*model.PgCheckerV1{}
	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{
			AgentType: null.StringFrom(pgAgent.Type.String),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading checkers: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	for _, pgChecker := range pgCheckers {
		pgCheckersByType[pgChecker.Type] = pgChecker
	}

	if pgCheckersByType[req.Checker.Type] != nil {
		// Checker already registered
		return &apiagent.Empty{}, nil
	}

	pgChecker := &model.PgCheckerV1{}

	pgChecker.UID, err = utils.UUID()
	if err != nil {
		c.log.Errorf("Error generating checker uid: %s", err)

		return nil, fmt.Errorf("internal server error")
	}
	pgChecker.AgentType = pgAgent.Type.String
	pgChecker.Name = req.Checker.Name
	pgChecker.Type = req.Checker.Type
	pgChecker.Version = req.Checker.Version
	pgChecker.Capabilities = &model.PgCheckerV1Capabilities{}
	pgChecker.Capabilities.Params = []*model.PgCheckerV1Param{}
	for _, reqCheckerParam := range req.Checker.Params {
		pgCheckerParam := &model.PgCheckerV1Param{}

		pgCheckerParam.Name = reqCheckerParam.Name
		pgCheckerParam.Label = reqCheckerParam.Label
		if reqCheckerParam.Hint != "" {
			pgCheckerParam.Hint.Scan(reqCheckerParam.Hint)
		}
		pgCheckerParam.Required = reqCheckerParam.Required
		switch reqCheckerParam.Type {
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypeNumber:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypeNumber
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypeText:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypeText
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypePassword:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypePassword
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypeSelect:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypeSelect
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypeDuration:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypeDuration
		case apiagent.CheckerV1ParamType_CheckerV1ParamTypeBoolean:
			pgCheckerParam.Type = model.PgCheckerV1ParamTypeBoolean
		default:
			c.log.Warnf("Invalid checker param type %s", reqCheckerParam.Type)

			return nil, fmt.Errorf("bad request")
		}

		pgCheckerParam.Options = reqCheckerParam.Options

		pgChecker.Capabilities.Params = append(pgChecker.Capabilities.Params, pgCheckerParam)
	}

	pgChecker.Capabilities.Values = []*model.PgCheckerV1Value{}
	for _, reqCheckerValue := range req.Checker.Values {
		pgCheckerValue := &model.PgCheckerV1Value{}

		pgCheckerValue.Name = reqCheckerValue.Name
		pgCheckerValue.Label = reqCheckerValue.Name // TODO
		switch reqCheckerValue.Type {
		case apiagent.CheckerV1ValueType_CheckerV1ValueTypeNumber:
			pgCheckerValue.Type = model.ValueTypeNumber
		case apiagent.CheckerV1ValueType_CheckerV1ValueTypeText:
			pgCheckerValue.Type = model.ValueTypeText
		case apiagent.CheckerV1ValueType_CheckerV1ValueTypeDate:
			pgCheckerValue.Type = model.ValueTypeDate
		case apiagent.CheckerV1ValueType_CheckerV1ValueTypeDateTime:
			pgCheckerValue.Type = model.ValueTypeDatetime
		case apiagent.CheckerV1ValueType_CheckerV1ValueTypeDuration:
			pgCheckerValue.Type = model.ValueTypeDuration
		default:
			c.log.Warnf("Invalid checker value type %s", reqCheckerValue.Type)

			return nil, fmt.Errorf("bad request")
		}

		if reqCheckerValue.MinWarn != "" {
			pgCheckerValue.MinWarn.Scan(reqCheckerValue.MinWarn)
		}

		if reqCheckerValue.MinCrit != "" {
			pgCheckerValue.MinCrit.Scan(reqCheckerValue.MinCrit)
		}

		if reqCheckerValue.MaxWarn != "" {
			pgCheckerValue.MaxWarn.Scan(reqCheckerValue.MaxWarn)
		}

		if reqCheckerValue.MaxCrit != "" {
			pgCheckerValue.MaxCrit.Scan(reqCheckerValue.MaxCrit)
		}

		pgChecker.Capabilities.Values = append(pgChecker.Capabilities.Values, pgCheckerValue)
	}

	pgChecker.CustomChecks = req.Checker.CustomChecks

	if req.Checker.DefaultSchedule != "" {
		pgChecker.Capabilities.DefaultSchedule.Scan(req.Checker.DefaultSchedule)
	}

	err = c.postgresService.AddChecker(
		ctx,
		pgChecker,
	)
	if err != nil {
		c.log.Errorf("Error adding new checker: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	return &apiagent.Empty{}, nil
}
