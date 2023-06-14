package cron

import (
	"fmt"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
)

func (c *Controller) notify() error {
	reUnnotifiedStatusChanges, err := c.cacheService.GetUnnotifiedStatusChanges()
	if err != nil {
		return fmt.Errorf("error loading unnotified status changes: %s", err)
	}

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

		pgHost, err := c.getHosts(reUnnotifiedStatusChange.CheckUID)
		if err != nil {
			c.log.Warnf("Error loading check for unnotified status change: %s", err)

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

		if isDue {
			// TODO: Recheck status
		}
	}

	return nil
}

func (c *Controller) notifyLoop() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	tickerNotify := time.NewTicker(5 * time.Second)
	defer tickerNotify.Stop()

	for !c.stop {
		select {
		case <-ticker.C:
			continue
		case <-tickerNotify.C:
			err := c.notify()
			if err != nil {
				return fmt.Errorf("error notifying: %s", err)
			}
		}
	}

	return nil
}
