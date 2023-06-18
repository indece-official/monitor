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
	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapReHostStatusV1ToAPIHostV1Status(reHostStatus *model.ReHostStatusV1) (*apipublic.HostV1Status, error) {
	apiHostStatus := &apipublic.HostV1Status{}

	apiHostStatus.CountCritical = 0
	apiHostStatus.CountWarning = 0
	apiHostStatus.CountOk = 0
	apiHostStatus.CountUnknown = 0

	if reHostStatus != nil {
		for _, checkStatus := range reHostStatus.Checks {
			switch checkStatus.Status {
			case model.PgCheckStatusV1StatusCrit:
				apiHostStatus.CountCritical++
			case model.PgCheckStatusV1StatusWarn:
				apiHostStatus.CountWarning++
			case model.PgCheckStatusV1StatusOK:
				apiHostStatus.CountOk++
			case model.PgCheckStatusV1StatusUnkn:
				apiHostStatus.CountUnknown++
			default:
				apiHostStatus.CountUnknown++
			}
		}
	}

	return apiHostStatus, nil
}

func (c *Controller) mapPgHostV1ToAPIHostV1(pgHost *model.PgHostV1, reHostStatus *model.ReHostStatusV1) (*apipublic.HostV1, error) {
	apiHost := &apipublic.HostV1{}

	apiHost.Uid = pgHost.UID
	apiHost.Name = pgHost.Name
	apiHost.Tags = []apipublic.TagV1{}

	for _, pgTag := range pgHost.Tags {
		apiTag, err := c.mapPgTagV1ToAPITagV1(pgTag)
		if err != nil {
			return nil, err
		}

		apiHost.Tags = append(apiHost.Tags, *apiTag)
	}

	apiHostStatus, err := c.mapReHostStatusV1ToAPIHostV1Status(reHostStatus)
	if err != nil {
		return nil, err
	}

	apiHost.Status = *apiHostStatus

	return apiHost, nil
}

func (c *Controller) mapPgHostV1ToAPIGetHostsV1ResponseBody(pgHosts []*model.PgHostV1, reHostStatuses map[string]*model.ReHostStatusV1) (*apipublic.V1GetHostsJSONResponseBody, error) {
	resp := &apipublic.V1GetHostsJSONResponseBody{}

	resp.Hosts = []apipublic.HostV1{}

	for _, pgHost := range pgHosts {
		apiHost, err := c.mapPgHostV1ToAPIHostV1(pgHost, reHostStatuses[pgHost.UID])
		if err != nil {
			return nil, err
		}

		resp.Hosts = append(resp.Hosts, *apiHost)
	}

	return resp, nil
}

func (c *Controller) mapPgHostV1ToAPIGetHostV1ResponseBody(pgHost *model.PgHostV1, reHostStatus *model.ReHostStatusV1) (*apipublic.V1GetHostJSONResponseBody, error) {
	resp := &apipublic.V1GetHostJSONResponseBody{}

	apiHost, err := c.mapPgHostV1ToAPIHostV1(pgHost, reHostStatus)
	if err != nil {
		return nil, err
	}

	resp.Host = *apiHost

	return resp, nil
}

func (c *Controller) mapAPIAddHostV1RequestBodyToPgHostV1(requestBody *apipublic.V1AddHostJSONRequestBody) (*model.PgHostV1, error) {
	pgHost := &model.PgHostV1{}

	pgHost.Name = requestBody.Name
	if len(requestBody.TagUids) > 0 {
		pgHost.TagUIDs = requestBody.TagUids
	} else {
		pgHost.TagUIDs = []string{}
	}

	return pgHost, nil
}

func (c *Controller) mapAPIUpdateHostV1RequestBodyToPgHostV1(requestBody *apipublic.V1UpdateHostJSONRequestBody, oldPgHost *model.PgHostV1) (*model.PgHostV1, error) {
	tmp := *oldPgHost
	pgHost := tmp

	pgHost.Name = requestBody.Name

	if len(requestBody.TagUids) > 0 {
		pgHost.TagUIDs = requestBody.TagUids
	} else {
		pgHost.TagUIDs = []string{}
	}

	return &pgHost, nil
}

func (c *Controller) mapPgHostV1ToAPIAddHostV1ResponseBody(pgHost *model.PgHostV1) (*apipublic.V1AddHostJSONResponseBody, error) {
	resp := &apipublic.V1AddHostJSONResponseBody{}

	resp.HostUid = pgHost.UID

	return resp, nil
}
