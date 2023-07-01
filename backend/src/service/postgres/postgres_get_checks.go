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

type GetChecksFilter struct {
	CheckUID   null.String
	CheckerUID null.String
	Type       null.String
	HostUID    null.String
	Disabled   null.Bool
}

func (s *Service) GetChecks(qctx context.Context, filter *GetChecksFilter) ([]*model.PgCheckV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_check_v1.datetime_deleted IS NULL")

	if filter.CheckUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_check_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.CheckUID.String)
	}

	if filter.CheckerUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_check_v1.checker_uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.CheckerUID.String)
	}

	if filter.Type.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_check_v1.type = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Type.String)
	}

	if filter.HostUID.Valid {
		conditions = append(conditions, fmt.Sprintf(`
			EXISTS (
				SELECT 1
				FROM mo_checker_v1
				INNER JOIN mo_agent_v1 ON
					mo_agent_v1.uid = mo_checker_v1.agent_uid AND
					mo_agent_v1.datetime_deleted IS NULL
				WHERE
					mo_checker_v1.uid = mo_check_v1.checker_uid AND
					mo_checker_v1.datetime_deleted IS NULL AND
					mo_agent_v1.host_uid = $%d
		)`, len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.HostUID.String)
	}

	if filter.Disabled.Valid && filter.Disabled.Bool {
		conditions = append(conditions, "mo_check_v1.datetime_disabled IS NOT NULL")
	} else if filter.Disabled.Valid && !filter.Disabled.Bool {
		conditions = append(conditions, "mo_check_v1.datetime_disabled IS NULL")
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_check_v1.uid,
			mo_check_v1.checker_uid,
			mo_check_v1.name,
			mo_check_v1.type,
			mo_check_v1.schedule,
			mo_check_v1.config,
			mo_check_v1.custom,
			mo_check_v1.datetime_disabled
		FROM mo_check_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_check_v1.name ASC`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgChecks := []*model.PgCheckV1{}

	for rows.Next() {
		pgCheck := &model.PgCheckV1{}
		pgCheck.Config = &model.PgCheckV1Config{}
		configJSON := []byte{}

		err = rows.Scan(
			&pgCheck.UID,
			&pgCheck.CheckerUID,
			&pgCheck.Name,
			&pgCheck.Type,
			&pgCheck.Schedule,
			&configJSON,
			&pgCheck.Custom,
			&pgCheck.DatetimeDisabled,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configJSON, pgCheck.Config)
		if err != nil {
			return nil, fmt.Errorf("can't decode config for check %s: %s", pgCheck.UID, err)
		}

		pgChecks = append(pgChecks, pgCheck)
	}

	return pgChecks, nil
}
