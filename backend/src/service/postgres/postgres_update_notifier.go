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

func (s *Service) UpdateNotifier(qctx context.Context, notifierUID string, pgNotifier *model.PgNotifierV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	configJSON, err := json.Marshal(pgNotifier.Config)
	if err != nil {
		return fmt.Errorf("error json encoding notifier config: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`UPDATE mo_notifier_v1
		SET
			name = $2,
			config = $3,
			datetime_updated = NOW()
		WHERE
			uid = $1 AND
			datetime_deleted IS NULL`,
		notifierUID,
		pgNotifier.Name,
		configJSON,
	)
	if err != nil {
		return fmt.Errorf("can't update notifier '%s': %s", notifierUID, err)
	}

	return nil
}
