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

func (s *Service) AddChecker(qctx context.Context, pgChecker *model.PgCheckerV1) error {
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
		`INSERT INTO mo_checker_v1 (
			uid,
			type,
			agent_type,
			version,
			name,
			capabilities,
			custom_checks,
			datetime_created,
			datetime_updated
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
			NOW()
		)`,
		pgChecker.UID,
		pgChecker.Type,
		pgChecker.AgentType,
		pgChecker.Version,
		pgChecker.Name,
		capabilitiesJSON,
		pgChecker.CustomChecks,
	)
	if err != nil {
		return fmt.Errorf("can't add checker '%s': %s", pgChecker.UID, err)
	}

	return nil
}
