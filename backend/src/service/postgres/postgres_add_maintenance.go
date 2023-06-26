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
	"time"

	"github.com/indece-official/monitor/backend/src/model"
)

func (s *Service) AddMaintenance(qctx context.Context, pgMaintenance *model.PgMaintenanceV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	detailsJSON, err := json.Marshal(pgMaintenance.Details)
	if err != nil {
		return fmt.Errorf("error json encoding maintenance details: %s", err)
	}

	pgMaintenance.DatetimeCreated = time.Now()
	pgMaintenance.DatetimeUpdated = time.Now()

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_maintenance_v1 (
			uid,
			title,
			message,
			details,
			datetime_created,
			datetime_updated,
			datetime_start,
			datetime_finish
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)`,
		pgMaintenance.UID,
		pgMaintenance.Title,
		pgMaintenance.Message,
		detailsJSON,
		pgMaintenance.DatetimeCreated,
		pgMaintenance.DatetimeUpdated,
		pgMaintenance.DatetimeStart,
		pgMaintenance.DatetimeFinish,
	)
	if err != nil {
		return fmt.Errorf("can't add maintenance '%s': %s", pgMaintenance.UID, err)
	}

	return nil
}
