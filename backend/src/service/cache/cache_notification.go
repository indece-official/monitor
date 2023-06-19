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
)

func (s *Service) getNotificationKey(checkUID string, notifierUID string) string {
	return fmt.Sprintf("%s:%s", checkUID, notifierUID)
}

func (s *Service) SetNotification(reNotification *model.ReNotificationV1) error {
	s.mutexNotifications.Lock()
	defer s.mutexNotifications.Unlock()

	key := s.getNotificationKey(
		reNotification.CheckUID,
		reNotification.NotifierUID,
	)

	s.notifications[key] = reNotification

	return nil
}

func (s *Service) GetNotifications() ([]*model.ReNotificationV1, error) {
	s.mutexNotifications.Lock()
	defer s.mutexNotifications.Unlock()

	reNotifications := []*model.ReNotificationV1{}

	for _, reNotification := range s.notifications {
		reNotifications = append(reNotifications, reNotification)
	}

	return reNotifications, nil
}

func (s *Service) DeleteNotification(checkUID string, notifierUID string) error {
	s.mutexNotifications.Lock()
	defer s.mutexNotifications.Unlock()

	key := s.getNotificationKey(
		checkUID,
		notifierUID,
	)

	_, ok := s.notifications[key]
	if !ok {
		return nil
	}

	delete(s.notifications, key)

	return nil
}
