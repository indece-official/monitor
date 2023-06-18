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

func (s *Service) getUnnotifiedStatusChangeKey(checkUID string, notifierUID string) string {
	return fmt.Sprintf("%s:%s", checkUID, notifierUID)
}

func (s *Service) SetUnnotifiedStatusChange(reUnnotifiedStatusChange *model.ReUnnotifiedStatusChangeV1) error {
	s.mutexUnnotifiedStatusChanges.Lock()
	defer s.mutexUnnotifiedStatusChanges.Unlock()

	key := s.getUnnotifiedStatusChangeKey(
		reUnnotifiedStatusChange.CheckUID,
		reUnnotifiedStatusChange.NotifierUID,
	)

	s.unnotifiedStatusChanges[key] = reUnnotifiedStatusChange

	return nil
}

func (s *Service) GetUnnotifiedStatusChanges() ([]*model.ReUnnotifiedStatusChangeV1, error) {
	s.mutexUnnotifiedStatusChanges.Lock()
	defer s.mutexUnnotifiedStatusChanges.Unlock()

	reUnnotifiedStatusChanges := []*model.ReUnnotifiedStatusChangeV1{}

	for _, reUnnotifiedStatusChange := range s.unnotifiedStatusChanges {
		reUnnotifiedStatusChanges = append(reUnnotifiedStatusChanges, reUnnotifiedStatusChange)
	}

	return reUnnotifiedStatusChanges, nil
}

func (s *Service) DeleteUnnotifiedStatusChange(checkUID string, notifierUID string) error {
	s.mutexUnnotifiedStatusChanges.Lock()
	defer s.mutexUnnotifiedStatusChanges.Unlock()

	key := s.getUnnotifiedStatusChangeKey(
		checkUID,
		notifierUID,
	)

	_, ok := s.unnotifiedStatusChanges[key]
	if !ok {
		return nil
	}

	delete(s.unnotifiedStatusChanges, key)

	return nil
}
