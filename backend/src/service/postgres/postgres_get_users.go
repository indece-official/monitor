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
	"gopkg.in/guregu/null.v4"
)

type GetUsersFilter struct {
	UserUID  null.String
	Username null.String
}

func (s *Service) GetUsers(qctx context.Context, filter *GetUsersFilter) ([]*model.PgUserV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	limits := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_user_v1.datetime_deleted IS NULL")

	if filter.UserUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_user_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.UserUID.String)
	}

	if filter.Username.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_user_v1.username = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Username.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_user_v1.uid,
			mo_user_v1.source,
			mo_user_v1.username,
			mo_user_v1.name,
			mo_user_v1.email,
			mo_user_v1.local_roles,
			mo_user_v1.password_hash,
			mo_user_v1.datetime_locked
		FROM mo_user_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_user_v1.username ASC
		`+strings.Join(limits, " "),
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgUsers := []*model.PgUserV1{}

	for rows.Next() {
		pgUser := &model.PgUserV1{}
		localRoles := []string{}

		err = rows.Scan(
			&pgUser.UID,
			&pgUser.Source,
			&pgUser.Username,
			&pgUser.Name,
			&pgUser.Email,
			pq.Array(&localRoles),
			&pgUser.PasswordHash,
			&pgUser.DatetimeLocked,
		)
		if err != nil {
			return nil, err
		}

		pgUser.LocalRoles = []model.PgUserV1Role{}
		for _, localRole := range localRoles {
			pgUser.LocalRoles = append(pgUser.LocalRoles, model.PgUserV1Role(localRole))
		}

		pgUsers = append(pgUsers, pgUser)
	}

	return pgUsers, nil
}
