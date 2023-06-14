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
)

// UpsertConfigProperty upserts config properties
func (s *Service) UpsertConfigProperty(qctx context.Context, pgConfigProperty *model.PgConfigPropertyV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_configproperty_v1 (
			key,
			value,
			datetime_created,
			datetime_updated
		)
		VALUES(
			$1,
			$2,
			NOW(),
			NOW()
		)
		ON CONFLICT (key) DO UPDATE
		SET
			value = $2,
			datetime_updated = NOW()`,
		pgConfigProperty.Key,
		pgConfigProperty.Value,
	)
	if err != nil {
		return fmt.Errorf("can't upsert confi property %s: %s", pgConfigProperty.Key, err)
	}

	return nil
}
