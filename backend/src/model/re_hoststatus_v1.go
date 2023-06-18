package model

type ReHostStatusV1Check struct {
	CheckUID string                `json:"check_uid"`
	Status   PgCheckStatusV1Status `json:"status"`
	Message  string                `json:"message"`
}

type ReHostStatusV1 struct {
	HostUID string                 `json:"host_uid"`
	Checks  []*ReHostStatusV1Check `json:"checks"`
}
