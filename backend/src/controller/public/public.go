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

// Package public API controller
//
// This is the public api of the service - all operations are projected with JWTs
// if required.
package public

import (
	"fmt"

	"github.com/indece-official/monitor/backend/src/service/cache"
	"github.com/indece-official/monitor/backend/src/service/cert"
	"github.com/indece-official/monitor/backend/src/service/postgres"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/namsral/flag"
)

const ControllerName = "public"

var (
	serverHost   = flag.String("server_host", "0.0.0.0", "")
	serverPort   = flag.Int("server_port", 8080, "")
	cookieSecure = flag.Bool("cookie_secure", true, "")
)

// IController is the interface of the public api controller
type IController interface {
	gousu.IController
}

// Controller is the public api controller
type Controller struct {
	baseController  *gousuchi.AbstractController
	log             *logger.Log
	postgresService postgres.IService
	cacheService    cache.IService
	certService     cert.IService
}

var _ IController = (*Controller)(nil)

func (c *Controller) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
	)

	router.Get("/*", c.reqStaticFile)
	router.Get("/env.js", c.baseController.Wrap(c.reqEnvJS))

	router.Get("/v1/health/sanity", c.baseController.Wrap(c.reqGetHealthSanity))

	router.Post("/api/v1/login", c.baseController.Wrap(c.reqV1Login))
	router.Post("/api/v1/logout", c.baseController.Wrap(c.reqV1Logout))

	router.Post("/api/v1/setup/finish", c.baseController.Wrap(c.reqV1FinishSetup))

	router.Get("/api/v1/connector", c.baseController.Wrap(c.reqV1GetConnectors))
	router.Post("/api/v1/connector", c.baseController.Wrap(c.reqV1AddConnector))
	router.Get("/api/v1/connector/{connectorUID}", c.baseController.Wrap(c.reqV1GetConnector))
	router.Delete("/api/v1/connector/{connectorUID}", c.baseController.Wrap(c.reqV1DeleteConnector))

	router.Get("/api/v1/user", c.baseController.Wrap(c.reqV1GetUsers))
	router.Get("/api/v1/user/self", c.baseController.Wrap(c.reqV1GetOwnUser))
	router.Get("/api/v1/user/{userUID}", c.baseController.Wrap(c.reqV1GetUser))
	router.Post("/api/v1/user", c.baseController.Wrap(c.reqV1AddUser))
	router.Put("/api/v1/user/{userUID}", c.baseController.Wrap(c.reqV1UpdateUser))
	router.Delete("/api/v1/user/{userUID}", c.baseController.Wrap(c.reqV1DeleteUser))

	router.Get("/api/v1/checker", c.baseController.Wrap(c.reqV1GetCheckers))
	router.Get("/api/v1/checker/{checkerUID}", c.baseController.Wrap(c.reqV1GetChecker))

	router.Get("/api/v1/check", c.baseController.Wrap(c.reqV1GetChecks))
	router.Post("/api/v1/check", c.baseController.Wrap(c.reqV1AddCheck))
	router.Put("/api/v1/check/{checkUID}", c.baseController.Wrap(c.reqV1UpdateCheck))
	router.Delete("/api/v1/check/{checkUID}", c.baseController.Wrap(c.reqV1DeleteCheck))
	router.Get("/api/v1/check/{checkUID}", c.baseController.Wrap(c.reqV1GetCheck))

	router.Get("/api/v1/tag", c.baseController.Wrap(c.reqV1GetTags))
	router.Post("/api/v1/tag", c.baseController.Wrap(c.reqV1AddTag))
	router.Put("/api/v1/tag/{tagUID}", c.baseController.Wrap(c.reqV1UpdateTag))
	router.Delete("/api/v1/tag/{tagUID}", c.baseController.Wrap(c.reqV1DeleteTag))
	router.Get("/api/v1/tag/{tagUID}", c.baseController.Wrap(c.reqV1GetTag))

	router.Get("/api/v1/config", c.baseController.Wrap(c.reqV1GetConfig))
	router.Put("/api/v1/config/{key}", c.baseController.Wrap(c.reqV1SetConfigProperty))

	router.Get("/api/v1/host", c.baseController.Wrap(c.reqV1GetHosts))
	router.Post("/api/v1/host", c.baseController.Wrap(c.reqV1AddHost))
	router.Put("/api/v1/host/{hostUID}", c.baseController.Wrap(c.reqV1UpdateHost))
	router.Delete("/api/v1/host/{hostUID}", c.baseController.Wrap(c.reqV1DeleteHost))
	router.Get("/api/v1/host/{hostUID}", c.baseController.Wrap(c.reqV1GetHost))

	router.Get("/api/v1/host/{hostUID}/check", c.baseController.Wrap(c.reqV1GetHostChecks))

	router.Get("/api/v1/notifier", c.baseController.Wrap(c.reqV1GetNotifiers))
	router.Post("/api/v1/notifier", c.baseController.Wrap(c.reqV1AddNotifier))
	router.Put("/api/v1/notifier/{notifierUID}", c.baseController.Wrap(c.reqV1UpdateNotifier))
	router.Delete("/api/v1/notifier/{notifierUID}", c.baseController.Wrap(c.reqV1DeleteNotifier))
	router.Get("/api/v1/notifier/{notifierUID}", c.baseController.Wrap(c.reqV1GetNotifier))

	return router
}

func (c *Controller) Name() string {
	return ControllerName
}

// Start starts the api server in a new go-func
func (c *Controller) Start() error {
	c.baseController.UseHost(*serverHost)
	c.baseController.UsePort(*serverPort)
	c.baseController.UseRouter(c.getRouter())

	return c.baseController.Start()
}

// Health checks if the api server has thrown unresolvable internal errors
func (c *Controller) Health() error {
	return c.baseController.Health()
}

// Stop currently does nothing
func (c *Controller) Stop() error {
	return c.baseController.Stop()
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	log := logger.GetLogger(fmt.Sprintf("controller.%s", ControllerName))

	return &Controller{
		log:             log,
		baseController:  gousuchi.NewAbstractController(log),
		postgresService: ctx.GetService(postgres.ServiceName).(postgres.IService),
		cacheService:    ctx.GetService(cache.ServiceName).(cache.IService),
		certService:     ctx.GetService(cert.ServiceName).(cert.IService),
	}
}

var _ gousu.ControllerFactory = NewController
