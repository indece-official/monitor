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

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/service/template"
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
			HostUID: null.StringFrom(hostUID),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error loading checks for host: %s", err)
	}

	hostStats := &HostStats{}

	for _, pgCheck := range pgChecks {
		if pgCheck.Status == nil {
			hostStats.CountUnknown++

			continue
		}

		switch pgCheck.Status.Status {
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
			CheckUID: null.StringFrom(checkUID),
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

	if pgCheck.Status == nil ||
		pgCheckStatus.Status != pgCheck.Status.Status {
		notify = true
	}

	err = c.postgresService.AddCheckStatus(ctx, pgCheckStatus)
	if err != nil {
		return fmt.Errorf("error adding check status: %s", err)
	}

	err = c.cacheService.UpsertHostCheckStatus(
		pgCheck.HostUID,
		&model.ReHostStatusV1Check{
			CheckUID: pgCheck.UID,
			Status:   pgCheckStatus.Status,
		},
	)
	if err != nil {
		return fmt.Errorf("error caching check status: %s", err)
	}

	if !notify {
		return nil
	}

	pgHosts, err := c.postgresService.GetHosts(
		ctx,
		&postgres.GetHostsFilter{
			HostUID: null.StringFrom(pgCheck.HostUID),
		},
	)
	if err != nil {
		return fmt.Errorf("error loading existing hosts: %s", err)
	}

	if len(pgHosts) != 1 {
		return fmt.Errorf("host not found")
	}

	pgHost := pgHosts[0]

	hostStats, err := c.getHostStats(ctx, pgHost.UID)
	if err != nil {
		return fmt.Errorf("error loading host stats: %s", err)
	}

	templateParams := map[string]string{}

	templateParams["host_name"] = pgHost.Name
	templateParams["count_ok"] = fmt.Sprintf("%d", hostStats.CountOK)
	templateParams["count_warning"] = fmt.Sprintf("%d", hostStats.CountWarning)
	templateParams["count_critical"] = fmt.Sprintf("%d", hostStats.CountCritical)
	templateParams["count_unknown"] = fmt.Sprintf("%d", hostStats.CountUnknown)
	templateParams["check_name"] = pgCheck.Name
	templateParams["checkstatus_status"] = string(pgCheckStatus.Status)
	templateParams["checkstatus_message"] = pgCheckStatus.Message

	for _, pgUser := range c.users {
		if !pgUser.Email.Valid {
			continue
		}

		err = c.sendEmail(
			ctx,
			model.LocaleEnUs,
			template.TemplateTypeStatusChanged,
			pgUser.Email.String,
			templateParams,
		)
		if err != nil {
			return fmt.Errorf("error sending email: %s", err)
		}
	}

	return nil
}
