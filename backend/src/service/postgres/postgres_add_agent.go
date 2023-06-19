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

	"github.com/indece-official/monitor/backend/src/model"
)

func (s *Service) AddAgent(qctx context.Context, pgAgent *model.PgAgentV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	certsJSON, err := json.Marshal(pgAgent.Certs)
	if err != nil {
		return fmt.Errorf("error json encoding certs: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_agent_v1 (
			uid,
			host_uid,
			type,
			version,
			certs,
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
			NOW(),
			NOW(),
			$6
		)`,
		pgAgent.UID,
		pgAgent.HostUID,
		pgAgent.Type,
		pgAgent.Version,
		certsJSON,
		pgAgent.DatetimeRegistered,
	)
	if err != nil {
		return fmt.Errorf("can't add agent '%s': %s", pgAgent.UID, err)
	}

	return nil
}
