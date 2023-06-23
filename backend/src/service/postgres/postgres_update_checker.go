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

func (s *Service) UpdateChecker(qctx context.Context, checkerUID string, pgChecker *model.PgCheckerV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	capabilitiesJSON, err := json.Marshal(pgChecker.Capabilities)
	if err != nil {
		return fmt.Errorf("error json encoding checker capabilities: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`UPDATE mo_checker_v1
		SET
			type = $2,
			version = $3,
			name = $4,
			capabilities = $5,
			custom_checks = $6,
			datetime_updated = NOW()
		WHERE
			uid = $1 AND
			datetime_deleted IS NULL`,
		checkerUID,
		pgChecker.Type,
		pgChecker.Version,
		pgChecker.Name,
		capabilitiesJSON,
		pgChecker.CustomChecks,
	)
	if err != nil {
		return fmt.Errorf("can't update checker '%s': %s", checkerUID, err)
	}

	return nil
}
