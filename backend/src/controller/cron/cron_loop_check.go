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

package cron

import (
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/utils"
)

func (c *Controller) getCheck(checkUID string) (*model.PgCheckV1, error) {
	c.mutexCheckers.Lock()
	defer c.mutexCheckers.Unlock()

	for _, pgCheck := range c.checks {
		if pgCheck.UID == checkUID {
			return pgCheck, nil
		}
	}

	return nil, fmt.Errorf("not matching check found")
}

func (c *Controller) getNotifier(notifierUID string) (*model.PgNotifierV1, error) {
	c.mutexNotifiers.Lock()
	defer c.mutexNotifiers.Unlock()

	for _, pgNotifier := range c.notifiers {
		if pgNotifier.UID == notifierUID {
			return pgNotifier, nil
		}
	}

	return nil, fmt.Errorf("not matching notifier found")
}

func (c *Controller) getChecker(checkerUID string) (*model.PgCheckerV1, error) {
	c.mutexCheckers.Lock()
	defer c.mutexCheckers.Unlock()

	for _, pgChecker := range c.checkers {
		if pgChecker.UID == checkerUID {
			return pgChecker, nil
		}
	}

	return nil, fmt.Errorf("not matching checker found")
}

func (c *Controller) getHost(hostUID string) (*model.PgHostV1, error) {
	c.mutexHosts.Lock()
	defer c.mutexHosts.Unlock()

	for _, pgHost := range c.hosts {
		if pgHost.UID == hostUID {
			return pgHost, nil
		}
	}

	return nil, fmt.Errorf("not matching host found")
}

func (c *Controller) getAgent(agentUID string) (*model.PgAgentV1, error) {
	c.mutexAgents.Lock()
	defer c.mutexAgents.Unlock()

	for _, pgAgent := range c.agents {
		if pgAgent.UID == agentUID {
			return pgAgent, nil
		}
	}

	return nil, fmt.Errorf("not matching agent found")
}

func (c *Controller) check(checkUID string) error {
	pgCheck, err := c.getCheck(checkUID)
	if err != nil {
		c.log.Warnf("Error loading check %s: %s", checkUID, err)

		return nil
	}

	pgChecker, err := c.getChecker(pgCheck.CheckerUID)
	if err != nil {
		c.log.Warnf("Error loading checker for check %s: %s", pgCheck.UID, err)

		return nil
	}

	pgAgent, err := c.getAgent(pgChecker.AgentUID)
	if err != nil {
		c.log.Warnf("Error loading agent for check %s: %s", pgCheck.UID, err)

		return nil
	}

	reAgentActionPayload := &model.ReAgentActionV1CheckPayload{}

	reAgentActionPayload.CheckUID = pgCheck.UID
	reAgentActionPayload.CheckerType = pgChecker.Type
	reAgentActionPayload.Params = []*model.ReAgentActionV1CheckPayloadParam{}
	for _, pgCheckParam := range pgCheck.Config.Params {
		reCheckParam := &model.ReAgentActionV1CheckPayloadParam{}

		reCheckParam.Name = pgCheckParam.Name
		reCheckParam.Value = pgCheckParam.Value

		reAgentActionPayload.Params = append(reAgentActionPayload.Params, reCheckParam)
	}

	reAgentActionPayload.TimeoutDuration = 30 * time.Second
	if pgCheck.Config.Timeout.Valid && pgCheck.Config.Timeout.String != "" {
		reAgentActionPayload.TimeoutDuration, err = time.ParseDuration(pgCheck.Config.Timeout.String)
		if err != nil {
			return fmt.Errorf("error parsing check timeout for action: %s", err)
		}
	} else if pgChecker.Capabilities.DefaultTimeout.Valid && pgChecker.Capabilities.DefaultTimeout.String != "" {
		reAgentActionPayload.TimeoutDuration, err = time.ParseDuration(pgChecker.Capabilities.DefaultTimeout.String)
		if err != nil {
			return fmt.Errorf("error parsing checker default timeout for action: %s", err)
		}
	}

	reAgentActionPayload.TimeoutAt = time.Now().Add(reAgentActionPayload.TimeoutDuration)

	reAgentAction := &model.ReAgentActionV1{}
	reAgentAction.Type = model.ReAgentActionV1TypeCheck
	reAgentAction.ActionUID, err = utils.UUID()
	if err != nil {
		return fmt.Errorf("error generating uid for action: %s", err)
	}
	reAgentAction.AgentUID = pgAgent.UID
	reAgentAction.Payload = reAgentActionPayload

	err = c.cacheService.PublishAgentAction(reAgentAction)
	if err != nil {
		return fmt.Errorf("error publishing agent action: %s", err)
	}

	err = c.cacheService.AddOpenAgentAction(reAgentAction)
	if err != nil {
		return fmt.Errorf("error adding open agent action: %s", err)
	}

	return nil
}
