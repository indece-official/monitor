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
	"time"

	"gopkg.in/guregu/null.v4"
)

type SetupToken struct {
	Token      string
	ValitUntil time.Time
}

func (s *Service) SetSetupToken(setupToken string, duration time.Duration) error {
	s.setupToken = &SetupToken{
		Token:      setupToken,
		ValitUntil: time.Now().Add(duration),
	}

	return nil
}

func (s *Service) GetSetupToken() (null.String, error) {
	if s.setupToken == nil ||
		s.setupToken.ValitUntil.Before(time.Now()) {
		return null.String{}, nil
	}

	return null.StringFrom(s.setupToken.Token), nil
}

func (s *Service) DeleteSetupToken() error {
	s.setupToken = nil

	return nil
}
