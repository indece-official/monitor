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

type GetAgentsFilter struct {
	AgentUID null.String
	Token    null.String
}

func (s *Service) GetAgents(qctx context.Context, filter *GetAgentsFilter) ([]*model.PgAgentV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	limits := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_agent_v1.datetime_deleted IS NULL")

	if filter.AgentUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_agent_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.AgentUID.String)
	}

	if filter.Token.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_agent_v1.token = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Token.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			agent.uid,
			agent.host_uid,
			agent.type,
			agent.version,
			agent.certs,
			agent.datetime_registered
		FROM (
			SELECT
				*
			FROM mo_agent_v1
			WHERE `+strings.Join(conditions, " AND ")+`
			ORDER BY mo_agent_v1.datetime_created DESC
			`+strings.Join(limits, " ")+`
		) as agent`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgAgents := []*model.PgAgentV1{}

	for rows.Next() {
		pgAgent := &model.PgAgentV1{}
		pgAgent.Certs = &model.PgAgentV1Certs{}
		certsJSON := []byte{}

		err = rows.Scan(
			&pgAgent.UID,
			&pgAgent.HostUID,
			&pgAgent.Type,
			&pgAgent.Version,
			&certsJSON,
			&pgAgent.DatetimeRegistered,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(certsJSON, pgAgent.Certs)
		if err != nil {
			return nil, fmt.Errorf("can't decode certs for agent %s: %s", pgAgent.UID, err)
		}

		pgAgents = append(pgAgents, pgAgent)
	}

	return pgAgents, nil
}
