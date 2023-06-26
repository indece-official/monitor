package model

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type PgMaintenanceV1DetailsAffected struct {
	HostUIDs  []string `json:"host_uids"`
	CheckUIDs []string `json:"check_uids"`
	TagUIDs   []string `json:"tag_uids"`
}

type PgMaintenanceV1Details struct {
	Affected *PgMaintenanceV1DetailsAffected `json:"affected"`
}

type PgMaintenanceV1 struct {
	UID             string
	Title           string
	Message         string
	Details         *PgMaintenanceV1Details
	DatetimeCreated time.Time
	DatetimeUpdated time.Time
	DatetimeStart   time.Time
	DatetimeFinish  null.Time
}
