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

package postgres

import (
	"context"
	"fmt"

	"github.com/indece-official/go-gousu/gousupostgres/v2"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"gopkg.in/guregu/null.v4"

	"github.com/indece-official/monitor/backend/src/assets"
	"github.com/indece-official/monitor/backend/src/model"
)

// ServiceName defines the name of the postgres service
const ServiceName = gousupostgres.ServiceName

type IService interface {
	gousu.IService

	UpsertConfigProperty(qctx context.Context, pgConfigProperty *model.PgConfigPropertyV1) error
	GetConfigProperties(qctx context.Context, filter *GetConfigPropertiesFilter) (map[model.PgConfigPropertyV1Key]*model.PgConfigPropertyV1, error)

	AddUser(qctx context.Context, pgUser *model.PgUserV1) error
	UpdateUser(qctx context.Context, userUID string, pgUser *model.PgUserV1) error
	UpdateUserPassword(qctx context.Context, userUID string, passwordHash null.String) error
	DeleteUser(qctx context.Context, userUID string) error
	GetUsers(qctx context.Context, filter *GetUsersFilter) ([]*model.PgUserV1, error)

	AddHost(qctx context.Context, pgHost *model.PgHostV1) error
	UpdateHost(qctx context.Context, hostUID string, pgHost *model.PgHostV1) error
	DeleteHost(qctx context.Context, hostUID string) error
	GetHosts(qctx context.Context, filter *GetHostsFilter) ([]*model.PgHostV1, error)

	AddConnector(qctx context.Context, pgConnector *model.PgConnectorV1) error
	UpdateConnector(qctx context.Context, connectorUID string, pgConnector *model.PgConnectorV1) error
	DeleteConnector(qctx context.Context, connectorUID string) error
	GetConnectors(qctx context.Context, filter *GetConnectorsFilter) ([]*model.PgConnectorV1, error)

	AddChecker(qctx context.Context, pgChecker *model.PgCheckerV1) error
	GetCheckers(qctx context.Context, filter *GetCheckersFilter) ([]*model.PgCheckerV1, error)

	AddCheck(qctx context.Context, pgCheck *model.PgCheckV1) error
	UpdateCheck(qctx context.Context, checkUID string, pgCheck *model.PgCheckV1) error
	DeleteCheck(qctx context.Context, checkUID string) error
	GetChecks(qctx context.Context, filter *GetChecksFilter) ([]*model.PgCheckV1, error)

	AddCheckStatus(qctx context.Context, pgCheckStatus *model.PgCheckStatusV1) error

	AddTag(qctx context.Context, pgTag *model.PgTagV1) error
	UpdateTag(qctx context.Context, tagUID string, pgTag *model.PgTagV1) error
	DeleteTag(qctx context.Context, tagUID string) error
	GetTags(qctx context.Context, filter *GetTagsFilter) ([]*model.PgTagV1, error)
}

// Service provides the interaction with the postgresql database
type Service struct {
	log             *logger.Log
	postgresService gousupostgres.IService
}

var _ IService = (*Service)(nil)

// Name returns the name of the postgres service defined by ServiceName
func (s *Service) Name() string {
	return ServiceName
}

// Start initializes the connection to the postgres database and executed both setup.sql and update.sql
// after connecting
func (s *Service) Start() error {
	return s.postgresService.Start()
}

// Stop closes the connection to the postgres database
func (s *Service) Stop() error {
	return s.postgresService.Stop()
}

// Health checks the health of the postgres-service by pinging the postgres database
func (s *Service) Health() error {
	return s.postgresService.Health()
}

// NewService creates a new instance of postgres-service, should be used instead
//
//	of generating it manually
func NewService(ctx gousu.IContext) gousu.IService {
	setupSQL, err := assets.ReadFile("sql/setup.sql")
	gousu.CheckError(err)

	updateSQL, err := assets.ReadFile("sql/update.sql")
	gousu.CheckError(err)

	options := &gousupostgres.Options{
		SetupSQL:  setupSQL,
		UpdateSQL: updateSQL,
	}

	return &Service{
		log:             logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
		postgresService: gousupostgres.NewServiceBase(ctx, options),
	}
}

var _ gousu.ServiceFactory = NewService
