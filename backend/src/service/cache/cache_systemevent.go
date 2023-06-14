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

func (s *Service) SubscribeForSystemEvents() (chan *model.ReSystemEventV1, *Subscription, error) {
	subscriptionUID, err := utils.UUID()
	if err != nil {
		return nil, nil, fmt.Errorf("error generating subscription uid: %s", err)
	}

	systemEvents := make(chan *model.ReSystemEventV1, 10)

	subscription := &Subscription{}
	subscription.uid = subscriptionUID
	subscription.onUnsubscribeFunc = func(uid string) {
		if _, ok := s.systemEvents[uid]; !ok {
			return
		}

		close(s.systemEvents[uid])
		delete(s.systemEvents, uid)
	}

	s.systemEvents[subscriptionUID] = systemEvents

	return systemEvents, subscription, nil
}

func (s *Service) PublishSystemEvent(reSystemEvent *model.ReSystemEventV1) error {
	for _, systemEvents := range s.systemEvents {
		systemEvents <- reSystemEvent
	}

	return nil
}
