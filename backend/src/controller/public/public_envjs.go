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
	"github.com/indece-official/monitor/backend/src/buildvars"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reqEnvJS(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		r.Context(),
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeySetupFinished,
			},
		},
	)
	if err != nil {
		return gousuchi.InternalServerError(r, "Error loading config properties: %s", err)
	}

	setupEnabled := "false"
	setupToken := null.String{}

	if pgConfigProperties[model.PgConfigPropertyV1KeySetupFinished] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeySetupFinished].Value == model.PgConfigPropertyV1False {
		setupEnabled = "true"

		setupToken, err = c.cacheService.GetSetupToken()
		if err != nil {
			return gousuchi.InternalServerError(r, "Error loading setup token: %s", err)
		}
	}

	body := `window.CONFIG = {
			SERVER_VERSION: '` + buildvars.BuildVersion + `',
			SETUP_ENABLED: ` + setupEnabled + `,
			SETUP_TOKEN: '` + setupToken.String + `'
		};
		`

	return gousuchi.NewResponse(
		r,
		http.StatusOK,
		gousuchi.ContentType("text/javascript"),
		[]byte(body),
	)
}
