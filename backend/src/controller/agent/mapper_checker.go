package agent

import (
	"fmt"

	"github.com/indece-official/monitor/backend/src/generated/model/apiagent"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapAPICheckerV1ToPgCheckerV1Capabilities(apiChecker *apiagent.CheckerV1) (*model.PgCheckerV1Capabilities, error) {
	pgCapabilities := &model.PgCheckerV1Capabilities{}

	pgCapabilities.Params = []*model.PgCheckerV1Param{}
	for _, reqCheckerParam := range apiChecker.Params {
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

		pgCapabilities.Params = append(pgCapabilities.Params, pgCheckerParam)
	}

	pgCapabilities.Values = []*model.PgCheckerV1Value{}
	for _, reqCheckerValue := range apiChecker.Values {
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

		pgCapabilities.Values = append(pgCapabilities.Values, pgCheckerValue)
	}

	if apiChecker.DefaultSchedule != "" {
		pgCapabilities.DefaultSchedule.Scan(apiChecker.DefaultSchedule)
	} else {
		pgCapabilities.DefaultSchedule.Scan(nil)
	}

	if apiChecker.DefaultTimeout != "" {
		pgCapabilities.DefaultTimeout.Scan(apiChecker.DefaultTimeout)
	} else {
		pgCapabilities.DefaultTimeout.Scan(nil)
	}

	return pgCapabilities, nil
}
