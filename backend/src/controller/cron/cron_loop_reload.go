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
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"gopkg.in/guregu/null.v4"
)

const defaultInterval = 60 * time.Second

func (c *Controller) reloadChecker(checkerUID string) error {
	ctx := context.Background()

	pgCheckers, err := c.postgresService.GetCheckers(
		ctx,
		&postgres.GetCheckersFilter{
			CheckerUID: null.StringFrom(checkerUID),
		},
	)
	if err != nil {
		return fmt.Errorf("error loading checker: %s", err)
	}

	if len(pgCheckers) == 0 {
		return fmt.Errorf("error loading checker - not found")
	}

	newPgChecker := pgCheckers[0]

	c.mutexCheckers.Lock()
	defer c.mutexCheckers.Unlock()

	var oldPgChecker *model.PgCheckerV1

	newPgCheckers := []*model.PgCheckerV1{}
	for _, pgChecker := range c.checkers {
		if pgChecker.UID == checkerUID {
			oldPgChecker = pgChecker
			continue
		}

		newPgCheckers = append(newPgCheckers, pgChecker)
	}

	newPgCheckers = append(newPgCheckers, newPgChecker)

	c.checkers = newPgCheckers

	if oldPgChecker != nil &&
		!newPgChecker.Capabilities.DefaultSchedule.Equal(oldPgChecker.Capabilities.DefaultSchedule) {
		// Reschedule affected checks

		c.mutexChecks.Lock()
		defer c.mutexChecks.Unlock()
		for _, pgCheck := range c.checks {
			if pgCheck.CheckerUID != checkerUID {
				continue
			}

			if pgCheck.Schedule.Valid {
				// Customized config, no need to reschedule
				continue
			}

			err = c.rescheduleCheck(pgCheck, newPgChecker)
			if err != nil {
				c.log.Warnf("Error rescheduling check %s: %s", pgCheck.UID, err)

				continue
			}
		}
	}

	return nil
}

func (c *Controller) rescheduleCheck(pgCheck *model.PgCheckV1, pgChecker *model.PgCheckerV1) error {
	err := c.scheduler.RemoveByTag(fmt.Sprintf("check:%s", pgCheck.UID))
	if err != nil && !errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return fmt.Errorf("error unscheduling existing check: %s", err)
	}

	if pgCheck.Schedule.Valid {
		_, err = c.scheduler.
			CronWithSeconds(pgCheck.Schedule.String).
			Tag(fmt.Sprintf("check:%s", pgCheck.UID)).
			Tag(fmt.Sprintf("checker:%s", pgChecker.UID)).
			Do(func() {
				err := c.check(pgCheck.UID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", pgCheck.UID, err)
				}
			})
		if err != nil {
			return fmt.Errorf("error scheduling check %s with own schedule: %s", pgCheck.UID, err)
		}
	} else if pgChecker.Capabilities.DefaultSchedule.Valid {
		_, err = c.scheduler.
			CronWithSeconds(pgChecker.Capabilities.DefaultSchedule.String).
			Tag(fmt.Sprintf("check:%s", pgCheck.UID)).
			Tag(fmt.Sprintf("checker:%s", pgChecker.UID)).
			Do(func() {
				err := c.check(pgCheck.UID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", pgCheck.UID, err)
				}
			})
		if err != nil {
			return fmt.Errorf("error scheduling check %s with default checker schedule: %s", pgCheck.UID, err)
		}
	} else {
		_, err = c.scheduler.
			Every(defaultInterval).
			Tag(fmt.Sprintf("check:%s", pgCheck.UID)).
			Tag(fmt.Sprintf("checker:%s", pgChecker.UID)).
			Do(func() {
				err := c.check(pgCheck.UID)
				if err != nil {
					c.log.Errorf("Error running check %s: %s", pgCheck.UID, err)
				}
			})
		if err != nil {
			return fmt.Errorf("error scheduling check %s with default schedule: %s", pgCheck.UID, err)
		}
	}

	return nil
}

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
	err = c.cacheService.DeleteAllCheckStatuses()
	if err != nil {
		return fmt.Errorf("error clearing host statuses: %s", err)
	}

	c.scheduler.Clear()

	for _, pgCheck := range c.checks {
		pgChecker, err := c.getChecker(pgCheck.CheckerUID)
		if err != nil {
			c.log.Warnf("Error loading checker for check %s: %s", pgCheck.UID, err)

			continue
		}

		pgAgent, err := c.getAgent(pgChecker.AgentUID)
		if err != nil {
			c.log.Warnf("Error loading agent for checker %s: %s", pgChecker.UID, err)

			continue
		}

		pgCheckStatuses, err := c.postgresService.GetCheckStatuses(
			ctx,
			&postgres.GetCheckStatusesFilter{
				CheckUID:    pgCheck.UID,
				CountStatus: 1,
			},
		)
		if err != nil {
			c.log.Warnf("Error loading check statuses for check %s: %s", pgCheck.UID, err)

			continue
		}

		checkStatusUID := ""
		checkStatus := model.PgCheckStatusV1StatusUnkn
		checkMessage := ""
		checkDatetimeCreated := time.Now()
		if len(pgCheckStatuses) > 0 {
			checkStatusUID = pgCheckStatuses[0].UID
			checkStatus = pgCheckStatuses[0].Status
			checkMessage = pgCheckStatuses[0].Message
			checkDatetimeCreated = pgCheckStatuses[0].DatetimeCreated
		}

		err = c.cacheService.SetCheckStatus(
			pgCheck.UID,
			&model.ReCheckStatusV1{
				CheckStatusUID:  checkStatusUID,
				CheckUID:        pgCheck.UID,
				HostUID:         pgAgent.HostUID,
				Status:          checkStatus,
				Message:         checkMessage,
				DatetimeCreated: checkDatetimeCreated,
			},
		)
		if err != nil {
			return fmt.Errorf("error caching check status: %s", err)
		}

		err = c.rescheduleCheck(pgCheck, pgChecker)
		if err != nil {
			c.log.Warnf("Error scheduling check %s with own schedule: %s", pgCheck.UID, err)

			continue
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
				model.ReSystemEventV1TypeCheckerDeleted,
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
			case model.ReSystemEventV1TypeCheckerAdded:
				payload, ok := reSystemEvent.Payload.(*model.ReSystemEventV1CheckerAddedPayload)
				if !ok {
					c.log.Errorf("Invalid payload for checker added event: %s", err)

					continue
				}

				err := c.reloadChecker(payload.CheckerUID)
				if err != nil {
					c.log.Errorf("Error running reload checker: %s", err)

					continue
				}
			case model.ReSystemEventV1TypeCheckerUpdated:
				payload, ok := reSystemEvent.Payload.(*model.ReSystemEventV1CheckerUpdatedPayload)
				if !ok {
					c.log.Errorf("Invalid payload for checker updated event: %s", err)

					continue
				}

				err := c.reloadChecker(payload.CheckerUID)
				if err != nil {
					c.log.Errorf("Error running reload checker: %s", err)

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
