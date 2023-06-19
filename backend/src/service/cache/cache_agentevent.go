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

func (s *Service) SubscribeForAgentEvents() (chan *model.ReAgentEventV1, *Subscription, error) {
	subscriptionUID, err := utils.UUID()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating subscription uid: %s", err)
	}

	agentEvents := make(chan *model.ReAgentEventV1, 10)

	subscription := &Subscription{}
	subscription.uid = subscriptionUID
	subscription.onUnsubscribeFunc = func(uid string) {
		if _, ok := s.agentEvents[uid]; !ok {
			return
		}

		close(s.agentEvents[uid])
		delete(s.agentEvents, uid)
	}

	s.agentEvents[subscriptionUID] = agentEvents

	return agentEvents, subscription, nil
}

func (s *Service) PublishAgentEvent(reAgentEvent *model.ReAgentEventV1) error {
	for _, agentEvents := range s.agentEvents {
		agentEvents <- reAgentEvent
	}

	return nil
}
