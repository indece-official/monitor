package cron

import (
	"context"
	"testing"
	"time"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/stretchr/testify/assert"
)

func TestFilterStatusCriticalOnly(t *testing.T) {
	controller := &Controller{}

	pgFilter := &model.PgNotifierV1ConfigFilter{}
	pgFilter.Critical = true
	pgFilter.Warning = false
	pgFilter.Unknown = false
	pgFilter.Decline = false

	matches := controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusWarn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusUnkn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusOK,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusWarn,
		model.PgCheckStatusV1StatusOK,
	)
	assert.False(t, matches)
}

func TestFilterStatusCriticalWithDecline(t *testing.T) {
	controller := &Controller{}

	pgFilter := &model.PgNotifierV1ConfigFilter{}
	pgFilter.Critical = true
	pgFilter.Warning = false
	pgFilter.Unknown = false
	pgFilter.Decline = true

	matches := controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusWarn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusUnkn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusOK,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusWarn,
		model.PgCheckStatusV1StatusOK,
	)
	assert.False(t, matches)
}

func TestFilterStatusCriticalWarningOnly(t *testing.T) {
	controller := &Controller{}

	pgFilter := &model.PgNotifierV1ConfigFilter{}
	pgFilter.Critical = true
	pgFilter.Warning = true
	pgFilter.Unknown = false
	pgFilter.Decline = false

	matches := controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusWarn,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusUnkn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusOK,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusWarn,
		model.PgCheckStatusV1StatusOK,
	)
	assert.False(t, matches)
}

func TestFilterStatusCriticalWarningWithDecline(t *testing.T) {
	controller := &Controller{}

	pgFilter := &model.PgNotifierV1ConfigFilter{}
	pgFilter.Critical = true
	pgFilter.Warning = true
	pgFilter.Unknown = false
	pgFilter.Decline = true

	matches := controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusCrit,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusWarn,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusOK,
		model.PgCheckStatusV1StatusUnkn,
	)
	assert.False(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusCrit,
		model.PgCheckStatusV1StatusOK,
	)
	assert.True(t, matches)

	matches = controller.filterStatus(
		pgFilter,
		model.PgCheckStatusV1StatusWarn,
		model.PgCheckStatusV1StatusOK,
	)
	assert.True(t, matches)
}

func TestCheckNotificationDueIsDue(t *testing.T) {
	controller := Controller{}

	controller.notifiers = []*model.PgNotifierV1{
		{
			UID: "ntf-01",
			Config: &model.PgNotifierV1Config{
				Filters: []*model.PgNotifierV1ConfigFilter{
					{
						TagUIDs: []string{
							"tag-01",
						},
						MinDuration: 30 * time.Second,
					},
					{
						TagUIDs: []string{
							"tag-01",
						},
						MinDuration: 60 * time.Second,
					},
					{
						TagUIDs: []string{
							"tag-02",
						},
						MinDuration: 10 * time.Second,
					},
					{
						TagUIDs:     []string{},
						MinDuration: 40 * time.Second,
					},
				},
			},
		},
	}

	controller.hosts = []*model.PgHostV1{
		{
			UID: "hst-01",
			TagUIDs: []string{
				"tag-01",
			},
		},
	}

	reNotification := &model.ReNotificationV1{
		CheckUID:        "chk-01",
		NotifierUID:     "ntf-01",
		HostUID:         "hst-01",
		DatetimeCreated: time.Now().Add(-50 * time.Second),
	}

	ctx := context.Background()

	isDue, matchingPgFilters, err := controller.checkNotificationDue(ctx, reNotification)
	assert.NoError(t, err)
	assert.True(t, isDue)
	assert.Len(t, matchingPgFilters, 2)
	assert.Len(t, matchingPgFilters[0].TagUIDs, 1)
	assert.Equal(t, "tag-01", matchingPgFilters[0].TagUIDs[0])
	assert.Len(t, matchingPgFilters[1].TagUIDs, 0)
}

func TestCheckNotificationDueNotDue(t *testing.T) {
	controller := Controller{}

	controller.notifiers = []*model.PgNotifierV1{
		{
			UID: "ntf-01",
			Config: &model.PgNotifierV1Config{
				Filters: []*model.PgNotifierV1ConfigFilter{
					{
						TagUIDs: []string{
							"tag-01",
						},
						MinDuration: 30 * time.Second,
					},
					{
						TagUIDs: []string{
							"tag-01",
						},
						MinDuration: 60 * time.Second,
					},
					{
						TagUIDs: []string{
							"tag-02",
						},
						MinDuration: 10 * time.Second,
					},
					{
						TagUIDs:     []string{},
						MinDuration: 40 * time.Second,
					},
				},
			},
		},
	}

	controller.hosts = []*model.PgHostV1{
		{
			UID: "hst-01",
			TagUIDs: []string{
				"tag-01",
			},
		},
	}

	reNotification := &model.ReNotificationV1{
		CheckUID:        "chk-01",
		NotifierUID:     "ntf-01",
		HostUID:         "hst-01",
		DatetimeCreated: time.Now().Add(-20 * time.Second),
	}

	ctx := context.Background()

	isDue, matchingPgFilters, err := controller.checkNotificationDue(ctx, reNotification)
	assert.NoError(t, err)
	assert.False(t, isDue)
	assert.Len(t, matchingPgFilters, 0)
}
