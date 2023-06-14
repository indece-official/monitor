package model

import "time"

type ReUnnotifiedStatusChangedV1 struct {
	CheckUID        string    `json:"check_uid"`
	NotifierUID     string    `json:"notifier_uid"`
	DatetimeCreated time.Time `json:"datetime_created"`
}
