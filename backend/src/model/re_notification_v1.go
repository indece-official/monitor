package model

import "time"

type ReNotificationV1 struct {
	HostUID         string                `json:"host_uid"`
	CheckUID        string                `json:"check_uid"`
	NotifierUID     string                `json:"notifier_uid"`
	Status          PgCheckStatusV1Status `json:"status"`
	PreviousStatus  PgCheckStatusV1Status `json:"previous_status"`
	DatetimeCreated time.Time             `json:"datetime_created"`
}
