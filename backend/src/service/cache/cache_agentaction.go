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

func (s *Service) SubscribeForAgentActions(agentUID string) (chan *model.ReAgentActionV1, *Subscription, error) {
	if _, ok := s.agentActions[agentUID]; !ok {
		s.agentActions[agentUID] = map[string]chan *model.ReAgentActionV1{}
	}

	subscriptionUID, err := utils.UUID()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating subscription uid: %s", err)
	}

	agentActions := make(chan *model.ReAgentActionV1, 10)

	subscription := &Subscription{}
	subscription.uid = subscriptionUID
	subscription.onUnsubscribeFunc = func(uid string) {
		if _, ok := s.agentActions[agentUID]; !ok {
			return
		}

		if _, ok := s.agentActions[agentUID][uid]; !ok {
			return
		}

		close(s.agentActions[agentUID][uid])
		delete(s.agentActions[agentUID], uid)
	}

	s.agentActions[agentUID][subscriptionUID] = agentActions

	return agentActions, subscription, nil
}

func (s *Service) PublishAgentAction(reAgentAction *model.ReAgentActionV1) error {
	if _, ok := s.agentActions[reAgentAction.AgentUID]; !ok {
		return nil
	}

	for _, agentActions := range s.agentActions[reAgentAction.AgentUID] {
		agentActions <- reAgentAction
	}

	return nil
}

func (s *Service) AddOpenAgentAction(reAgentAction *model.ReAgentActionV1) error {
	s.mutexOpenAgentAction.Lock()
	defer s.mutexOpenAgentAction.Unlock()

	if _, ok := s.openAgentActions[reAgentAction.AgentUID]; !ok {
		s.openAgentActions[reAgentAction.AgentUID] = map[string]*model.ReAgentActionV1{}
	}

	s.openAgentActions[reAgentAction.AgentUID][reAgentAction.ActionUID] = reAgentAction

	return nil
}

func (s *Service) GetOpenAgentAction(agentUID string, actionUID string) (*model.ReAgentActionV1, error) {
	s.mutexOpenAgentAction.Lock()
	defer s.mutexOpenAgentAction.Unlock()

	if _, ok := s.openAgentActions[agentUID]; !ok {
		return nil, nil
	}

	if _, ok := s.openAgentActions[agentUID][actionUID]; !ok {
		return nil, nil
	}

	return s.openAgentActions[agentUID][actionUID], nil
}

func (s *Service) GetOpenAgentActions(agentUID string) ([]*model.ReAgentActionV1, error) {
	s.mutexOpenAgentAction.Lock()
	defer s.mutexOpenAgentAction.Unlock()

	actions := []*model.ReAgentActionV1{}

	if _, ok := s.openAgentActions[agentUID]; !ok {
		return actions, nil
	}

	for _, action := range s.openAgentActions[agentUID] {
		actions = append(actions, action)
	}

	return actions, nil
}

func (s *Service) GetAllOpenAgentActions() ([]*model.ReAgentActionV1, error) {
	s.mutexOpenAgentAction.Lock()
	defer s.mutexOpenAgentAction.Unlock()

	actions := []*model.ReAgentActionV1{}

	for _, actionMap := range s.openAgentActions {
		for _, action := range actionMap {
			actions = append(actions, action)
		}
	}

	return actions, nil
}

func (s *Service) DeleteOpenAgentAction(agentUID string, actionUID string) error {
	s.mutexOpenAgentAction.Lock()
	defer s.mutexOpenAgentAction.Unlock()

	if _, ok := s.openAgentActions[agentUID]; !ok {
		return nil
	}

	if _, ok := s.openAgentActions[agentUID][actionUID]; !ok {
		return nil
	}

	delete(s.openAgentActions[agentUID], actionUID)

	return nil
}
