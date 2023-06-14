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

func (s *Service) AddCheck(qctx context.Context, pgCheck *model.PgCheckV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	configJSON, err := json.Marshal(pgCheck.Config)
	if err != nil {
		return fmt.Errorf("error json encoding check config: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_check_v1 (
			uid,
			host_uid,
			checker_uid,
			name,
			type,
			schedule,
			config,
			custom,
			datetime_created,
			datetime_updated,
			datetime_disabled
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			NOW(),
			NOW(),
			$9
		)`,
		pgCheck.UID,
		pgCheck.HostUID,
		pgCheck.CheckerUID,
		pgCheck.Name,
		pgCheck.Type,
		pgCheck.Schedule,
		configJSON,
		pgCheck.Custom,
		pgCheck.DatetimeDisabled,
	)
	if err != nil {
		return fmt.Errorf("can't add check '%s': %s", pgCheck.UID, err)
	}

	return nil
}
