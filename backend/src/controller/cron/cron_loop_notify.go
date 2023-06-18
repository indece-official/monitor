package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/template"
)

func (c *Controller) notify(ctx context.Context) error {
	reUnnotifiedStatusChanges, err := c.cacheService.GetUnnotifiedStatusChanges()
	if err != nil {
		return fmt.Errorf("error loading unnotified status changes: %s", err)
	}

	relevantReUnnotifiedStatusChanges := []*model.ReUnnotifiedStatusChangeV1{}

	for _, reUnnotifiedStatusChange := range reUnnotifiedStatusChanges {
		pgNotifier, err := c.getNotifier(reUnnotifiedStatusChange.NotifierUID)
		if err != nil {
			c.log.Warnf("Error loading notifier for unnotified status change: %s", err)

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		pgCheck, err := c.getCheck(reUnnotifiedStatusChange.CheckUID)
		if err != nil {
			c.log.Warnf("Error loading check for unnotified status change: %s", err)

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		pgHost, err := c.getHost(reUnnotifiedStatusChange.HostUID)
		if err != nil {
			c.log.Warnf("Error loading host for unnotified status change: %s", err)

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		filterMatches := false
		isDue := false

		for _, pgFilter := range pgNotifier.Config.Filters {
			tagsMatch := true

			for _, tagUID := range pgFilter.TagUIDs {
				if !gousu.ContainsString(pgHost.TagUIDs, tagUID) {
					tagsMatch = false
					break
				}
			}

			// TODO: Check for status

			if tagsMatch {
				filterMatches = true

				if reUnnotifiedStatusChange.DatetimeCreated.Add(pgFilter.MinDuration).Before(time.Now()) {
					isDue = true
				}
			}
		}

		if !filterMatches {
			c.log.Warnf("No filter matches for unnotified status change")

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		if !isDue {
			continue
		}

		// Recheck status

		reHostStatus, err := c.cacheService.GetHostStatus(pgCheck.HostUID)
		if err != nil {
			c.log.Warnf("Error loading host status: %s", err)

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		if reHostStatus == nil {
			c.log.Warnf("Error no host status found: %s", err)

			err = c.cacheService.DeleteUnnotifiedStatusChange(reUnnotifiedStatusChange.CheckUID, reUnnotifiedStatusChange.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting unnotified status change: %s", err)
			}

			continue
		}

		statusMatches := false

		for _, reCheckStatus := range reHostStatus.Checks {
			if reCheckStatus.CheckUID == pgCheck.UID {
				// TODO: warn -> crit etc.
				statusMatches = reCheckStatus.Status == reUnnotifiedStatusChange.Status
				break
			}
		}

		if statusMatches {
			relevantReUnnotifiedStatusChanges = append(relevantReUnnotifiedStatusChanges, reUnnotifiedStatusChange)
		}
	}

	if len(relevantReUnnotifiedStatusChanges) > 0 {
		// Group by host & notifier
		changesMap := map[string][]*model.ReUnnotifiedStatusChangeV1{}
		for _, change := range relevantReUnnotifiedStatusChanges {
			key := fmt.Sprintf("%s:%s", change.HostUID, change.NotifierUID)
			_, ok := changesMap[key]
			if !ok {
				changesMap[key] = []*model.ReUnnotifiedStatusChangeV1{}
			}

			changesMap[key] = append(changesMap[key], change)
		}

		for _, changes := range changesMap {
			pgHost, err := c.getHost(changes[0].HostUID)
			if err != nil {
				return fmt.Errorf("error loading host: %s", err)
			}

			templateParamChecks := []map[string]string{}
			for _, change := range changes {
				pgCheck, err := c.getCheck(change.CheckUID)
				if err != nil {
					return fmt.Errorf("error loading check: %s", err)
				}

				reHostStatus, err := c.cacheService.GetHostStatus(pgCheck.HostUID)
				if err != nil {
					c.log.Warnf("Error loading host status: %s", err)

					continue
				}

				if reHostStatus == nil {
					c.log.Warnf("Error no host status found: %s", err)

					continue
				}

				var matchingReHostCheck *model.ReHostStatusV1Check

				for _, reCheckStatus := range reHostStatus.Checks {
					if reCheckStatus.CheckUID == pgCheck.UID {
						matchingReHostCheck = reCheckStatus
						break
					}
				}

				if matchingReHostCheck == nil {
					c.log.Warnf("Error no check status found: %s", err)

					continue
				}

				templateParamChecks = append(templateParamChecks, map[string]string{
					"check_name":          pgCheck.Name,
					"checkstatus_status":  string(matchingReHostCheck.Status),
					"checkstatus_message": matchingReHostCheck.Message,
				})
			}

			hostStats, err := c.getHostStats(ctx, pgHost.UID)
			if err != nil {
				return fmt.Errorf("error loading host stats: %s", err)
			}

			templateParams := map[string]interface{}{}

			templateParams["host_name"] = pgHost.Name
			templateParams["count_ok"] = fmt.Sprintf("%d", hostStats.CountOK)
			templateParams["count_warning"] = fmt.Sprintf("%d", hostStats.CountWarning)
			templateParams["count_critical"] = fmt.Sprintf("%d", hostStats.CountCritical)
			templateParams["count_unknown"] = fmt.Sprintf("%d", hostStats.CountUnknown)
			templateParams["checks"] = templateParamChecks

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
		}
	}

	return nil
}

func (c *Controller) notifyLoop() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	tickerNotify := time.NewTicker(10 * time.Second)
	defer tickerNotify.Stop()

	for !c.stop {
		select {
		case <-ticker.C:
			continue
		case <-tickerNotify.C:
			ctx := context.Background()

			err := c.notify(ctx)
			if err != nil {
				return fmt.Errorf("error notifying: %s", err)
			}
		}
	}

	return nil
}
