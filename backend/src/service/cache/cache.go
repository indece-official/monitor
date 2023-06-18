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

package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/model"
	"gopkg.in/guregu/null.v4"
)

const ServiceName = "cache"

type IService interface {
	SubscribeForConnectorActions(connectorUID string) (chan *model.ReConnectorActionV1, *Subscription, error)
	PublishConnectorAction(reConnectorAction *model.ReConnectorActionV1) error

	AddOpenConnectorAction(reConnectorAction *model.ReConnectorActionV1) error
	GetOpenConnectorAction(connectorUID string, actionUID string) (*model.ReConnectorActionV1, error)
	GetOpenConnectorActions(connectorUID string) ([]*model.ReConnectorActionV1, error)
	GetAllOpenConnectorActions() ([]*model.ReConnectorActionV1, error)
	DeleteOpenConnectorAction(connectorUID string, actionUID string) error

	SubscribeForConnectorEvents() (chan *model.ReConnectorEventV1, *Subscription, error)
	PublishConnectorEvent(reConnectorEvent *model.ReConnectorEventV1) error

	SetConnectorStatus(connectorUID string, reStatus *model.ReConnectorStatusV1) error
	GetConnectorStatus(connectorUID string) (*model.ReConnectorStatusV1, error)
	DeleteConnectorStatus(connectorUID string) error

	UpsertHostCheckStatus(hostUID string, reStatusCheck *model.ReHostStatusV1Check) error
	GetHostStatus(hostUID string) (*model.ReHostStatusV1, error)
	DeleteHostStatus(hostUID string) error
	DeleteAllHostStatuses() error

	SubscribeForSystemEvents() (chan *model.ReSystemEventV1, *Subscription, error)
	PublishSystemEvent(reSystemEvent *model.ReSystemEventV1) error

	SetUserSession(sessionKey string, sessionData *model.ReUserSessionV1) error
	GetUserSession(sessionKey string) (*model.ReUserSessionV1, error)
	DeleteUserSession(sessionKey string) error

	SetSetupToken(setupToken string, duration time.Duration) error
	GetSetupToken() (null.String, error)
	DeleteSetupToken() error
}

type Service struct {
	log                      *logger.Log
	connectorActions         map[string]map[string]chan *model.ReConnectorActionV1
	openConnectorActions     map[string]map[string]*model.ReConnectorActionV1
	connectorEvents          map[string]chan *model.ReConnectorEventV1
	connectorStatus          map[string]*model.ReConnectorStatusV1
	hostStatus               map[string]*model.ReHostStatusV1
	systemEvents             map[string]chan *model.ReSystemEventV1
	userSessions             map[string]*model.ReUserSessionV1
	mutexOpenConnectorAction sync.Mutex
	mutexConnectorStatus     sync.Mutex
	mutexHostStatus          sync.Mutex
	setupToken               *SetupToken
}

var _ IService = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Health() error {
	return nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log:                  logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
		connectorActions:     map[string]map[string]chan *model.ReConnectorActionV1{},
		openConnectorActions: map[string]map[string]*model.ReConnectorActionV1{},
		connectorEvents:      map[string]chan *model.ReConnectorEventV1{},
		connectorStatus:      map[string]*model.ReConnectorStatusV1{},
		hostStatus:           map[string]*model.ReHostStatusV1{},
		systemEvents:         map[string]chan *model.ReSystemEventV1{},
		userSessions:         map[string]*model.ReUserSessionV1{},
	}
}

var _ gousu.ServiceFactory = NewService
