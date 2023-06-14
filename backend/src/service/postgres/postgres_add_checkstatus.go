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

func (s *Service) AddCheckStatus(qctx context.Context, pgCheckStatus *model.PgCheckStatusV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	dataJSON, err := json.Marshal(pgCheckStatus.Data)
	if err != nil {
		return fmt.Errorf("error json encoding check config: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_checkstatus_v1 (
			uid,
			check_uid,
			status,
			message,
			data,
			datetime_created,
			datetime_updated
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			NOW(),
			NOW()
		)`,
		pgCheckStatus.UID,
		pgCheckStatus.CheckUID,
		pgCheckStatus.Status,
		pgCheckStatus.Message,
		dataJSON,
	)
	if err != nil {
		return fmt.Errorf("can't add check status '%s': %s", pgCheckStatus.UID, err)
	}

	return nil
}
