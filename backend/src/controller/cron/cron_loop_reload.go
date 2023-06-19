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
	"context"
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) reload() error {
	ctx := context.Background()

	pgChecks, err := c.postgresService.GetChecks(
		ctx,
		&postgres.GetChecksFilter{
			Disabled:    null.BoolFrom(false),
			CountStatus: 1,
		},
	)
	if err != nil {
		return fmt.Errorf("error loading checks: %s", err)
	}

	c.mutexChecks.Lock()
	c.checks = pgChecks
	c.mutexChecks.Unlock()

	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{},
	)
	if err != nil {
		return fmt.Errorf("error loading checkers: %s", err)
	}

	c.mutexCheckers.Lock()
	c.checkers = pgCheckers
	c.mutexCheckers.Unlock()

	pgAgents, err := c.postgresService.GetAgents(
		ctx,
		&postgres.GetAgentsFilter{},
	)
	if err != nil {
		return fmt.Errorf("error loading agents: %s", err)
	}

	c.mutexAgents.Lock()
	c.agents = pgAgents
	c.mutexAgents.Unlock()

	pgHosts, err := c.postgresService.GetHosts(
		ctx,
		&postgres.GetHostsFilter{},
	)
	if err != nil {
		return fmt.Errorf("error loading hosts: %s", err)
	}

	c.mutexHosts.Lock()
	c.hosts = pgHosts
	c.mutexHosts.Unlock()

	pgUsers, err := c.postgresService.GetUsers(
		ctx,
		&postgres.GetUsersFilter{},
	)
	if err != nil {
		return fmt.Errorf("error loading users: %s", err)
	}

	c.mutexUsers.Lock()
	c.users = pgUsers
	c.mutexUsers.Unlock()

	pgNotifiers, err := c.postgresService.GetNotifiers(
		ctx,
		&postgres.GetNotifiersFilter{
			Disabled: null.BoolFrom(false),
		},
	)
	if err != nil {
		return fmt.Errorf("error loading notifiers: %s", err)
	}

	c.mutexNotifiers.Lock()
	c.notifiers = pgNotifiers
	c.mutexNotifiers.Unlock()

	c.mutexChecks.Lock()
	c.mutexNotifier.Lock()
	defer c.mutexNotifier.Unlock()
	defer c.mutexChecks.Unlock()
	err = c.cacheService.DeleteAllHostStatuses()
	if err != nil {
		return fmt.Errorf("error clearing host statuses: %s", err)
	}

	defaultInterval := 60 * time.Second

	c.scheduler.Clear()

	for _, pgCheck := range c.checks {
		checkUID := pgCheck.UID

		pgChecker, err := c.getChecker(pgCheck.CheckerUID)
		if err != nil {
			c.log.Warnf("Error loading checker for check %s: %s", pgCheck.UID, err)

			continue
		}

		checkStatus := model.PgCheckStatusV1StatusUnkn
		checkMessage := ""
		if len(pgCheck.Statuses) > 0 {
			checkStatus = pgCheck.Statuses[0].Status
			checkMessage = pgCheck.Statuses[0].Message
		}

		err = c.cacheService.UpsertHostCheckStatus(
			pgCheck.HostUID,
			&model.ReHostStatusV1Check{
				CheckUID: pgCheck.UID,
				Status:   checkStatus,
				Message:  checkMessage,
			},
		)
		if err != nil {
			return fmt.Errorf("error caching check status: %s", err)
		}

		if pgCheck.Schedule.Valid {
			_, err = c.scheduler.CronWithSeconds(pgCheck.Schedule.String).Do(func() {
				err := c.check(checkUID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", checkUID, err)
				}
			})
			if err != nil {
				c.log.Warnf("Error scheduling check %s with own schedule: %s", pgCheck.UID, err)

				continue
			}
		} else if pgChecker.Capabilities.DefaultSchedule.Valid {
			_, err = c.scheduler.CronWithSeconds(pgChecker.Capabilities.DefaultSchedule.String).Do(func() {
				err := c.check(checkUID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", checkUID, err)
				}
			})
			if err != nil {
				c.log.Warnf("Error scheduling check %s with default checker schedule: %s", pgCheck.UID, err)

				continue
			}
		} else {
			_, err = c.scheduler.Every(defaultInterval).Do(func() {
				err := c.check(checkUID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", checkUID, err)
				}
			})
			if err != nil {
				c.log.Warnf("Error scheduling check %s with default schedule: %s", pgCheck.UID, err)

				continue
			}
		}
	}

	return nil
}

func (c *Controller) reloadLoop() error {
	reSystemEvents, subscription, err := c.cacheService.SubscribeForSystemEvents()
	if err != nil {
		return fmt.Errorf("subscribing for system events failed: %s", err)
	}
	defer subscription.Unsubscribe()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for !c.stop {
		select {
		case <-ticker.C:
			continue
		case reSystemEvent, more := <-reSystemEvents:
			if !more {
				return fmt.Errorf("subscription closed the channel")
			}

			switch reSystemEvent.Type {
			case model.ReSystemEventV1TypeCheckAdded,
				model.ReSystemEventV1TypeCheckUpdated,
				model.ReSystemEventV1TypeCheckDeleted,
				model.ReSystemEventV1TypeHostAdded,
				model.ReSystemEventV1TypeHostUpdated,
				model.ReSystemEventV1TypeHostDeleted,
				model.ReSystemEventV1TypeNotifierAdded,
				model.ReSystemEventV1TypeNotifierUpdated,
				model.ReSystemEventV1TypeNotifierDeleted:
				err := c.reload()
				if err != nil {
					c.log.Errorf("Error running reload: %s", err)

					continue
				}
			case model.ReSystemEventV1TypeCheckExecute:
				payload, ok := reSystemEvent.Payload.(*model.ReSystemEventV1CheckExecutePayload)
				if !ok {
					c.log.Errorf("Invalid payload for execute check event")

					continue
				}

				err := c.check(payload.CheckUID)
				if err != nil {
					c.log.Errorf("Error executing check: %s", err)

					continue
				}
			}
		}
	}

	return nil
}
