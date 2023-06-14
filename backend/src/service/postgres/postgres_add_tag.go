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

func (s *Service) AddTag(qctx context.Context, pgTag *model.PgTagV1) error {
	db, err := s.postgresService.GetDBSafe()
	if err != nil {
		return fmt.Errorf("error acquiring db connection: %s", err)
	}

	_, err = db.ExecContext(
		qctx,
		`INSERT INTO mo_tag_v1 (
			uid,
			name,
			color,
			datetime_created,
			datetime_updated
		)
		VALUES(
			$1,
			$2,
			$3,
			NOW(),
			NOW()
		)`,
		pgTag.UID,
		pgTag.Name,
		pgTag.Color,
	)
	if err != nil {
		return fmt.Errorf("can't add tag '%s': %s", pgTag.UID, err)
	}

	return nil
}
