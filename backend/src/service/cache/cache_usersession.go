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

func (s *Service) SetUserSession(sessionKey string, sessionData *model.ReUserSessionV1) error {
	s.userSessions[sessionKey] = sessionData

	return nil
}

func (s *Service) GetUserSession(sessionKey string) (*model.ReUserSessionV1, error) {
	sessionData, ok := s.userSessions[sessionKey]
	if !ok {
		return nil, nil
	}

	return sessionData, nil
}

func (s *Service) DeleteUserSession(sessionKey string) error {
	_, ok := s.userSessions[sessionKey]
	if !ok {
		return nil
	}

	delete(s.userSessions, sessionKey)

	return nil
}
