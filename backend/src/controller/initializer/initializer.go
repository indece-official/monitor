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

package initializer

import (
	"context"
	"fmt"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/cache"
	"github.com/indece-official/monitor/backend/src/service/cert"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
)

// ControllerName defines the name of the api controller used for dependency exception
const ControllerName = "initializer"

// IController is the interface of the api controller
type IController interface {
	gousu.IController
}

// Controller is the admin api controller
type Controller struct {
	log             *logger.Log
	postgresService postgres.IService
	cacheService    cache.IService
	certService     cert.IService
}

var _ IController = (*Controller)(nil)

func (c *Controller) Name() string {
	return ControllerName
}

func (c *Controller) Start() error {
	ctx := context.Background()

	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		ctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeySetupFinished,
				model.PgConfigPropertyV1KeyTLSCaCrt,
				model.PgConfigPropertyV1KeyTLSServerCrt,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error loading config properties: %s", err)
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeySetupFinished] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeySetupFinished].Value == model.PgConfigPropertyV1False {
		setupToken, err := c.cacheService.GetSetupToken()
		if err != nil {
			return fmt.Errorf("error loading setup token: %s", err)
		}

		if !setupToken.Valid || setupToken.String == "" {
			c.log.Infof("Generating setup token")

			setupTokenStr, err := utils.RandString(128)
			if err != nil {
				return fmt.Errorf("error generating setup token: %s", err)
			}

			err = c.cacheService.SetSetupToken(setupTokenStr, 60*time.Minute)
			if err != nil {
				return fmt.Errorf("error storing setup token: %s", err)
			}
		}
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value == "" ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerCrt].Value == "" {
		c.log.Infof("Generating certs")

		err = c.generateCerts(ctx)
		if err != nil {
			return fmt.Errorf("error generating certs: %s", err)
		}
	}

	return nil
}

func (c *Controller) Health() error {
	return nil
}

func (c *Controller) Stop() error {
	return nil
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	log := logger.GetLogger(fmt.Sprintf("controller.%s", ControllerName))

	return &Controller{
		log:             log,
		postgresService: ctx.GetService(postgres.ServiceName).(postgres.IService),
		cacheService:    ctx.GetService(cache.ServiceName).(cache.IService),
		certService:     ctx.GetService(cert.ServiceName).(cert.IService),
	}
}

var _ gousu.ControllerFactory = NewController
