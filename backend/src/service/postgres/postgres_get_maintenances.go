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

package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/indece-official/monitor/backend/src/model"
	"gopkg.in/guregu/null.v4"
)

type GetMaintenancesFilter struct {
	MaintenanceUID null.String
	NowActive      bool
	From           null.Int
	Size           null.Int
}

func (s *Service) GetMaintenances(qctx context.Context, filter *GetMaintenancesFilter) ([]*model.PgMaintenanceV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	limits := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_maintenance_v1.datetime_deleted IS NULL")

	if filter.MaintenanceUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_maintenance_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.MaintenanceUID.String)
	}

	if filter.NowActive {
		conditions = append(conditions, `(
			mo_maintenance_v1.datetime_start <= NOW() AND 
			(mo_maintenance_v1.datetime_finish IS NULL OR 
			mo_maintenance_v1.datetime_finish >= NOW())
		)`)
	}

	if filter.From.Valid {
		limits = append(limits, fmt.Sprintf("OFFSET $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.From.Int64)
	}

	if filter.Size.Valid {
		limits = append(limits, fmt.Sprintf("LIMIT $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Size.Int64)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_maintenance_v1.uid,
			mo_maintenance_v1.title,
			mo_maintenance_v1.message,
			mo_maintenance_v1.details,
			mo_maintenance_v1.datetime_created,
			mo_maintenance_v1.datetime_updated,
			mo_maintenance_v1.datetime_start,
			mo_maintenance_v1.datetime_finish
		FROM mo_maintenance_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_maintenance_v1.datetime_start DESC
		`+strings.Join(limits, " "),
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgMaintenances := []*model.PgMaintenanceV1{}

	for rows.Next() {
		pgMaintenance := &model.PgMaintenanceV1{}
		pgMaintenance.Details = &model.PgMaintenanceV1Details{}
		detailsJSON := []byte{}

		err = rows.Scan(
			&pgMaintenance.UID,
			&pgMaintenance.Title,
			&pgMaintenance.Message,
			&detailsJSON,
			&pgMaintenance.DatetimeCreated,
			&pgMaintenance.DatetimeUpdated,
			&pgMaintenance.DatetimeStart,
			&pgMaintenance.DatetimeFinish,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(detailsJSON, pgMaintenance.Details)
		if err != nil {
			return nil, fmt.Errorf("can't decode details for maintenance %s: %s", pgMaintenance.UID, err)
		}

		pgMaintenances = append(pgMaintenances, pgMaintenance)
	}

	return pgMaintenances, nil
}
