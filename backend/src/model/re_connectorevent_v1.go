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

	"gopkg.in/guregu/null.v4"
)

type ReConnectorEventV1Type string

const (
	ReConnectorEventV1TypeCheckResult ReConnectorEventV1Type = "check_result"
)

type ReConnectorEventV1CheckResultPayloadValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ReConnectorEventV1CheckResultPayload struct {
	ActionUID string                                       `json:"action_uid"`
	CheckUID  string                                       `json:"check_uid"`
	Message   string                                       `json:"message"`
	Values    []*ReConnectorEventV1CheckResultPayloadValue `json:"values"`
	Error     null.String                                  `json:"error"`
}

type ReConnectorEventV1 struct {
	Type         ReConnectorEventV1Type `json:"type"`
	ConnectorUID string                 `json:"connector_uid"`
	Payload      interface{}            `json:"payload"`
}

func (r *ReConnectorEventV1) UnmarshalJSON(b []byte) error {
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
	case ReConnectorEventV1TypeCheckResult:
		payload = &ReConnectorEventV1CheckResultPayload{}
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
