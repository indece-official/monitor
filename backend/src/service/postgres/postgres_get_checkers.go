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

type GetCheckersFilter struct {
	CheckerUID null.String
	Type       null.String
	AgentType  null.String
}

func (s *Service) GetCheckers(qctx context.Context, filter *GetCheckersFilter) ([]*model.PgCheckerV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_checker_v1.datetime_deleted IS NULL")

	if filter.CheckerUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_checker_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.CheckerUID.String)
	}

	if filter.Type.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_checker_v1.type = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Type.String)
	}

	if filter.AgentType.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_checker_v1.agent_type = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.AgentType.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			mo_checker_v1.uid,
			mo_checker_v1.type,
			mo_checker_v1.agent_type,
			mo_checker_v1.version,
			mo_checker_v1.name,
			mo_checker_v1.capabilities,
			mo_checker_v1.custom_checks
		FROM mo_checker_v1
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY mo_checker_v1.name ASC`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgCheckers := []*model.PgCheckerV1{}

	for rows.Next() {
		pgChecker := &model.PgCheckerV1{}
		pgChecker.Capabilities = &model.PgCheckerV1Capabilities{}
		capabilitiesJSON := []byte{}

		err = rows.Scan(
			&pgChecker.UID,
			&pgChecker.Type,
			&pgChecker.AgentType,
			&pgChecker.Name,
			&pgChecker.Version,
			&capabilitiesJSON,
			&pgChecker.CustomChecks,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(capabilitiesJSON, pgChecker.Capabilities)
		if err != nil {
			return nil, fmt.Errorf("can't decode capabilities for checker %s: %s", pgChecker.UID, err)
		}

		pgCheckers = append(pgCheckers, pgChecker)
	}

	return pgCheckers, nil
}
