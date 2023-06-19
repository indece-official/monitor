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
	SubscribeForAgentActions(agentUID string) (chan *model.ReAgentActionV1, *Subscription, error)
	PublishAgentAction(reAgentAction *model.ReAgentActionV1) error

	AddOpenAgentAction(reAgentAction *model.ReAgentActionV1) error
	GetOpenAgentAction(agentUID string, actionUID string) (*model.ReAgentActionV1, error)
	GetOpenAgentActions(agentUID string) ([]*model.ReAgentActionV1, error)
	GetAllOpenAgentActions() ([]*model.ReAgentActionV1, error)
	DeleteOpenAgentAction(agentUID string, actionUID string) error

	SubscribeForAgentEvents() (chan *model.ReAgentEventV1, *Subscription, error)
	PublishAgentEvent(reAgentEvent *model.ReAgentEventV1) error

	SetAgentStatus(agentUID string, reStatus *model.ReAgentStatusV1) error
	GetAgentStatus(agentUID string) (*model.ReAgentStatusV1, error)
	DeleteAgentStatus(agentUID string) error

	UpsertHostCheckStatus(hostUID string, reStatusCheck *model.ReHostStatusV1Check) error
	GetHostStatus(hostUID string) (*model.ReHostStatusV1, error)
	GetHostCheckStatus(hostUID string, checkUID string) (*model.ReHostStatusV1Check, error)
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

	SetNotification(reNotification *model.ReNotificationV1) error
	GetNotifications() ([]*model.ReNotificationV1, error)
	DeleteNotification(checkUID string, notifierUID string) error
}

type Service struct {
	log                  *logger.Log
	agentActions         map[string]map[string]chan *model.ReAgentActionV1
	openAgentActions     map[string]map[string]*model.ReAgentActionV1
	agentEvents          map[string]chan *model.ReAgentEventV1
	agentStatus          map[string]*model.ReAgentStatusV1
	hostStatus           map[string]*model.ReHostStatusV1
	systemEvents         map[string]chan *model.ReSystemEventV1
	userSessions         map[string]*model.ReUserSessionV1
	notifications        map[string]*model.ReNotificationV1
	mutexOpenAgentAction sync.Mutex
	mutexAgentStatus     sync.Mutex
	mutexHostStatus      sync.Mutex
	mutexNotifications   sync.Mutex
	setupToken           *SetupToken
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
		log:              logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
		agentActions:     map[string]map[string]chan *model.ReAgentActionV1{},
		openAgentActions: map[string]map[string]*model.ReAgentActionV1{},
		agentEvents:      map[string]chan *model.ReAgentEventV1{},
		agentStatus:      map[string]*model.ReAgentStatusV1{},
		hostStatus:       map[string]*model.ReHostStatusV1{},
		systemEvents:     map[string]chan *model.ReSystemEventV1{},
		userSessions:     map[string]*model.ReUserSessionV1{},
		notifications:    map[string]*model.ReNotificationV1{},
	}
}

var _ gousu.ServiceFactory = NewService
