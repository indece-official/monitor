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
)

type GetCheckStatusesFilter struct {
	CheckUID    string
	CountStatus int
}

func (s *Service) GetCheckStatuses(qctx context.Context, filter *GetCheckStatusesFilter) ([]*model.PgCheckStatusV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	limitStatus := fmt.Sprintf("LIMIT $%d", len(conditionParams)+1)
	conditionParams = append(conditionParams, filter.CountStatus)

	conditions = append(conditions, "mo_checkstatus_v1.datetime_deleted IS NULL")

	conditions = append(conditions, fmt.Sprintf("mo_checkstatus_v1.check_uid = $%d", len(conditionParams)+1))
	conditionParams = append(conditionParams, filter.CheckUID)

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_checkstatus_v1.uid,
			mo_checkstatus_v1.check_uid,
			mo_checkstatus_v1.status,
			mo_checkstatus_v1.message,
			mo_checkstatus_v1.data,
			mo_checkstatus_v1.datetime_created
		FROM mo_checkstatus_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_checkstatus_v1.datetime_created DESC
		`+limitStatus,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgCheckStatuses := []*model.PgCheckStatusV1{}

	for rows.Next() {
		pgCheckStatus := &model.PgCheckStatusV1{}
		pgCheckStatus.Data = map[string]interface{}{}
		dataJSON := []byte{}

		err = rows.Scan(
			&pgCheckStatus.UID,
			&pgCheckStatus.CheckUID,
			&pgCheckStatus.Status,
			&pgCheckStatus.Message,
			&dataJSON,
			&pgCheckStatus.DatetimeCreated,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(dataJSON, &pgCheckStatus.Data)
		if err != nil {
			return nil, fmt.Errorf("can't decode data for check status %s: %s", pgCheckStatus.UID, err)
		}

		pgCheckStatuses = append(pgCheckStatuses, pgCheckStatus)
	}

	return pgCheckStatuses, nil
}
