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
	"fmt"
	"strings"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/lib/pq"
)

type GetConfigPropertiesFilter struct {
	Keys []model.PgConfigPropertyV1Key
}

func (s *Service) GetConfigProperties(qctx context.Context, filter *GetConfigPropertiesFilter) (map[model.PgConfigPropertyV1Key]*model.PgConfigPropertyV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "1 = 1")

	if len(filter.Keys) > 0 {
		conditions = append(conditions, fmt.Sprintf("mo_configproperty_v1.key = ANY($%d)", len(conditionParams)+1))

		keys := []string{}
		for _, key := range filter.Keys {
			keys = append(keys, string(key))
		}

		conditionParams = append(conditionParams, pq.Array(keys))
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_configproperty_v1.key,
			mo_configproperty_v1.value
		FROM mo_configproperty_v1
		WHERE `+strings.Join(conditions, " AND "),
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgConfigProperties := map[model.PgConfigPropertyV1Key]*model.PgConfigPropertyV1{}

	for rows.Next() {
		pgConfigProperty := &model.PgConfigPropertyV1{}

		err = rows.Scan(
			&pgConfigProperty.Key,
			&pgConfigProperty.Value,
		)
		if err != nil {
			return nil, err
		}

		pgConfigProperties[pgConfigProperty.Key] = pgConfigProperty
	}

	return pgConfigProperties, nil
}
