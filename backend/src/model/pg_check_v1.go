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

type PgCheckV1Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PgCheckV1Config struct {
	Params  []*PgCheckV1Param `json:"params"`
	Timeout null.String       `json:"timeout"`
}

type PgCheckV1 struct {
	UID              string
	CheckerUID       string
	Name             string
	Type             null.String
	Schedule         null.String
	Config           *PgCheckV1Config
	Custom           bool
	DatetimeDisabled null.Time

	Statuses []*PgCheckStatusV1
}
