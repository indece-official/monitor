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

type GetNotifiersFilter struct {
	NotifierUID null.String
	Disabled    null.Bool
}

func (s *Service) GetNotifiers(qctx context.Context, filter *GetNotifiersFilter) ([]*model.PgNotifierV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_notifier_v1.datetime_deleted IS NULL")

	if filter.NotifierUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_notifier_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.NotifierUID.String)
	}

	if filter.Disabled.Valid && filter.Disabled.Bool {
		conditions = append(conditions, "mo_notifier_v1.datetime_disabled IS NOT NULL")
	} else if filter.Disabled.Valid && !filter.Disabled.Bool {
		conditions = append(conditions, "mo_notifier_v1.datetime_disabled IS NULL")
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_notifier_v1.uid,
			mo_notifier_v1.name,
			mo_notifier_v1.type,
			mo_notifier_v1.config,
			mo_notifier_v1.datetime_disabled
		FROM mo_notifier_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_notifier_v1.name ASC`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgNotifiers := []*model.PgNotifierV1{}

	for rows.Next() {
		pgNotifier := &model.PgNotifierV1{}
		pgNotifier.Config = &model.PgNotifierV1Config{}
		configJSON := []byte{}

		err = rows.Scan(
			&pgNotifier.UID,
			&pgNotifier.Name,
			&pgNotifier.Type,
			&configJSON,
			&pgNotifier.DatetimeDisabled,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configJSON, pgNotifier.Config)
		if err != nil {
			return nil, fmt.Errorf("can't decode config for notifier %s: %s", pgNotifier.UID, err)
		}

		pgNotifiers = append(pgNotifiers, pgNotifier)
	}

	return pgNotifiers, nil
}
