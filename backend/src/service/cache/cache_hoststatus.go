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

func (s *Service) UpsertHostCheckStatus(hostUID string, reStatusCheck *model.ReHostStatusV1Check) error {
	s.mutexHostStatus.Lock()
	defer s.mutexHostStatus.Unlock()

	hostStatus := s.hostStatus[hostUID]
	if hostStatus == nil {
		s.hostStatus[hostUID] = &model.ReHostStatusV1{}
		s.hostStatus[hostUID].HostUID = hostUID
		s.hostStatus[hostUID].Checks = []*model.ReHostStatusV1Check{}

		hostStatus = s.hostStatus[hostUID]
	}

	found := false

	for i, checkStatus := range hostStatus.Checks {
		if checkStatus.CheckUID == reStatusCheck.CheckUID {
			hostStatus.Checks[i] = reStatusCheck
			found = true
			break
		}
	}

	if !found {
		hostStatus.Checks = append(hostStatus.Checks, reStatusCheck)
	}

	return nil
}

func (s *Service) GetHostStatus(hostUID string) (*model.ReHostStatusV1, error) {
	s.mutexHostStatus.Lock()
	defer s.mutexHostStatus.Unlock()

	reStatus, ok := s.hostStatus[hostUID]
	if !ok {
		return nil, nil
	}

	return reStatus, nil
}

func (s *Service) DeleteHostStatus(hostUID string) error {
	s.mutexHostStatus.Lock()
	defer s.mutexHostStatus.Unlock()

	_, ok := s.hostStatus[hostUID]
	if !ok {
		return nil
	}

	delete(s.hostStatus, hostUID)

	return nil
}

func (s *Service) DeleteAllHostStatuses() error {
	s.mutexHostStatus.Lock()
	defer s.mutexHostStatus.Unlock()

	s.hostStatus = map[string]*model.ReHostStatusV1{}

	return nil
}
