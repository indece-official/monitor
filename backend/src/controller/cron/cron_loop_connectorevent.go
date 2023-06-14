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

func (c *Controller) processConnectorEvent(reConnectorEvent *model.ReConnectorEventV1) error {
	ctx := context.Background()

	c.log.Infof("Received connector event %s from connector %s", reConnectorEvent.Type, reConnectorEvent.ConnectorUID)

	switch reConnectorEvent.Type {
	case model.ReConnectorEventV1TypeCheckResult:
		rePayload, ok := reConnectorEvent.Payload.(*model.ReConnectorEventV1CheckResultPayload)
		if !ok {
			return fmt.Errorf("invalid connector event payload")
		}

		reConnectorAction, err := c.cacheService.GetOpenConnectorAction(reConnectorEvent.ConnectorUID, rePayload.ActionUID)
		if err != nil {
			return fmt.Errorf("error getting open connector action: %s", err)
		}

		if reConnectorAction == nil {
			c.log.Infof("No active action found for action uid")

			return nil
		}

		err = c.cacheService.DeleteOpenConnectorAction(reConnectorEvent.ConnectorUID, rePayload.ActionUID)
		if err != nil {
			return fmt.Errorf("error deleting open connector action: %s", err)
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
		return fmt.Errorf("invalid connector event type %s", reConnectorEvent.Type)
	}

	return nil
}

func (c *Controller) connectorEventLoop() error {
	reConnectorEvents, subscription, err := c.cacheService.SubscribeForConnectorEvents()
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
		case reConnectorEvent, more := <-reConnectorEvents:
			if !more {
				return fmt.Errorf("subscription closed the channel")
			}

			c.log.Infof("Received event from connector %s: %v", reConnectorEvent.ConnectorUID, reConnectorEvent)

			err := c.processConnectorEvent(reConnectorEvent)
			if err != nil {
				c.log.Errorf("Error processing connector event: %s", err)

				continue
			}
		}
	}

	return nil
}
