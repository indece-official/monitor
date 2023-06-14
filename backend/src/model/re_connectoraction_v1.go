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
	"time"
)

type ReConnectorActionV1Type string

const (
	ReConnectorActionV1TypeCheck ReConnectorActionV1Type = "check"
)

type ReConnectorActionV1CheckPayloadParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ReConnectorActionV1CheckPayload struct {
	CheckUID        string                                  `json:"check_uid"`
	CheckerType     string                                  `json:"checker_type"`
	Params          []*ReConnectorActionV1CheckPayloadParam `json:"params"`
	TimeoutAt       time.Time                               `json:"timeout_at"`
	TimeoutDuration time.Duration                           `json:"timeout_duration"`
}

type ReConnectorActionV1 struct {
	ActionUID    string                  `json:"action_uid"`
	Type         ReConnectorActionV1Type `json:"type"`
	ConnectorUID string                  `json:"connector_uid"`
	Payload      interface{}             `json:"payload"`
}

func (r *ReConnectorActionV1) UnmarshalJSON(b []byte) error {
	objMap := map[string]*json.RawMessage{}

	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	objType, ok := objMap["type"]
	if !ok {
		return fmt.Errorf("missing type")
	}

	objConnectorUID, ok := objMap["connector_uid"]
	if !ok {
		return fmt.Errorf("missing connector_uid")
	}

	objPayload, ok := objMap["payload"]
	if !ok {
		return fmt.Errorf("missing payload")
	}

	err = json.Unmarshal(*objType, &r.Type)
	if err != nil {
		return fmt.Errorf("can't decode type")
	}

	err = json.Unmarshal(*objConnectorUID, &r.ConnectorUID)
	if err != nil {
		return fmt.Errorf("can't decode connector_uid")
	}

	var payload interface{}

	switch r.Type {
	case ReConnectorActionV1TypeCheck:
		payload = &ReConnectorActionV1CheckPayload{}
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
