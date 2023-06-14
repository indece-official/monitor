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

func (c *Controller) getConnectorByConnectorTypeAndHostUID(connectorType string, hostUID string) (*model.PgConnectorV1, error) {
	c.mutexConnectors.Lock()
	defer c.mutexConnectors.Unlock()

	for _, pgConnector := range c.connectors {
		if pgConnector.Type.Valid &&
			pgConnector.Type.String == connectorType &&
			pgConnector.HostUID == hostUID {
			return pgConnector, nil
		}
	}

	return nil, fmt.Errorf("not matching connector found")
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

	pgConnector, err := c.getConnectorByConnectorTypeAndHostUID(pgChecker.ConnectorType, pgCheck.HostUID)
	if err != nil {
		c.log.Warnf("Error loading connector for check %s: %s", pgCheck.UID, err)

		return nil
	}

	reConnectorActionPayload := &model.ReConnectorActionV1CheckPayload{}

	reConnectorActionPayload.CheckUID = pgCheck.UID
	reConnectorActionPayload.CheckerType = pgChecker.Type
	reConnectorActionPayload.Params = []*model.ReConnectorActionV1CheckPayloadParam{}
	for _, pgCheckParam := range pgCheck.Config.Params {
		reCheckParam := &model.ReConnectorActionV1CheckPayloadParam{}

		reCheckParam.Name = pgCheckParam.Name
		reCheckParam.Value = pgCheckParam.Value

		reConnectorActionPayload.Params = append(reConnectorActionPayload.Params, reCheckParam)
	}
	reConnectorActionPayload.TimeoutAt = time.Now().Add(10 * time.Second)
	reConnectorActionPayload.TimeoutDuration = 30 * time.Second

	reConnectorAction := &model.ReConnectorActionV1{}
	reConnectorAction.Type = model.ReConnectorActionV1TypeCheck
	reConnectorAction.ActionUID, err = utils.UUID()
	if err != nil {
		return fmt.Errorf("error generating uid for action: %s", err)
	}
	reConnectorAction.ConnectorUID = pgConnector.UID
	reConnectorAction.Payload = reConnectorActionPayload

	err = c.cacheService.PublishConnectorAction(reConnectorAction)
	if err != nil {
		return fmt.Errorf("error publishing connector action: %s", err)
	}

	err = c.cacheService.AddOpenConnectorAction(reConnectorAction)
	if err != nil {
		return fmt.Errorf("error adding open connector action: %s", err)
	}

	c.log.Infof("Triggered check %s on connector %s", pgCheck.UID, pgConnector.UID)

	return nil
}
