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
	"fmt"
	"io"

	"github.com/indece-official/monitor/backend/src/generated/model/apiagent"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) processAgentAction(stream apiagent.Agent_CheckV1Server, reAgentAction *model.ReAgentActionV1) error {
	if reAgentAction.Type != model.ReAgentActionV1TypeCheck {
		return nil
	}

	rePayload := reAgentAction.Payload.(*model.ReAgentActionV1CheckPayload)

	req := &apiagent.CheckV1Request{}
	req.ActionUID = reAgentAction.ActionUID
	req.CheckUID = rePayload.CheckUID
	req.CheckerType = rePayload.CheckerType
	req.Params = []*apiagent.CheckV1Param{}

	for _, reParam := range rePayload.Params {
		reqParam := &apiagent.CheckV1Param{}

		reqParam.Name = reParam.Name
		reqParam.Value = reParam.Value

		req.Params = append(req.Params, reqParam)
	}

	err := stream.Send(req)
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error sending data: %s", err)
	}

	return nil
}

func (c *Controller) CheckV1(stream apiagent.Agent_CheckV1Server) error {
	reSession, err := c.checkAuth(stream.Context())
	if err != nil {
		return err
	}

	c.log.Infof("Check connected")
	defer c.log.Infof("Check disconnected")

	reAgentActions, subscription, err := c.cacheService.SubscribeForAgentActions(reSession.AgentUID)
	if err != nil {
		return fmt.Errorf("error subscribing for agent actions: %s", err)
	}
	defer subscription.Unsubscribe()

	stop := make(chan bool)
	defer func() {
		stop <- true
	}()

	go func() {
		openReAgentActions, err := c.cacheService.GetOpenAgentActions(reSession.AgentUID)
		if err != nil {
			c.log.Errorf("Error loading open agent actions: %s", err)
		}

		for _, reAgentAction := range openReAgentActions {
			err = c.processAgentAction(stream, reAgentAction)
			if err != nil {
				c.log.Warnf("Error processing open agent action: %s", err)
			}
		}

		for {
			select {
			case reAgentAction := <-reAgentActions:
				err = c.processAgentAction(stream, reAgentAction)
				if err != nil {
					c.log.Warnf("Error processing agent action: %s", err)
				}
			case <-stop:
				return
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			c.log.Errorf("Error receiving data: %s", err)

			return fmt.Errorf("internal server error")
		}

		reAgentEvent := &model.ReAgentEventV1{}
		reAgentEvent.AgentUID = reSession.AgentUID
		reAgentEvent.Type = model.ReAgentEventV1TypeCheckResult

		rePayload := &model.ReAgentEventV1CheckResultPayload{}

		rePayload.ActionUID = resp.ActionUID
		rePayload.CheckUID = resp.CheckUID
		rePayload.Message = resp.Message
		if resp.Error != "" {
			rePayload.Error.Scan(resp.Error)
		}
		rePayload.Values = []*model.ReAgentEventV1CheckResultPayloadValue{}

		for _, respValue := range resp.Values {
			rePayloadValue := &model.ReAgentEventV1CheckResultPayloadValue{}

			rePayloadValue.Name = respValue.Name
			rePayloadValue.Value = respValue.Value

			rePayload.Values = append(rePayload.Values, rePayloadValue)
		}

		reAgentEvent.Payload = rePayload

		err = c.cacheService.PublishAgentEvent(reAgentEvent)
		if err != nil {
			c.log.Errorf("Error publishing agent event: %s", err)

			continue
		}
	}
}
