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
	"gopkg.in/guregu/null.v4"
)

type PgUserV1Source string

const (
	PgUserV1SourceLocal PgUserV1Source = "local"
)

type PgUserV1Role string

const (
	PgUserV1RoleShow  PgUserV1Role = "show"
	PgUserV1RoleAdmin PgUserV1Role = "admin"
	PgUserV1RoleSetup PgUserV1Role = "setup"
)

type PgUserV1 struct {
	UID            string
	Source         PgUserV1Source
	Username       string
	Name           null.String
	Email          null.String
	LocalRoles     []PgUserV1Role
	PasswordHash   null.String
	DatetimeLocked null.Time
}
