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

type PgConfigPropertyV1Key string

const (
	PgConfigPropertyV1KeySetupFinished PgConfigPropertyV1Key = "setup_finished"
	PgConfigPropertyV1KeyCompanyName   PgConfigPropertyV1Key = "company_name"
	PgConfigPropertyV1KeyTLSCaCrt      PgConfigPropertyV1Key = "tls_ca_crt"
	PgConfigPropertyV1KeyTLSCaKey      PgConfigPropertyV1Key = "tls_ca_key"
	PgConfigPropertyV1KeyTLSServerCrt  PgConfigPropertyV1Key = "tls_server_crt"
	PgConfigPropertyV1KeyTLSServerKey  PgConfigPropertyV1Key = "tls_server_key"
	PgConfigPropertyV1KeyAgentHost     PgConfigPropertyV1Key = "agent_host"
	PgConfigPropertyV1KeyAgentPort     PgConfigPropertyV1Key = "agent_port"
)

const (
	PgConfigPropertyV1False = "false"
	PgConfigPropertyV1True  = "true"
)

type PgConfigPropertyV1Protection string

const (
	PgConfigPropertyV1ProtectionPrivate   PgConfigPropertyV1Protection = "private"
	PgConfigPropertyV1ProtectionProtected PgConfigPropertyV1Protection = "protected"
	PgConfigPropertyV1ProtectionPublic    PgConfigPropertyV1Protection = "public"
)

var PgConfigPropertyV1Protections = map[PgConfigPropertyV1Key]PgConfigPropertyV1Protection{
	PgConfigPropertyV1KeySetupFinished: PgConfigPropertyV1ProtectionProtected,
	PgConfigPropertyV1KeyTLSCaCrt:      PgConfigPropertyV1ProtectionProtected,
	PgConfigPropertyV1KeyTLSServerCrt:  PgConfigPropertyV1ProtectionProtected,
	PgConfigPropertyV1KeyAgentHost:     PgConfigPropertyV1ProtectionPublic,
	PgConfigPropertyV1KeyAgentPort:     PgConfigPropertyV1ProtectionPublic,
}

type PgConfigPropertyV1 struct {
	Key   PgConfigPropertyV1Key
	Value string
}
