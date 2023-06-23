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

import "gopkg.in/guregu/null.v4"

type PgCheckerV1ParamType string

const (
	PgCheckerV1ParamTypeText     PgCheckerV1ParamType = "text"
	PgCheckerV1ParamTypePassword PgCheckerV1ParamType = "password"
	PgCheckerV1ParamTypeNumber   PgCheckerV1ParamType = "number"
	PgCheckerV1ParamTypeSelect   PgCheckerV1ParamType = "select"
	PgCheckerV1ParamTypeDuration PgCheckerV1ParamType = "duration"
	PgCheckerV1ParamTypeBoolean  PgCheckerV1ParamType = "boolean"
)

type PgCheckerV1Param struct {
	Name     string               `json:"name"`
	Label    string               `json:"label"`
	Hint     null.String          `json:"hint"`
	Required bool                 `json:"required"`
	Type     PgCheckerV1ParamType `json:"type"`
	Options  []string             `json:"options"`
}

type PgCheckerV1Value struct {
	Name    string      `json:"name"`
	Label   string      `json:"label"`
	Type    ValueType   `json:"type"`
	MinWarn null.String `json:"min_warn"`
	MinCrit null.String `json:"min_crit"`
	MaxWarn null.String `json:"max_warn"`
	MaxCrit null.String `json:"max_crit"`
}

type PgCheckerV1Capabilities struct {
	Params          []*PgCheckerV1Param `json:"params"`
	Values          []*PgCheckerV1Value `json:"values"`
	DefaultSchedule null.String         `json:"default_schedule"`
	DefaultTimeout  null.String         `json:"default_timeout"`
}

type PgCheckerV1 struct {
	UID          string
	Type         string
	AgentUID     string
	Version      string
	Name         string
	Capabilities *PgCheckerV1Capabilities
	CustomChecks bool
}
