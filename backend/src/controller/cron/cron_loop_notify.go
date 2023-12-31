package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/template"
)

func (c *Controller) sendNotifications(
	ctx context.Context,
	reNotifications []*model.ReNotificationV1,
) error {
	c.log.Infof("Sending notifications (%d) for host %s", len(reNotifications), reNotifications[0].HostUID)

	pgHost, err := c.getHost(reNotifications[0].HostUID)
	if err != nil {
		return fmt.Errorf("error loading host: %s", err)
	}

	pgNotifier, err := c.getNotifier(reNotifications[0].NotifierUID)
	if err != nil {
		return fmt.Errorf("error loading notifier: %s", err)
	}

	templateParamChecks := []map[string]string{}
	for _, change := range reNotifications {
		pgCheck, err := c.getCheck(change.CheckUID)
		if err != nil {
			return fmt.Errorf("error loading check: %s", err)
		}

		reCheckStatus, err := c.cacheService.GetCheckStatus(pgCheck.UID)
		if err != nil {
			c.log.Warnf("Error loading host status check: %s", err)

			continue
		}

		if reCheckStatus == nil {
			c.log.Warnf("Error no host status check found")

			continue
		}

		templateParamChecks = append(templateParamChecks, map[string]string{
			"check_name":          pgCheck.Name,
			"checkstatus_status":  string(reCheckStatus.Status),
			"checkstatus_message": reCheckStatus.Message,
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

	switch pgNotifier.Type {
	case model.PgNotifierV1TypeEmailSmtp:
		sender := c.smtpService.Open(pgNotifier.Config.Params.EmailSmtp)
		defer sender.Close()

		for _, pgUser := range c.users {
			if !pgUser.Email.Valid {
				continue
			}

			err = c.sendEmail(
				ctx,
				sender,
				model.LocaleEnUs,
				template.TemplateTypeStatusChanged,
				pgUser.Email.String,
				templateParams,
			)
			if err != nil {
				return fmt.Errorf("error sending email: %s", err)
			}
		}
	case model.PgNotifierV1TypeHttp:
		err = c.sendHttpQuery(
			ctx,
			model.LocaleEnUs,
			template.TemplateTypeStatusChanged,
			templateParams,
			pgNotifier.Config.Params.Http.Method,
			pgNotifier.Config.Params.Http.URL,
			pgNotifier.Config.Params.Http.Headers,
			pgNotifier.Config.Params.Http.Body,
		)
		if err != nil {
			return fmt.Errorf("error sending teams notification: %s", err)
		}

	case model.PgNotifierV1TypeMicrosoftTeams:
		err = c.sendTeamsNotification(
			ctx,
			model.LocaleEnUs,
			pgNotifier.Config.Params.MicrosoftTeams.WebhookURL,
			template.TemplateTypeStatusChanged,
			templateParams,
		)
		if err != nil {
			return fmt.Errorf("error sending teams notification: %s", err)
		}
	}

	return nil
}

func (c *Controller) groupNotificatinsByHostAndNotifier(reNotifications []*model.ReNotificationV1) [][]*model.ReNotificationV1 {
	changesMap := map[string][]*model.ReNotificationV1{}
	for _, change := range reNotifications {
		key := fmt.Sprintf("%s:%s", change.HostUID, change.NotifierUID)
		_, ok := changesMap[key]
		if !ok {
			changesMap[key] = []*model.ReNotificationV1{}
		}

		changesMap[key] = append(changesMap[key], change)
	}

	changesArr := [][]*model.ReNotificationV1{}
	for _, changes := range changesMap {
		changesArr = append(changesArr, changes)
	}

	return changesArr
}

func (c *Controller) filterStatus(
	pgFilter *model.PgNotifierV1ConfigFilter,
	statusOld model.PgCheckStatusV1Status,
	statusNew model.PgCheckStatusV1Status,
) bool {
	if statusOld == statusNew {
		return false
	}

	if pgFilter.Critical && statusNew == model.PgCheckStatusV1StatusCrit {
		return true
	}

	if pgFilter.Warning && statusNew == model.PgCheckStatusV1StatusWarn {
		return true
	}

	if pgFilter.Unknown && statusNew == model.PgCheckStatusV1StatusUnkn {
		return true
	}

	if pgFilter.Decline && statusNew == model.PgCheckStatusV1StatusOK {
		if pgFilter.Critical && statusOld == model.PgCheckStatusV1StatusCrit {
			return true
		}

		if pgFilter.Warning && statusOld == model.PgCheckStatusV1StatusWarn {
			return true
		}

		if pgFilter.Unknown && statusOld == model.PgCheckStatusV1StatusUnkn {
			return true
		}
	}

	return false
}

func (c *Controller) checkNotificationDue(ctx context.Context, reNotification *model.ReNotificationV1) (bool, bool, []*model.PgNotifierV1ConfigFilter, error) {
	matchingPgFilters := []*model.PgNotifierV1ConfigFilter{}

	pgNotifier, err := c.getNotifier(reNotification.NotifierUID)
	if err != nil {
		return false, false, nil, fmt.Errorf("error loading notifier for unnotified status change: %s", err)
	}

	pgHost, err := c.getHost(reNotification.HostUID)
	if err != nil {
		return false, false, nil, fmt.Errorf("error loading host for unnotified status change: %s", err)
	}

	filterMatch := false

	for _, pgFilter := range pgNotifier.Config.Filters {
		tagsMatch := true

		for _, tagUID := range pgFilter.TagUIDs {
			if !gousu.ContainsString(pgHost.TagUIDs, tagUID) {
				tagsMatch = false
				break
			}
		}

		if tagsMatch {
			filterMatch = true

			if reNotification.DatetimeCreated.Add(pgFilter.MinDuration).Before(time.Now()) {
				matchingPgFilters = append(matchingPgFilters, pgFilter)
			}
		}
	}

	if filterMatch {
		return true, len(matchingPgFilters) > 0, matchingPgFilters, nil
	}

	return false, false, nil, nil
}

func (c *Controller) notify(ctx context.Context) error {
	c.mutexNotifier.Lock()
	defer c.mutexNotifier.Unlock()

	reNotifications, err := c.cacheService.GetNotifications()
	if err != nil {
		return fmt.Errorf("error loading notifications: %s", err)
	}

	relevantReNotifications := []*model.ReNotificationV1{}

	for _, reNotification := range reNotifications {
		hasMatches, isDue, matchingPgFilters, err := c.checkNotificationDue(ctx, reNotification)
		if err != nil {
			c.log.Warnf("Error checking if notification is due: %s", err)

			err = c.cacheService.DeleteNotification(reNotification.CheckUID, reNotification.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting notification: %s", err)
			}

			continue
		}

		if !hasMatches {
			err = c.cacheService.DeleteNotification(reNotification.CheckUID, reNotification.NotifierUID)
			if err != nil {
				return fmt.Errorf("error deleting notification: %s", err)
			}

			continue
		}

		if !isDue {
			continue
		}

		// Recheck status
		err = c.cacheService.DeleteNotification(reNotification.CheckUID, reNotification.NotifierUID)
		if err != nil {
			return fmt.Errorf("error deleting sent notification: %s", err)
		}

		reCheckStatus, err := c.cacheService.GetCheckStatus(reNotification.CheckUID)
		if err != nil {
			c.log.Warnf("Error loading host status: %s", err)

			continue
		}

		if reCheckStatus == nil {
			c.log.Warnf("Error no host status found: %s", err)

			continue
		}

		for _, pgFilter := range matchingPgFilters {
			if c.filterStatus(pgFilter, reNotification.PreviousStatus, reCheckStatus.Status) {
				relevantReNotifications = append(relevantReNotifications, reNotification)
				break
			}
		}
	}

	if len(relevantReNotifications) > 0 {
		// Group by host & notifier
		changesArr := c.groupNotificatinsByHostAndNotifier(relevantReNotifications)

		for _, changes := range changesArr {
			err = c.sendNotifications(ctx, changes)
			if err != nil {
				return fmt.Errorf("error loading check: %s", err)
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
