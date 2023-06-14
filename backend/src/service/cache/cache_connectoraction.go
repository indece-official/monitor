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

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/utils"
)

func (s *Service) SubscribeForConnectorActions(connectorUID string) (chan *model.ReConnectorActionV1, *Subscription, error) {
	if _, ok := s.connectorActions[connectorUID]; !ok {
		s.connectorActions[connectorUID] = map[string]chan *model.ReConnectorActionV1{}
	}

	subscriptionUID, err := utils.UUID()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating subscription uid: %s", err)
	}

	connectorActions := make(chan *model.ReConnectorActionV1, 10)

	subscription := &Subscription{}
	subscription.uid = subscriptionUID
	subscription.onUnsubscribeFunc = func(uid string) {
		if _, ok := s.connectorActions[connectorUID]; !ok {
			return
		}

		if _, ok := s.connectorActions[connectorUID][uid]; !ok {
			return
		}

		close(s.connectorActions[connectorUID][uid])
		delete(s.connectorActions[connectorUID], uid)
	}

	s.connectorActions[connectorUID][subscriptionUID] = connectorActions

	return connectorActions, subscription, nil
}

func (s *Service) PublishConnectorAction(reConnectorAction *model.ReConnectorActionV1) error {
	if _, ok := s.connectorActions[reConnectorAction.ConnectorUID]; !ok {
		return nil
	}

	for _, connectorActions := range s.connectorActions[reConnectorAction.ConnectorUID] {
		connectorActions <- reConnectorAction
	}

	return nil
}

func (s *Service) AddOpenConnectorAction(reConnectorAction *model.ReConnectorActionV1) error {
	s.mutexOpenConnectorAction.Lock()
	defer s.mutexOpenConnectorAction.Unlock()

	if _, ok := s.openConnectorActions[reConnectorAction.ConnectorUID]; !ok {
		s.openConnectorActions[reConnectorAction.ConnectorUID] = map[string]*model.ReConnectorActionV1{}
	}

	s.openConnectorActions[reConnectorAction.ConnectorUID][reConnectorAction.ActionUID] = reConnectorAction

	return nil
}

func (s *Service) GetOpenConnectorAction(connectorUID string, actionUID string) (*model.ReConnectorActionV1, error) {
	s.mutexOpenConnectorAction.Lock()
	defer s.mutexOpenConnectorAction.Unlock()

	if _, ok := s.openConnectorActions[connectorUID]; !ok {
		return nil, nil
	}

	if _, ok := s.openConnectorActions[connectorUID][actionUID]; !ok {
		return nil, nil
	}

	return s.openConnectorActions[connectorUID][actionUID], nil
}

func (s *Service) GetOpenConnectorActions(connectorUID string) ([]*model.ReConnectorActionV1, error) {
	s.mutexOpenConnectorAction.Lock()
	defer s.mutexOpenConnectorAction.Unlock()

	actions := []*model.ReConnectorActionV1{}

	if _, ok := s.openConnectorActions[connectorUID]; !ok {
		return actions, nil
	}

	for _, action := range s.openConnectorActions[connectorUID] {
		actions = append(actions, action)
	}

	return actions, nil
}

func (s *Service) GetAllOpenConnectorActions() ([]*model.ReConnectorActionV1, error) {
	s.mutexOpenConnectorAction.Lock()
	defer s.mutexOpenConnectorAction.Unlock()

	actions := []*model.ReConnectorActionV1{}

	for _, actionMap := range s.openConnectorActions {
		for _, action := range actionMap {
			actions = append(actions, action)
		}
	}

	return actions, nil
}

func (s *Service) DeleteOpenConnectorAction(connectorUID string, actionUID string) error {
	s.mutexOpenConnectorAction.Lock()
	defer s.mutexOpenConnectorAction.Unlock()

	if _, ok := s.openConnectorActions[connectorUID]; !ok {
		return nil
	}

	if _, ok := s.openConnectorActions[connectorUID][actionUID]; !ok {
		return nil
	}

	delete(s.openConnectorActions[connectorUID], actionUID)

	return nil
}
