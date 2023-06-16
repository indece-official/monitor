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
			Disabled: null.BoolFrom(false),
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

	pgConnectors, err := c.postgresService.GetConnectors(
		ctx,
		&postgres.GetConnectorsFilter{},
	)
	if err != nil {
		return fmt.Errorf("error loading connectors: %s", err)
	}

	c.mutexConnectors.Lock()
	c.connectors = pgConnectors
	c.mutexConnectors.Unlock()

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

	err = c.cacheService.DeleteAllHostStatuses()
	if err != nil {
		return fmt.Errorf("error clearing host statuses: %s", err)
	}

	defaultInterval := 60 * time.Second

	c.scheduler.Clear()

	c.mutexChecks.Lock()
	for _, pgCheck := range c.checks {
		checkUID := pgCheck.UID

		pgChecker, err := c.getChecker(pgCheck.CheckerUID)
		if err != nil {
			c.log.Warnf("Error loading checker for check %s: %s", pgCheck.UID, err)

			continue
		}

		checkStatus := model.PgCheckStatusV1StatusUnkn
		if pgCheck.Status != nil {
			checkStatus = pgCheck.Status.Status
		}

		err = c.cacheService.UpsertHostCheckStatus(
			pgCheck.HostUID,
			&model.ReHostStatusV1Check{
				CheckUID: pgCheck.UID,
				Status:   checkStatus,
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
	c.mutexChecks.Unlock()

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
				model.ReSystemEventV1TypeCheckDeleted:
				err := c.reload()
				if err != nil {
					c.log.Errorf("Error running checks: %s", err)

					continue
				}
			}
		}
	}

	return nil
}
