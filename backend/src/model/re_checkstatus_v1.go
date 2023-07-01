package model

import "time"

type ReCheckStatusV1 struct {
	CheckStatusUID  string                 `json:"checkstatus_uid"`
	CheckUID        string                 `json:"check_uid"`
	HostUID         string                 `json:"host_uid"`
	Status          PgCheckStatusV1Status  `json:"status"`
	Message         string                 `json:"message"`
	Data            map[string]interface{} `json:"data"`
	DatetimeCreated time.Time              `json:"datetime_created"`
}
