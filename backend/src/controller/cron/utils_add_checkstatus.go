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

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/utils"
	"gopkg.in/guregu/null.v4"
)

type HostStats struct {
	CountOK       int
	CountWarning  int
	CountCritical int
	CountUnknown  int
}

func (c *Controller) getHostStats(
	ctx context.Context,
	hostUID string,
) (*HostStats, error) {
	pgChecks, err := c.postgresService.GetChecks(
		ctx,
		&postgres.GetChecksFilter{
			HostUID:     null.StringFrom(hostUID),
			CountStatus: 1,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error loading checks for host: %s", err)
	}

	hostStats := &HostStats{}

	for _, pgCheck := range pgChecks {
		if len(pgCheck.Statuses) == 0 {
			hostStats.CountUnknown++

			continue
		}

		switch pgCheck.Statuses[0].Status {
		case model.PgCheckStatusV1StatusCrit:
			hostStats.CountCritical++
		case model.PgCheckStatusV1StatusWarn:
			hostStats.CountWarning++
		case model.PgCheckStatusV1StatusOK:
			hostStats.CountOK++
		default:
			hostStats.CountUnknown++
		}
	}

	return hostStats, nil
}

func (c *Controller) addCheckStatus(
	ctx context.Context,
	checkUID string,
	message string,
	checkError error,
	values map[string]string,
) error {
	var err error

	pgChecks, err := c.postgresService.GetChecks(
		ctx,
		&postgres.GetChecksFilter{
			CheckUID:    null.StringFrom(checkUID),
			CountStatus: 1,
		},
	)
	if err != nil {
		return fmt.Errorf("error loading existing check status: %s", err)
	}

	if len(pgChecks) != 1 {
		return fmt.Errorf("check not found")
	}

	pgCheck := pgChecks[0]

	pgChecker, err := c.getChecker(pgCheck.CheckerUID)
	if err != nil {
		return fmt.Errorf("checker not found for check")
	}

	pgAgent, err := c.getAgent(pgChecker.AgentUID)
	if err != nil {
		return fmt.Errorf("agent not found for checker")
	}

	pgCheckStatus := &model.PgCheckStatusV1{}
	pgCheckStatus.UID, err = utils.UUID()
	if err != nil {
		return fmt.Errorf("error generating uid for check status: %s", err)
	}
	pgCheckStatus.CheckUID = checkUID
	pgCheckStatus.Data = map[string]interface{}{}

	// Evaluate status

	isCrit := false
	isWarn := false

	for _, pgCheckerValue := range pgChecker.Capabilities.Values {
		valueStr, ok := values[pgCheckerValue.Name]
		if !ok {
			// TODO: Error?

			continue
		}

		value, err := model.NewValue(
			pgCheckerValue.Type,
			valueStr,
		)
		if err != nil {
			c.log.Warnf("Error parsing value %s: %s", pgCheckerValue.Name, err)

			// TODO: Error?

			continue
		}

		if pgCheckerValue.MinCrit.Valid {
			minCrit, err := model.NewValue(
				pgCheckerValue.Type,
				pgCheckerValue.MinCrit.String,
			)
			if err != nil {
				c.log.Warnf("Error parsing min-crit: %s", err)

				// TODO: Error?

				continue
			}

			isLessThan, err := value.LessThanEqual(minCrit)
			if err != nil {
				c.log.Warnf("Error comparing min-crit: %s", err)

				// TODO: Error?

				continue
			}

			if isLessThan {
				isCrit = true
			}
		}

		if pgCheckerValue.MinWarn.Valid {
			minWarn, err := model.NewValue(
				pgCheckerValue.Type,
				pgCheckerValue.MinWarn.String,
			)
			if err != nil {
				c.log.Warnf("Error parsing min-warn: %s", err)

				// TODO: Error?

				continue
			}

			isLessThan, err := value.LessThanEqual(minWarn)
			if err != nil {
				c.log.Warnf("Error comparing min-warn: %s", err)

				// TODO: Error?

				continue
			}

			if isLessThan {
				isWarn = true
			}
		}

		if pgCheckerValue.MaxCrit.Valid {
			maxCrit, err := model.NewValue(
				pgCheckerValue.Type,
				pgCheckerValue.MaxCrit.String,
			)
			if err != nil {
				c.log.Warnf("Error parsing max-crit: %s", err)

				// TODO: Error?

				continue
			}

			isGreaterThan, err := value.GreaterThanEqual(maxCrit)
			if err != nil {
				c.log.Warnf("Error comparing max-crit: %s", err)

				// TODO: Error?

				continue
			}

			if isGreaterThan {
				isCrit = true
			}
		}

		if pgCheckerValue.MaxWarn.Valid {
			maxWarn, err := model.NewValue(
				pgCheckerValue.Type,
				pgCheckerValue.MaxWarn.String,
			)
			if err != nil {
				c.log.Warnf("Error parsing max-warn: %s", err)

				// TODO: Error?

				continue
			}

			isGreaterThan, err := value.GreaterThanEqual(maxWarn)
			if err != nil {
				c.log.Warnf("Error comparing min-warn: %s", err)

				// TODO: Error?

				continue
			}

			if isGreaterThan {
				isWarn = true
			}
		}

		pgCheckStatus.Data[pgCheckerValue.Name] = value.Raw()
	}

	if checkError != nil {
		pgCheckStatus.Status = model.PgCheckStatusV1StatusCrit
		pgCheckStatus.Message = checkError.Error()
	} else {
		switch {
		case isCrit:
			pgCheckStatus.Status = model.PgCheckStatusV1StatusCrit
		case isWarn:
			pgCheckStatus.Status = model.PgCheckStatusV1StatusWarn
		default:
			pgCheckStatus.Status = model.PgCheckStatusV1StatusOK
		}
		pgCheckStatus.Message = message
	}

	notify := false
	prevStatus := model.PgCheckStatusV1StatusUnkn
	if len(pgCheck.Statuses) > 0 {
		prevStatus = pgCheck.Statuses[0].Status
	}

	if pgCheckStatus.Status != prevStatus {
		notify = true
	}

	err = c.postgresService.AddCheckStatus(ctx, pgCheckStatus)
	if err != nil {
		return fmt.Errorf("error adding check status: %s", err)
	}

	err = c.cacheService.UpsertHostCheckStatus(
		pgAgent.HostUID,
		&model.ReHostStatusV1Check{
			CheckUID: pgCheck.UID,
			Status:   pgCheckStatus.Status,
			Message:  pgCheckStatus.Message,
		},
	)
	if err != nil {
		return fmt.Errorf("error caching check status: %s", err)
	}

	if !notify {
		return nil
	}

	// Check if maintenance active
	pgMaintenances, err := c.postgresService.GetMaintenances(
		ctx,
		&postgres.GetMaintenancesFilter{
			NowActive: true,
		},
	)
	if err != nil {
		return fmt.Errorf("error loading active maintenances: %s", err)
	}

	inMaintenance := false

	for _, pgMaintenance := range pgMaintenances {
		if len(pgMaintenance.Details.Affected.CheckUIDs) > 0 &&
			gousu.ContainsString(pgMaintenance.Details.Affected.CheckUIDs, pgCheck.UID) {
			inMaintenance = true
			break
		}

		if len(pgMaintenance.Details.Affected.HostUIDs) > 0 &&
			gousu.ContainsString(pgMaintenance.Details.Affected.HostUIDs, pgAgent.UID) {
			inMaintenance = true
			break
		}

		if len(pgMaintenance.Details.Affected.TagUIDs) > 0 {
			pgHost, err := c.getHost(pgAgent.HostUID)
			if err != nil {
				return fmt.Errorf("error loading host: %s", err)
			}

			if utils.Intersects[string](pgMaintenance.Details.Affected.TagUIDs, pgHost.TagUIDs) {
				inMaintenance = true
				break
			}
		}
	}

	if inMaintenance {
		return nil
	}

	for _, pgNotifier := range c.notifiers {
		err = c.cacheService.SetNotification(
			&model.ReNotificationV1{
				HostUID:         pgAgent.HostUID,
				NotifierUID:     pgNotifier.UID,
				CheckUID:        pgCheck.UID,
				Status:          pgCheckStatus.Status,
				PreviousStatus:  prevStatus,
				DatetimeCreated: time.Now(),
			},
		)
		if err != nil {
			return fmt.Errorf("error adding unnotified status change: %s", err)
		}
	}

	return nil
}
