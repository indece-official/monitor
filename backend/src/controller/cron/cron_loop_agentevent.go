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

func (c *Controller) processAgentEvent(reAgentEvent *model.ReAgentEventV1) error {
	ctx := context.Background()

	switch reAgentEvent.Type {
	case model.ReAgentEventV1TypeCheckResult:
		rePayload, ok := reAgentEvent.Payload.(*model.ReAgentEventV1CheckResultPayload)
		if !ok {
			return fmt.Errorf("invalid agent event payload")
		}

		reAgentAction, err := c.cacheService.GetOpenAgentAction(reAgentEvent.AgentUID, rePayload.ActionUID)
		if err != nil {
			return fmt.Errorf("error getting open agent action: %s", err)
		}

		if reAgentAction == nil {
			c.log.Infof("No active action found for action uid")

			return nil
		}

		err = c.cacheService.DeleteOpenAgentAction(reAgentEvent.AgentUID, rePayload.ActionUID)
		if err != nil {
			return fmt.Errorf("error deleting open agent action: %s", err)
		}

		pgCheckStatusData := map[string]string{}
		for _, rePayloadValue := range rePayload.Values {
			pgCheckStatusData[rePayloadValue.Name] = rePayloadValue.Value
		}

		var checkError error
		if rePayload.Error.Valid {
			checkError = fmt.Errorf("%s", rePayload.Error.String)
		}

		err = c.addCheckStatus(
			ctx,
			rePayload.CheckUID,
			rePayload.Message,
			checkError,
			pgCheckStatusData,
		)
		if err != nil {
			return fmt.Errorf("error adding check status: %s", err)
		}
	default:
		return fmt.Errorf("invalid agent event type %s", reAgentEvent.Type)
	}

	return nil
}

func (c *Controller) agentEventLoop() error {
	reAgentEvents, subscription, err := c.cacheService.SubscribeForAgentEvents()
	if err != nil {
		return fmt.Errorf("subscribing for system events failed: %s", err)
	}
	defer subscription.Unsubscribe()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for !c.stop {
		select {
		case <-ticker.C:
			continue
		case reAgentEvent, more := <-reAgentEvents:
			if !more {
				return fmt.Errorf("subscription closed the channel")
			}

			err := c.processAgentEvent(reAgentEvent)
			if err != nil {
				c.log.Errorf("Error processing agent event: %s", err)

				continue
			}
		}
	}

	return nil
}
