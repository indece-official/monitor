// indece Monitor
// Copyright (C) 2023 indece UG (haftungsbeschränkt)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License or any
// later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"encoding/json"
	"fmt"
)

type ReSystemEventV1Type string

const (
	ReSystemEventV1TypeCheckAdded      ReSystemEventV1Type = "check_added"
	ReSystemEventV1TypeCheckUpdated    ReSystemEventV1Type = "check_updated"
	ReSystemEventV1TypeCheckDeleted    ReSystemEventV1Type = "check_deleted"
	ReSystemEventV1TypeHostAdded       ReSystemEventV1Type = "host_added"
	ReSystemEventV1TypeHostUpdated     ReSystemEventV1Type = "host_updated"
	ReSystemEventV1TypeHostDeleted     ReSystemEventV1Type = "host_deleted"
	ReSystemEventV1TypeTagAdded        ReSystemEventV1Type = "tag_added"
	ReSystemEventV1TypeTagUpdated      ReSystemEventV1Type = "tag_updated"
	ReSystemEventV1TypeTagDeleted      ReSystemEventV1Type = "tag_deleted"
	ReSystemEventV1TypeNotifierAdded   ReSystemEventV1Type = "notifier_added"
	ReSystemEventV1TypeNotifierUpdated ReSystemEventV1Type = "notifier_updated"
	ReSystemEventV1TypeNotifierDeleted ReSystemEventV1Type = "notifier_deleted"
	ReSystemEventV1TypeCertsUpdated    ReSystemEventV1Type = "certs_updates"
)

type ReSystemEventV1CheckAddedPayload struct {
	CheckUID string `json:"check_uid"`
}

type ReSystemEventV1CheckUpdatedPayload struct {
	CheckUID string `json:"check_uid"`
}

type ReSystemEventV1CheckDeletedPayload struct {
	CheckUID string `json:"check_uid"`
}

type ReSystemEventV1HostAddedPayload struct {
	HostUID string `json:"host_uid"`
}

type ReSystemEventV1HostUpdatedPayload struct {
	HostUID string `json:"host_uid"`
}

type ReSystemEventV1HostDeletedPayload struct {
	HostUID string `json:"host_uid"`
}

type ReSystemEventV1TagAddedPayload struct {
	TagUID string `json:"tag_uid"`
}

type ReSystemEventV1TagUpdatedPayload struct {
	TagUID string `json:"tag_uid"`
}

type ReSystemEventV1TagDeletedPayload struct {
	TagUID string `json:"tag_uid"`
}

type ReSystemEventV1NotifierAddedPayload struct {
	NotifierUID string `json:"notifier_uid"`
}

type ReSystemEventV1NotifierUpdatedPayload struct {
	NotifierUID string `json:"notifier_uid"`
}

type ReSystemEventV1NotifierDeletedPayload struct {
	NotifierUID string `json:"notifier_uid"`
}

type ReSystemEventV1CertsUpdatedPayload struct {
}

type ReSystemEventV1 struct {
	Type    ReSystemEventV1Type `json:"type"`
	Payload interface{}         `json:"payload"`
}

func (r *ReSystemEventV1) UnmarshalJSON(b []byte) error {
	objMap := map[string]*json.RawMessage{}

	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	objType, ok := objMap["type"]
	if !ok {
		return fmt.Errorf("missing type")
	}

	objPayload, ok := objMap["payload"]
	if !ok {
		return fmt.Errorf("missing payload")
	}

	err = json.Unmarshal(*objType, &r.Type)
	if err != nil {
		return fmt.Errorf("can't decode type")
	}

	var payload interface{}

	switch r.Type {
	case ReSystemEventV1TypeCheckAdded:
		payload = &ReSystemEventV1CheckAddedPayload{}
	case ReSystemEventV1TypeCheckUpdated:
		payload = &ReSystemEventV1CheckUpdatedPayload{}
	case ReSystemEventV1TypeCheckDeleted:
		payload = &ReSystemEventV1CheckDeletedPayload{}
	case ReSystemEventV1TypeHostAdded:
		payload = &ReSystemEventV1HostAddedPayload{}
	case ReSystemEventV1TypeHostUpdated:
		payload = &ReSystemEventV1HostUpdatedPayload{}
	case ReSystemEventV1TypeHostDeleted:
		payload = &ReSystemEventV1HostDeletedPayload{}
	case ReSystemEventV1TypeTagAdded:
		payload = &ReSystemEventV1TagAddedPayload{}
	case ReSystemEventV1TypeTagUpdated:
		payload = &ReSystemEventV1TagUpdatedPayload{}
	case ReSystemEventV1TypeTagDeleted:
		payload = &ReSystemEventV1TagDeletedPayload{}
	case ReSystemEventV1TypeNotifierAdded:
		payload = &ReSystemEventV1NotifierAddedPayload{}
	case ReSystemEventV1TypeNotifierUpdated:
		payload = &ReSystemEventV1NotifierUpdatedPayload{}
	case ReSystemEventV1TypeNotifierDeleted:
		payload = &ReSystemEventV1NotifierDeletedPayload{}
	case ReSystemEventV1TypeCertsUpdated:
		payload = &ReSystemEventV1CertsUpdatedPayload{}
	default:
		return fmt.Errorf("unknown type '%s'", r.Type)
	}

	err = json.Unmarshal(*objPayload, payload)
	if err != nil {
		return err
	}

	r.Payload = payload

	return nil
}
