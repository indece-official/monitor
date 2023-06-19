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
	"context"
	"fmt"
	"time"

	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) checkActionTimeouts() error {
	var err error

	now := time.Now()
	ctx := context.Background()

	openReAgentActions, err := c.cacheService.GetAllOpenAgentActions()
	if err != nil {
		return fmt.Errorf("error getting open agent actions: %s", err)
	}

	for _, openReAgentAction := range openReAgentActions {
		if openReAgentAction.Type != model.ReAgentActionV1TypeCheck {
			continue
		}

		reActionPayload, ok := openReAgentAction.Payload.(*model.ReAgentActionV1CheckPayload)
		if !ok {
			c.log.Errorf("Invalid agent action payload")

			continue
		}

		if reActionPayload.TimeoutAt.Before(now) {
			err = c.cacheService.DeleteOpenAgentAction(openReAgentAction.AgentUID, openReAgentAction.ActionUID)
			if err != nil {
				return fmt.Errorf("error deleting open agent actions: %s", err)
			}

			err = c.addCheckStatus(
				ctx,
				reActionPayload.CheckUID,
				"",
				fmt.Errorf("timeout after %ds", reActionPayload.TimeoutDuration/time.Second),
				map[string]string{},
			)
			if err != nil {
				return fmt.Errorf("error adding check status: %s", err)
			}
		}
	}

	return nil
}

func (c *Controller) actionTimeoutLoop() error {
	for !c.stop {
		time.Sleep(2 * time.Second)

		err := c.checkActionTimeouts()
		if err != nil {
			c.log.Errorf("Error running action timeouts: %s", err)

			continue
		}
	}

	return nil
}
