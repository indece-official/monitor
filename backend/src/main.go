//go:generate go run assets/generate.go
//go:generate sh -c "mkdir -p generated/model/apipublic && oapi-codegen --package=apipublic --generate=types ../assets/swagger/apipublic.yml > ./generated/model/apipublic/apipublic.gen.go"
//go:generate /bin/sh -c "mkdir -p generated/model/apiconnector && protoc --go_out=./generated/model/ --go-grpc_out=./generated/model/ --proto_path=../assets/grpc/ ../assets/grpc/connector.proto"

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

package main

import (
	"fmt"

	"github.com/indece-official/monitor/backend/src/buildvars"
	"github.com/indece-official/monitor/backend/src/controller/connector"
	"github.com/indece-official/monitor/backend/src/controller/cron"
	"github.com/indece-official/monitor/backend/src/controller/initializer"
	"github.com/indece-official/monitor/backend/src/controller/public"
	"github.com/indece-official/monitor/backend/src/service/cache"
	"github.com/indece-official/monitor/backend/src/service/cert"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"github.com/indece-official/monitor/backend/src/service/smtp"
	"github.com/indece-official/monitor/backend/src/service/template"

	"github.com/indece-official/go-gousu/v2/gousu"
)

func main() {
	runner := gousu.NewRunner(buildvars.ProjectName, fmt.Sprintf("%s (Build %s)", buildvars.BuildVersion, buildvars.BuildDate))

	runner.CreateService(postgres.NewService)
	runner.CreateService(smtp.NewService)
	runner.CreateService(cache.NewService)
	runner.CreateService(template.NewService)
	runner.CreateService(cert.NewService)
	runner.CreateController(initializer.NewController)
	runner.CreateController(public.NewController)
	runner.CreateController(connector.NewController)
	runner.CreateController(cron.NewController)
	runner.CreateController(gousu.NewActuatorController)

	runner.Run()
}
