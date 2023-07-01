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

import "github.com/indece-official/monitor/backend/src/model"

func (s *Service) SetCheckStatus(checkUID string, reStatus *model.ReCheckStatusV1) error {
	s.mutexCheckStatus.Lock()
	defer s.mutexCheckStatus.Unlock()

	s.checkStatus[checkUID] = reStatus

	return nil
}

func (s *Service) GetCheckStatus(checkUID string) (*model.ReCheckStatusV1, error) {
	s.mutexCheckStatus.Lock()
	defer s.mutexCheckStatus.Unlock()

	reStatus, ok := s.checkStatus[checkUID]
	if !ok {
		return nil, nil
	}

	return reStatus, nil
}

func (s *Service) GetAllCheckStatuses() ([]*model.ReCheckStatusV1, error) {
	s.mutexCheckStatus.Lock()
	defer s.mutexCheckStatus.Unlock()

	reStatuses := []*model.ReCheckStatusV1{}

	for _, reStatus := range s.checkStatus {
		reStatuses = append(reStatuses, reStatus)
	}

	return reStatuses, nil
}

func (s *Service) DeleteCheckStatus(checkUID string) error {
	s.mutexCheckStatus.Lock()
	defer s.mutexCheckStatus.Unlock()

	_, ok := s.checkStatus[checkUID]
	if !ok {
		return nil
	}

	delete(s.checkStatus, checkUID)

	return nil
}

func (s *Service) DeleteAllCheckStatuses() error {
	s.mutexCheckStatus.Lock()
	defer s.mutexCheckStatus.Unlock()

	s.checkStatus = map[string]*model.ReCheckStatusV1{}

	return nil
}
