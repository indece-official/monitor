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

type GetHostsFilter struct {
	HostUID null.String
}

func (s *Service) GetHosts(qctx context.Context, filter *GetHostsFilter) ([]*model.PgHostV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_host_v1.datetime_deleted IS NULL")

	if filter.HostUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_host_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.HostUID.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_host_v1.uid,
			mo_host_v1.name,
			mo_host_v1.tag_uids,
			mo_tag_v1.uid,
			mo_tag_v1.name,
			mo_tag_v1.color
		FROM mo_host_v1
		LEFT JOIN mo_tag_v1 ON
			mo_tag_v1.uid = ANY(mo_host_v1.tag_uids) AND
			mo_tag_v1.datetime_deleted IS NULL
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_host_v1.name ASC`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgHosts := []*model.PgHostV1{}
	var pgLastHost *model.PgHostV1

	for rows.Next() {
		pgHost := &model.PgHostV1{}
		pgTagUID := null.String{}
		pgTagName := null.String{}
		pgTagColor := null.String{}

		err = rows.Scan(
			&pgHost.UID,
			&pgHost.Name,
			pq.Array(&pgHost.TagUIDs),
			&pgTagUID,
			&pgTagName,
			&pgTagColor,
		)
		if err != nil {
			return nil, err
		}

		if pgLastHost == nil || pgLastHost.UID != pgHost.UID {
			pgHosts = append(pgHosts, pgHost)
			pgLastHost = pgHost
		}

		if pgTagUID.Valid {
			pgTag := &model.PgTagV1{}

			pgTag.UID = pgTagUID.String
			pgTag.Name = pgTagName.String
			pgTag.Color = pgTagColor.String

			pgLastHost.Tags = append(pgLastHost.Tags, pgTag)
		}
	}

	return pgHosts, nil
}
