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

package public

import (
	"net/http"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
)

func (c *Controller) reqV1GetConfig(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	_, errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		r.Context(),
		&postgres.GetConfigPropertiesFilter{},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading config properties: %s", err)
	}

	filteredPgConfigProperties := []*model.PgConfigPropertyV1{}

	for _, pgConfigProperty := range pgConfigProperties {
		if model.PgConfigPropertyV1Protections[pgConfigProperty.Key] != model.PgConfigPropertyV1ProtectionProtected &&
			model.PgConfigPropertyV1Protections[pgConfigProperty.Key] != model.PgConfigPropertyV1ProtectionPublic {
			continue
		}

		filteredPgConfigProperties = append(filteredPgConfigProperties, pgConfigProperty)
	}

	respData, err := c.mapPgConfigPropertyV1ToAPIGetConfigV1ResponseBody(filteredPgConfigProperties)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error mapping response: %s", err)
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d config properties", len(pgConfigProperties))
}
