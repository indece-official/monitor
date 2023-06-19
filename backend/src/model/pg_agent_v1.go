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
	"time"

	"gopkg.in/guregu/null.v4"
)

type PgAgentV1Cert struct {
	TLSClientCrt string    `json:"tls_client_crt"`
	TLSClientKey string    `json:"tls_client_key"`
	CreateAt     time.Time `json:"created_at"`
	ValidUntil   time.Time `json:"valid_until"`
}

type PgAgentV1Certs struct {
	Certs []*PgAgentV1Cert `json:"certs"`
}

type PgAgentV1 struct {
	UID                string
	HostUID            string
	Type               null.String
	Version            null.String
	Certs              *PgAgentV1Certs
	DatetimeRegistered null.Time
}
