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

package cert

import (
	"fmt"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
)

const ServiceName = "cert"

type IService interface {
	GenerateCACert() (*PEMCert, error)
	GenerateServerCert(hostname string, ca *PEMCert) (*PEMCert, error)
	GenerateClientCert(agentUID string, clientsPEM *PEMCert) (*PEMCert, error)
}

type Service struct {
	log *logger.Log
}

var _ IService = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Health() error {
	return nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ gousu.ServiceFactory = NewService
