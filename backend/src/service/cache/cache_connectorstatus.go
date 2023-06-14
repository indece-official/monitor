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

func (s *Service) SetConnectorStatus(connectorUID string, reStatus *model.ReConnectorStatusV1) error {
	s.mutexConnectorStatus.Lock()
	defer s.mutexConnectorStatus.Unlock()

	s.connectorStatus[connectorUID] = reStatus

	return nil
}

func (s *Service) GetConnectorStatus(connectorUID string) (*model.ReConnectorStatusV1, error) {
	s.mutexConnectorStatus.Lock()
	defer s.mutexConnectorStatus.Unlock()

	reStatus, ok := s.connectorStatus[connectorUID]
	if !ok {
		return nil, nil
	}

	return reStatus, nil
}

func (s *Service) DeleteConnectorStatus(connectorUID string) error {
	s.mutexConnectorStatus.Lock()
	defer s.mutexConnectorStatus.Unlock()

	_, ok := s.connectorStatus[connectorUID]
	if !ok {
		return nil
	}

	delete(s.connectorStatus, connectorUID)

	return nil
}
