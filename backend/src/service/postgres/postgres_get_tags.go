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
	"gopkg.in/guregu/null.v4"
)

type GetTagsFilter struct {
	TagUID null.String
}

func (s *Service) GetTags(qctx context.Context, filter *GetTagsFilter) ([]*model.PgTagV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	limits := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_tag_v1.datetime_deleted IS NULL")

	if filter.TagUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_tag_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.TagUID.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_tag_v1.uid,
			mo_tag_v1.name,
			mo_tag_v1.color
		FROM mo_tag_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_tag_v1.name ASC
		`+strings.Join(limits, " "),
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgTags := []*model.PgTagV1{}

	for rows.Next() {
		pgTag := &model.PgTagV1{}

		err = rows.Scan(
			&pgTag.UID,
			&pgTag.Name,
			&pgTag.Color,
		)
		if err != nil {
			return nil, err
		}

		pgTags = append(pgTags, pgTag)
	}

	return pgTags, nil
}
