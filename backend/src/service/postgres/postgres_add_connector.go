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

	"github.com/indece-official/monitor/backend/src/model"
)

func (s *Service) AddConnector(qctx context.Context, pgConnector *model.PgConnectorV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_connector_v1 (
			uid,
			host_uid,
			type,
			version,
			tls_client_crt,
			tls_client_key,
			datetime_created,
			datetime_updated,
			datetime_registered
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			NOW(),
			NOW(),
			$7
		)`,
		pgConnector.UID,
		pgConnector.HostUID,
		pgConnector.Type,
		pgConnector.Version,
		pgConnector.TLSClientCrt,
		pgConnector.TLSClientKey,
		pgConnector.DatetimeRegistered,
	)
	if err != nil {
		return fmt.Errorf("can't add connector '%s': %s", pgConnector.UID, err)
	}

	return nil
}
