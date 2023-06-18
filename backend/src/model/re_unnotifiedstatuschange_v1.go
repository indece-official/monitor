package model

import "time"

type ReUnnotifiedStatusChangeV1 struct {
	HostUID         string                `json:"host_uid"`
	CheckUID        string                `json:"check_uid"`
	NotifierUID     string                `json:"notifier_uid"`
	Status          PgCheckStatusV1Status `json:"status"`
	DatetimeCreated time.Time             `json:"datetime_created"`
}
