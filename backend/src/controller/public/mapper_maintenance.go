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

func (c *Controller) mapPgMaintenanceV1DetailsAffectedToAPIMaintenanceV1Affected(pgMaintenanceAffected *model.PgMaintenanceV1DetailsAffected) (*apipublic.MaintenanceV1Affected, error) {
	apiMaintenanceAffected := &apipublic.MaintenanceV1Affected{}

	if len(pgMaintenanceAffected.HostUIDs) > 0 {
		apiMaintenanceAffected.HostUids = pgMaintenanceAffected.HostUIDs
	} else {
		apiMaintenanceAffected.HostUids = []string{}
	}

	if len(pgMaintenanceAffected.CheckUIDs) > 0 {
		apiMaintenanceAffected.CheckUids = pgMaintenanceAffected.CheckUIDs
	} else {
		apiMaintenanceAffected.CheckUids = []string{}
	}

	if len(pgMaintenanceAffected.TagUIDs) > 0 {
		apiMaintenanceAffected.TagUids = pgMaintenanceAffected.TagUIDs
	} else {
		apiMaintenanceAffected.TagUids = []string{}
	}

	return apiMaintenanceAffected, nil
}

func (c *Controller) mapAPIMaintenanceV1AffectedToPgMaintenanceV1DetailsAffected(apiMaintenanceAffected *apipublic.MaintenanceV1Affected) (*model.PgMaintenanceV1DetailsAffected, error) {
	pgMaintenanceAffected := &model.PgMaintenanceV1DetailsAffected{}

	pgMaintenanceAffected.HostUIDs = apiMaintenanceAffected.HostUids
	pgMaintenanceAffected.CheckUIDs = apiMaintenanceAffected.CheckUids
	pgMaintenanceAffected.TagUIDs = apiMaintenanceAffected.TagUids

	return pgMaintenanceAffected, nil
}

func (c *Controller) mapPgMaintenanceV1ToAPIMaintenanceV1(pgMaintenance *model.PgMaintenanceV1, addConfig bool) (*apipublic.MaintenanceV1, error) {
	var err error

	apiMaintenance := &apipublic.MaintenanceV1{}

	apiMaintenance.Uid = pgMaintenance.UID
	apiMaintenance.Title = pgMaintenance.Title
	apiMaintenance.Message = pgMaintenance.Message

	apiMaintenanceAffected, err := c.mapPgMaintenanceV1DetailsAffectedToAPIMaintenanceV1Affected(pgMaintenance.Details.Affected)
	if err != nil {
		return nil, err
	}

	apiMaintenance.Affected = *apiMaintenanceAffected

	apiMaintenance.DatetimeCreated = pgMaintenance.DatetimeCreated
	apiMaintenance.DatetimeUpdated = pgMaintenance.DatetimeUpdated
	apiMaintenance.DatetimeStart = pgMaintenance.DatetimeStart
	apiMaintenance.DatetimeFinish = pgMaintenance.DatetimeFinish.Ptr()

	return apiMaintenance, nil
}

func (c *Controller) mapPgMaintenanceV1ToAPIGetMaintenancesV1ResponseBody(pgMaintenances []*model.PgMaintenanceV1, addConfig bool) (*apipublic.V1GetMaintenancesJSONResponseBody, error) {
	resp := &apipublic.V1GetMaintenancesJSONResponseBody{}

	resp.Maintenances = []apipublic.MaintenanceV1{}

	for _, pgMaintenance := range pgMaintenances {
		apiMaintenance, err := c.mapPgMaintenanceV1ToAPIMaintenanceV1(pgMaintenance, addConfig)
		if err != nil {
			return nil, err
		}

		resp.Maintenances = append(resp.Maintenances, *apiMaintenance)
	}

	return resp, nil
}

func (c *Controller) mapPgMaintenanceV1ToAPIGetMaintenanceV1ResponseBody(pgMaintenance *model.PgMaintenanceV1, addConfig bool) (*apipublic.V1GetMaintenanceJSONResponseBody, error) {
	resp := &apipublic.V1GetMaintenanceJSONResponseBody{}

	apiMaintenance, err := c.mapPgMaintenanceV1ToAPIMaintenanceV1(pgMaintenance, addConfig)
	if err != nil {
		return nil, err
	}

	resp.Maintenance = *apiMaintenance

	return resp, nil
}

func (c *Controller) mapAPIAddMaintenanceV1RequestBodyToPgMaintenanceV1(requestBody *apipublic.V1AddMaintenanceJSONRequestBody) (*model.PgMaintenanceV1, error) {
	var err error

	pgMaintenance := &model.PgMaintenanceV1{}

	pgMaintenance.Title = requestBody.Title
	pgMaintenance.Message = requestBody.Message

	pgMaintenance.Details = &model.PgMaintenanceV1Details{}

	pgMaintenance.Details.Affected, err = c.mapAPIMaintenanceV1AffectedToPgMaintenanceV1DetailsAffected(&requestBody.Affected)
	if err != nil {
		return nil, fmt.Errorf("error mapping maintenance type: %s", err)
	}

	pgMaintenance.DatetimeStart = requestBody.DatetimeStart

	if requestBody.DatetimeFinish != nil {
		pgMaintenance.DatetimeFinish.Scan(requestBody.DatetimeFinish)
	}

	return pgMaintenance, nil
}

func (c *Controller) mapAPIUpdateMaintenanceV1RequestBodyToPgMaintenanceV1(requestBody *apipublic.V1UpdateMaintenanceJSONRequestBody, oldPgMaintenance *model.PgMaintenanceV1) (*model.PgMaintenanceV1, error) {
	var err error

	tmp := *oldPgMaintenance
	pgMaintenance := tmp

	pgMaintenance.Title = requestBody.Title
	pgMaintenance.Message = requestBody.Message

	pgMaintenance.Details.Affected, err = c.mapAPIMaintenanceV1AffectedToPgMaintenanceV1DetailsAffected(&requestBody.Affected)
	if err != nil {
		return nil, fmt.Errorf("error mapping maintenance type: %s", err)
	}

	pgMaintenance.DatetimeStart = requestBody.DatetimeStart

	if requestBody.DatetimeFinish != nil {
		pgMaintenance.DatetimeFinish.Scan(requestBody.DatetimeFinish)
	} else {
		pgMaintenance.DatetimeFinish.Scan(nil)
	}

	return &pgMaintenance, nil
}

func (c *Controller) mapPgMaintenanceV1ToAPIAddMaintenanceV1ResponseBody(pgMaintenance *model.PgMaintenanceV1) (*apipublic.V1AddMaintenanceJSONResponseBody, error) {
	resp := &apipublic.V1AddMaintenanceJSONResponseBody{}

	resp.MaintenanceUid = pgMaintenance.UID

	return resp, nil
}
