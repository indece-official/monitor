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
	"fmt"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgCheckerV1ParamTypeToAPICheckerV1ParamType(pgCheckerParamType model.PgCheckerV1ParamType) (apipublic.CheckerV1ParamType, error) {
	switch pgCheckerParamType {
	case model.PgCheckerV1ParamTypeText:
		return apipublic.TEXT, nil
	case model.PgCheckerV1ParamTypePassword:
		return apipublic.PASSWORD, nil
	case model.PgCheckerV1ParamTypeNumber:
		return apipublic.NUMBER, nil
	case model.PgCheckerV1ParamTypeSelect:
		return apipublic.SELECT, nil
	case model.PgCheckerV1ParamTypeDuration:
		return apipublic.DURATION, nil
	case model.PgCheckerV1ParamTypeBoolean:
		return apipublic.BOOLEAN, nil
	default:
		return "", fmt.Errorf("invalid checker param type: %s", pgCheckerParamType)
	}
}

func (c *Controller) mapPgCheckerV1ParamToAPICheckerV1Param(pgCheckerParam *model.PgCheckerV1Param) (*apipublic.CheckerV1Param, error) {
	var err error

	apiCheckerParam := &apipublic.CheckerV1Param{}

	apiCheckerParam.Name = pgCheckerParam.Name
	apiCheckerParam.Label = pgCheckerParam.Label
	apiCheckerParam.Hint = pgCheckerParam.Hint.Ptr()
	apiCheckerParam.Required = pgCheckerParam.Required
	apiCheckerParam.Type, err = c.mapPgCheckerV1ParamTypeToAPICheckerV1ParamType(pgCheckerParam.Type)
	if err != nil {
		return nil, fmt.Errorf("error mapping param type: %s", err)
	}

	apiCheckerParam.Options = &pgCheckerParam.Options

	return apiCheckerParam, nil
}

func (c *Controller) mapPgCheckerV1ToAPICheckerV1(pgChecker *model.PgCheckerV1) (*apipublic.CheckerV1, error) {
	apiChecker := &apipublic.CheckerV1{}

	apiChecker.Uid = pgChecker.UID
	apiChecker.Name = pgChecker.Name
	apiChecker.Type = pgChecker.Type
	apiChecker.AgentType = pgChecker.AgentType
	apiChecker.Version = pgChecker.Version
	apiChecker.Capabilities = apipublic.CheckerV1Capabilities{}
	apiChecker.Capabilities.Params = []apipublic.CheckerV1Param{}
	apiChecker.Capabilities.DefaultSchedule = pgChecker.Capabilities.DefaultSchedule.Ptr()

	for _, pgCheckerParam := range pgChecker.Capabilities.Params {
		apiCheckerParam, err := c.mapPgCheckerV1ParamToAPICheckerV1Param(pgCheckerParam)
		if err != nil {
			return nil, err
		}

		apiChecker.Capabilities.Params = append(apiChecker.Capabilities.Params, *apiCheckerParam)
	}

	apiChecker.CustomChecks = pgChecker.CustomChecks

	return apiChecker, nil
}

func (c *Controller) mapPgCheckerV1ToAPIGetCheckersV1ResponseBody(pgCheckers []*model.PgCheckerV1) (*apipublic.V1GetCheckersJSONResponseBody, error) {
	resp := &apipublic.V1GetCheckersJSONResponseBody{}

	resp.Checkers = []apipublic.CheckerV1{}

	for _, pgChecker := range pgCheckers {
		apiChecker, err := c.mapPgCheckerV1ToAPICheckerV1(pgChecker)
		if err != nil {
			return nil, err
		}

		resp.Checkers = append(resp.Checkers, *apiChecker)
	}

	return resp, nil
}

func (c *Controller) mapPgCheckerV1ToAPIGetCheckerV1ResponseBody(pgChecker *model.PgCheckerV1) (*apipublic.V1GetCheckerJSONResponseBody, error) {
	resp := &apipublic.V1GetCheckerJSONResponseBody{}

	apiChecker, err := c.mapPgCheckerV1ToAPICheckerV1(pgChecker)
	if err != nil {
		return nil, err
	}

	resp.Checker = *apiChecker

	return resp, nil
}
