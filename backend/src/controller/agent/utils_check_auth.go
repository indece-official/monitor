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

package agent

import (
	"context"
	"fmt"

	"github.com/indece-official/monitor/backend/src/service/postgres"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"gopkg.in/guregu/null.v4"
)

type Session struct {
	AgentUID string
}

func (c *Controller) checkAuth(ctx context.Context) (*Session, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		c.log.Warnf("No peer found")

		return nil, fmt.Errorf("unauthorized")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		c.log.Warnf("Unexpected peer transport credentials")

		return nil, fmt.Errorf("unauthorized")
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		c.log.Warnf("Could not verify peer certificate")

		return nil, fmt.Errorf("unauthorized")
	}

	pgAgents, err := c.postgresService.GetAgents(
		ctx,
		&postgres.GetAgentsFilter{
			AgentUID: null.StringFrom(tlsAuth.State.VerifiedChains[0][0].Subject.CommonName),
		},
	)
	if err != nil {
		c.log.Errorf("Error loading agents: %s", err)

		return nil, fmt.Errorf("internal server error")
	}

	if len(pgAgents) != 1 {
		c.log.Errorf("Found %d matching agents: %s", len(pgAgents), err)

		return nil, fmt.Errorf("unauthorized")
	}

	pgAgent := pgAgents[0]

	return &Session{
		AgentUID: pgAgent.UID,
	}, nil
}
