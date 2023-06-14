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

package connector

import (
	"fmt"
	"io"

	"github.com/indece-official/monitor/backend/src/generated/model/apiconnector"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) processConnectorAction(stream apiconnector.Connector_CheckV1Server, reConnectorAction *model.ReConnectorActionV1) error {
	if reConnectorAction.Type != model.ReConnectorActionV1TypeCheck {
		return nil
	}

	rePayload := reConnectorAction.Payload.(*model.ReConnectorActionV1CheckPayload)

	req := &apiconnector.CheckV1Request{}
	req.ActionUID = reConnectorAction.ActionUID
	req.CheckUID = rePayload.CheckUID
	req.CheckerType = rePayload.CheckerType
	req.Params = []*apiconnector.CheckV1Param{}

	for _, reParam := range rePayload.Params {
		reqParam := &apiconnector.CheckV1Param{}

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

	c.log.Infof("Sent check request: %v", req)

	return nil
}

func (c *Controller) CheckV1(stream apiconnector.Connector_CheckV1Server) error {
	reSession, err := c.checkAuth(stream.Context())
	if err != nil {
		return err
	}

	c.log.Infof("Check connected")
	defer c.log.Infof("Check disconnected")

	reConnectorActions, subscription, err := c.cacheService.SubscribeForConnectorActions(reSession.ConnectorUID)
	if err != nil {
		return fmt.Errorf("error subscribing for connector actions: %s", err)
	}
	defer subscription.Unsubscribe()

	stop := make(chan bool)
	defer func() {
		stop <- true
	}()

	go func() {
		openReConnectorActions, err := c.cacheService.GetOpenConnectorActions(reSession.ConnectorUID)
		if err != nil {
			c.log.Errorf("Error loading open connector actions: %s", err)
		}

		for _, reConnectorAction := range openReConnectorActions {
			err = c.processConnectorAction(stream, reConnectorAction)
			if err != nil {
				c.log.Warnf("Error processing open connector action: %s", err)
			}
		}

		for {
			select {
			case reConnectorAction := <-reConnectorActions:
				err = c.processConnectorAction(stream, reConnectorAction)
				if err != nil {
					c.log.Warnf("Error processing connector action: %s", err)
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

		c.log.Infof("Received check response: %v", resp)

		reConnectorEvent := &model.ReConnectorEventV1{}
		reConnectorEvent.ConnectorUID = reSession.ConnectorUID
		reConnectorEvent.Type = model.ReConnectorEventV1TypeCheckResult

		rePayload := &model.ReConnectorEventV1CheckResultPayload{}

		rePayload.ActionUID = resp.ActionUID
		rePayload.CheckUID = resp.CheckUID
		rePayload.Message = resp.Message
		if resp.Error != "" {
			rePayload.Error.Scan(resp.Error)
		}
		rePayload.Values = []*model.ReConnectorEventV1CheckResultPayloadValue{}

		for _, respValue := range resp.Values {
			rePayloadValue := &model.ReConnectorEventV1CheckResultPayloadValue{}

			rePayloadValue.Name = respValue.Name
			rePayloadValue.Value = respValue.Value

			rePayload.Values = append(rePayload.Values, rePayloadValue)
		}

		reConnectorEvent.Payload = rePayload

		err = c.cacheService.PublishConnectorEvent(reConnectorEvent)
		if err != nil {
			c.log.Errorf("Error publishing connector event: %s", err)

			continue
		}
	}
}
