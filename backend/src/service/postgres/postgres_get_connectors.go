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

type GetConnectorsFilter struct {
	ConnectorUID null.String
	Token        null.String
}

func (s *Service) GetConnectors(qctx context.Context, filter *GetConnectorsFilter) ([]*model.PgConnectorV1, error) {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return nil, fmt.Errorf("error acquiring db connection: %s", err)
	}

	conditions := []string{}
	limits := []string{}
	conditionParams := []interface{}{}

	conditions = append(conditions, "mo_connector_v1.datetime_deleted IS NULL")

	if filter.ConnectorUID.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_connector_v1.uid = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.ConnectorUID.String)
	}

	if filter.Token.Valid {
		conditions = append(conditions, fmt.Sprintf("mo_connector_v1.token = $%d", len(conditionParams)+1))
		conditionParams = append(conditionParams, filter.Token.String)
	}

	// #nosec G202 -- Query parameters are used for all input data
	rows, err := db.QueryContext(
		qctx,
		`SELECT
			connector.uid,
			connector.host_uid,
			connector.type,
			connector.version,
			connector.tls_client_crt,
			connector.tls_client_key,
			connector.datetime_registered
		FROM (
			SELECT
				*
			FROM mo_connector_v1
			WHERE `+strings.Join(conditions, " AND ")+`
			ORDER BY mo_connector_v1.datetime_created DESC
			`+strings.Join(limits, " ")+`
		) as connector`,
		conditionParams...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pgConnectors := []*model.PgConnectorV1{}

	for rows.Next() {
		pgConnector := &model.PgConnectorV1{}

		err = rows.Scan(
			&pgConnector.UID,
			&pgConnector.HostUID,
			&pgConnector.Type,
			&pgConnector.Version,
			&pgConnector.TLSClientCrt,
			&pgConnector.TLSClientKey,
			&pgConnector.DatetimeRegistered,
		)
		if err != nil {
			return nil, err
		}

		pgConnectors = append(pgConnectors, pgConnector)
	}

	return pgConnectors, nil
}
