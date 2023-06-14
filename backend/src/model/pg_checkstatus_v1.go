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

import "time"

type PgCheckStatusV1Status string

const (
	PgCheckStatusV1StatusOK   PgCheckStatusV1Status = "ok"
	PgCheckStatusV1StatusWarn PgCheckStatusV1Status = "warn"
	PgCheckStatusV1StatusCrit PgCheckStatusV1Status = "crit"
	PgCheckStatusV1StatusUnkn PgCheckStatusV1Status = "unkn"
)

type PgCheckStatusV1 struct {
	UID             string
	CheckUID        string
	Status          PgCheckStatusV1Status
	Message         string
	Data            map[string]interface{}
	DatetimeCreated time.Time
}
