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
	"time"

	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) restart() error {

	err := c.stopServer()
	if err != nil {
		return fmt.Errorf("error stopping server for restart: %s", err)
	}

	err = c.startServer()
	if err != nil {
		return fmt.Errorf("error restarting server: %s", err)
	}

	return nil
}

func (c *Controller) restartLoop() error {
	reSystemEvents, subscription, err := c.cacheService.SubscribeForSystemEvents()
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
		case reSystemEvent, more := <-reSystemEvents:
			if !more {
				return fmt.Errorf("subscription closed the channel")
			}

			switch reSystemEvent.Type {
			case model.ReSystemEventV1TypeCertsUpdated:
				err := c.restart()
				if err != nil {
					c.log.Errorf("Error restarting server: %s", err)

					continue
				}
			}
		}
	}

	return nil
}
