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

package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/cache"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/service/smtp"
	"github.com/indece-official/monitor/backend/src/service/template"
)

const ControllerName = "cron"

type IController interface {
	gousu.IController
}

type Controller struct {
	log             *logger.Log
	postgresService postgres.IService
	cacheService    cache.IService
	templateService template.IService
	smtpService     smtp.IService
	checks          []*model.PgCheckV1
	checkers        []*model.PgCheckerV1
	agents          []*model.PgAgentV1
	hosts           []*model.PgHostV1
	users           []*model.PgUserV1
	notifiers       []*model.PgNotifierV1
	scheduler       *gocron.Scheduler
	mutexChecks     sync.Mutex
	mutexCheckers   sync.Mutex
	mutexAgents     sync.Mutex
	mutexHosts      sync.Mutex
	mutexUsers      sync.Mutex
	mutexNotifiers  sync.Mutex
	stop            bool
	error           error
	waitGroupStop   sync.WaitGroup
}

var _ IController = (*Controller)(nil)

func (c *Controller) Name() string {
	return ControllerName
}

func (c *Controller) Start() error {
	c.scheduler = gocron.NewScheduler(time.UTC)

	err := c.reload()
	if err != nil {
		return fmt.Errorf("error loading cron controller: %s", err)
	}

	c.waitGroupStop.Add(1)
	go func() {
		for !c.stop {
			c.error = nil

			err := c.reloadLoop()
			if err != nil {
				c.error = err

				c.log.Errorf("Error in reload loop: %s", err)
			}

			time.Sleep(5 * time.Second)
		}

		c.waitGroupStop.Done()
	}()

	c.waitGroupStop.Add(1)
	go func() {
		time.Sleep(10 * time.Second)

		for !c.stop {
			c.error = nil

			err := c.agentEventLoop()
			if err != nil {
				c.error = err

				c.log.Errorf("Error in agent event loop: %s", err)
			}

			time.Sleep(5 * time.Second)
		}

		c.waitGroupStop.Done()
	}()

	c.waitGroupStop.Add(1)
	go func() {
		for !c.stop {
			c.error = nil

			err := c.actionTimeoutLoop()
			if err != nil {
				c.error = err

				c.log.Errorf("Error in check action timeout loop: %s", err)
			}

			time.Sleep(5 * time.Second)
		}

		c.waitGroupStop.Done()
	}()

	c.waitGroupStop.Add(1)
	go func() {
		for !c.stop {
			c.error = nil

			err := c.notifyLoop()
			if err != nil {
				c.error = err

				c.log.Errorf("Error in notify loop: %s", err)
			}

			time.Sleep(5 * time.Second)
		}

		c.waitGroupStop.Done()
	}()

	c.scheduler.StartAsync()

	return nil
}

func (c *Controller) Health() error {
	return c.error
}

func (c *Controller) Stop() error {
	c.stop = true

	c.scheduler.Stop()
	c.scheduler.Clear()

	c.waitGroupStop.Wait()

	return nil
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	log := logger.GetLogger(fmt.Sprintf("controller.%s", ControllerName))

	return &Controller{
		log:             log,
		postgresService: ctx.GetService(postgres.ServiceName).(postgres.IService),
		cacheService:    ctx.GetService(cache.ServiceName).(cache.IService),
		templateService: ctx.GetService(template.ServiceName).(template.IService),
		smtpService:     ctx.GetService(smtp.ServiceName).(smtp.IService),
	}
}

var _ gousu.ControllerFactory = NewController
