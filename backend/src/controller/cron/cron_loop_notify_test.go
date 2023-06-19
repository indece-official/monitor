package cron

import (
	"testing"

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
