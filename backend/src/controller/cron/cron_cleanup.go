package cron

import (
	"context"
	"time"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/postgres"
)

func (c *Controller) cleanup() {
	qctx := context.Background()

	c.log.Infof("Running cleanup task ...")

	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		qctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyHistoryMaxAge,
			},
		},
	)
	if err != nil {
		c.log.Errorf("Error loading history max age from config: %s", err)

		return
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyHistoryMaxAge] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyHistoryMaxAge].Value == "" {
		c.log.Warnf("No history max age configured - not running cleanup")

		return
	}

	maxAge, err := time.ParseDuration(pgConfigProperties[model.PgConfigPropertyV1KeyHistoryMaxAge].Value)
	if err != nil {
		c.log.Errorf("Invalid value for history max age in config: %s", err)

		return
	}

	if maxAge <= 0 {
		c.log.Warnf("Configured history max age is <= 0 - not running cleanup")

		return
	}

	err = c.postgresService.DeleteCheckStatusByAge(qctx, maxAge)
	if err != nil {
		c.log.Errorf("Error deleting old status checks from postgres: %s", err)

		return
	}

	c.log.Infof("Finished running cleanup task")
}

func (c *Controller) cleanupLoop() error {
	tickerKeepAlive := time.NewTicker(1 * time.Second)
	defer tickerKeepAlive.Stop()

	tickerTrigger := time.NewTicker(24 * time.Hour)
	defer tickerTrigger.Stop()

	for !c.stop {
		select {
		case <-tickerKeepAlive.C:
			continue
		case <-tickerTrigger.C:
			c.cleanup()
		}
	}

	return nil
}
