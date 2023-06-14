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
	"github.com/lib/pq"
)

func (s *Service) AddUser(qctx context.Context, pgUser *model.PgUserV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_user_v1 (
			uid,
			source,
			username,
			name,
			email,
			local_roles,
			password_hash,
			datetime_created,
			datetime_updated,
			datetime_locked
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			NOW(),
			NOW(),
			$8
		)`,
		pgUser.UID,
		pgUser.Source,
		pgUser.Username,
		pgUser.Name,
		pgUser.Email,
		pq.Array(pgUser.LocalRoles),
		pgUser.PasswordHash,
		pgUser.DatetimeLocked,
	)
	if err != nil {
		return fmt.Errorf("can't add user '%s': %s", pgUser.UID, err)
	}

	return nil
}
