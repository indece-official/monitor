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

func (s *Service) SetAgentStatus(agentUID string, reStatus *model.ReAgentStatusV1) error {
	s.mutexAgentStatus.Lock()
	defer s.mutexAgentStatus.Unlock()

	s.agentStatus[agentUID] = reStatus

	return nil
}

func (s *Service) GetAgentStatus(agentUID string) (*model.ReAgentStatusV1, error) {
	s.mutexAgentStatus.Lock()
	defer s.mutexAgentStatus.Unlock()

	reStatus, ok := s.agentStatus[agentUID]
	if !ok {
		return nil, nil
	}

	return reStatus, nil
}

func (s *Service) DeleteAgentStatus(agentUID string) error {
	s.mutexAgentStatus.Lock()
	defer s.mutexAgentStatus.Unlock()

	_, ok := s.agentStatus[agentUID]
	if !ok {
		return nil
	}

	delete(s.agentStatus, agentUID)

	return nil
}
