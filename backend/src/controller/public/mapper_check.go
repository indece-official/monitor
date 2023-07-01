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

func (c *Controller) mapPgCheckV1ToAPICheckV1(pgCheck *model.PgCheckV1, reCheckStatus *model.ReCheckStatusV1, addParams bool) (*apipublic.CheckV1, error) {
	var err error

	apiCheck := &apipublic.CheckV1{}

	apiCheck.Uid = pgCheck.UID
	apiCheck.Name = pgCheck.Name
	apiCheck.CheckerUid = pgCheck.CheckerUID
	apiCheck.Custom = pgCheck.Custom
	apiCheck.Schedule = pgCheck.Schedule.Ptr()
	apiCheck.Disabled = pgCheck.DatetimeDisabled.Valid

	if addParams {
		apiCheckParams := []apipublic.CheckV1Param{}

		for _, pgCheckParam := range pgCheck.Config.Params {
			apiCheckParam := apipublic.CheckV1Param{}

			apiCheckParam.Name = pgCheckParam.Name
			apiCheckParam.Value = pgCheckParam.Value

			apiCheckParams = append(apiCheckParams, apiCheckParam)
		}

		apiCheck.Params = &apiCheckParams
	}

	if reCheckStatus != nil {
		apiCheck.Status, err = c.mapReCheckStatusV1ToAPICheckStatusV1(reCheckStatus)
		if err != nil {
			return nil, fmt.Errorf("error mapping check status: %s", err)
		}
	}

	return apiCheck, nil
}

func (c *Controller) mapPgCheckV1ToAPIGetChecksV1ResponseBody(pgChecks []*model.PgCheckV1, reCheckStatuses map[string]*model.ReCheckStatusV1) (*apipublic.V1GetChecksJSONResponseBody, error) {
	resp := &apipublic.V1GetChecksJSONResponseBody{}

	resp.Checks = []apipublic.CheckV1{}

	for _, pgCheck := range pgChecks {
		apiCheck, err := c.mapPgCheckV1ToAPICheckV1(pgCheck, reCheckStatuses[pgCheck.UID], false)
		if err != nil {
			return nil, err
		}

		resp.Checks = append(resp.Checks, *apiCheck)
	}

	return resp, nil
}

func (c *Controller) mapPgCheckV1ToAPIGetHostChecksV1ResponseBody(pgChecks []*model.PgCheckV1, reCheckStatuses map[string]*model.ReCheckStatusV1) (*apipublic.V1GetHostChecksJSONResponseBody, error) {
	resp := &apipublic.V1GetHostChecksJSONResponseBody{}

	resp.Checks = []apipublic.CheckV1{}

	for _, pgCheck := range pgChecks {
		apiCheck, err := c.mapPgCheckV1ToAPICheckV1(pgCheck, reCheckStatuses[pgCheck.UID], false)
		if err != nil {
			return nil, err
		}

		resp.Checks = append(resp.Checks, *apiCheck)
	}

	return resp, nil
}

func (c *Controller) mapPgCheckV1ToAPIGetCheckV1ResponseBody(pgCheck *model.PgCheckV1, reCheckStatus *model.ReCheckStatusV1) (*apipublic.V1GetCheckJSONResponseBody, error) {
	resp := &apipublic.V1GetCheckJSONResponseBody{}

	apiCheck, err := c.mapPgCheckV1ToAPICheckV1(pgCheck, reCheckStatus, true)
	if err != nil {
		return nil, err
	}

	resp.Check = *apiCheck

	return resp, nil
}

func (c *Controller) mapAPIAddCheckV1RequestBodyToPgCheckV1(requestBody *apipublic.V1AddCheckJSONRequestBody) (*model.PgCheckV1, error) {
	pgCheck := &model.PgCheckV1{}

	pgCheck.Name = requestBody.Name
	pgCheck.CheckerUID = requestBody.CheckerUid
	pgCheck.Config = &model.PgCheckV1Config{}
	pgCheck.Config.Params = []*model.PgCheckV1Param{}

	for _, reqParam := range requestBody.Params {
		pgCheckParam := &model.PgCheckV1Param{}

		pgCheckParam.Name = reqParam.Name
		pgCheckParam.Value = reqParam.Value

		pgCheck.Config.Params = append(pgCheck.Config.Params, pgCheckParam)
	}

	if requestBody.Schedule != nil {
		pgCheck.Schedule.Scan(*requestBody.Schedule)
	}
	pgCheck.Custom = true

	return pgCheck, nil
}

func (c *Controller) mapAPIUpdateCheckV1RequestBodyToPgCheckV1(requestBody *apipublic.V1UpdateCheckJSONRequestBody, oldPgCheck *model.PgCheckV1) (*model.PgCheckV1, error) {
	tmp := *oldPgCheck
	pgCheck := tmp

	pgCheck.Name = requestBody.Name
	pgCheck.Config = &model.PgCheckV1Config{}
	pgCheck.Config.Params = []*model.PgCheckV1Param{}

	for _, reqParam := range requestBody.Params {
		pgCheckParam := &model.PgCheckV1Param{}

		pgCheckParam.Name = reqParam.Name
		pgCheckParam.Value = reqParam.Value

		pgCheck.Config.Params = append(pgCheck.Config.Params, pgCheckParam)
	}

	if requestBody.Schedule != nil {
		pgCheck.Schedule.Scan(*requestBody.Schedule)
	} else {
		pgCheck.Schedule.Scan(nil)
	}

	return &pgCheck, nil
}
func (c *Controller) mapPgCheckV1ToAPIAddCheckV1ResponseBody(pgCheck *model.PgCheckV1) (*apipublic.V1AddCheckJSONResponseBody, error) {
	resp := &apipublic.V1AddCheckJSONResponseBody{}

	resp.CheckUid = pgCheck.UID

	return resp, nil
}
