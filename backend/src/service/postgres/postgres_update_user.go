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

func (s *Service) UpdateUser(qctx context.Context, userUID string, pgUser *model.PgUserV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`UPDATE mo_user_v1
		SET
			name = $2,
			email = $3,
			local_roles = $4,
			datetime_updated = NOW(),
			datetime_locked = $5
		WHERE
			uid = $1 AND
			datetime_deleted IS NULL`,
		userUID,
		pgUser.Name,
		pgUser.Email,
		pq.Array(pgUser.LocalRoles),
		pgUser.DatetimeLocked,
	)
	if err != nil {
		return fmt.Errorf("can't update user '%s': %s", userUID, err)
	}

	return nil
}
