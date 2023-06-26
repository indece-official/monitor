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

func (s *Service) UpdateMaintenance(qctx context.Context, maintenanceUID string, pgMaintenance *model.PgMaintenanceV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	detailsJSON, err := json.Marshal(pgMaintenance.Details)
	if err != nil {
		return fmt.Errorf("error json encoding maintenance details: %s", err)
	}

	pgMaintenance.DatetimeUpdated = time.Now()

	_, err = db.ExecContext(
		qctx,
		`UPDATE mo_maintenance_v1
		SET
			title = $2,
			message = $3,
			details = $4,
			datetime_updated = $5,
			datetime_start = $6,
			datetime_finish = $7
		WHERE
			uid = $1 AND
			datetime_deleted IS NULL`,
		maintenanceUID,
		pgMaintenance.Title,
		pgMaintenance.Message,
		detailsJSON,
		pgMaintenance.DatetimeUpdated,
		pgMaintenance.DatetimeStart,
		pgMaintenance.DatetimeFinish,
	)
	if err != nil {
		return fmt.Errorf("can't update maintenance '%s': %s", maintenanceUID, err)
	}

	return nil
}
