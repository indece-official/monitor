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

package public

import (
	"fmt"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgCheckStatusV1StatusToAPICheckStatusV1Status(pgCheckStatus model.PgCheckStatusV1Status) (apipublic.CheckStatusV1Status, error) {
	switch pgCheckStatus {
	case model.PgCheckStatusV1StatusOK:
		return apipublic.CheckStatusV1StatusOK, nil
	case model.PgCheckStatusV1StatusWarn:
		return apipublic.CheckStatusV1StatusWARNING, nil
	case model.PgCheckStatusV1StatusCrit:
		return apipublic.CheckStatusV1StatusCRITICAL, nil
	case model.PgCheckStatusV1StatusUnkn:
		return apipublic.CheckStatusV1StatusUNKNOWN, nil
	default:
		return "", fmt.Errorf("invalid check status: %s", pgCheckStatus)
	}
}

func (c *Controller) mapReCheckStatusV1ToAPICheckStatusV1(reCheckStatus *model.ReCheckStatusV1) (*apipublic.CheckStatusV1, error) {
	var err error

	apiCheckStatus := &apipublic.CheckStatusV1{}

	apiCheckStatus.Uid = reCheckStatus.CheckStatusUID
	apiCheckStatus.Status, err = c.mapPgCheckStatusV1StatusToAPICheckStatusV1Status(reCheckStatus.Status)
	if err != nil {
		return nil, fmt.Errorf("error mapping source: %s", err)
	}
	apiCheckStatus.Message = reCheckStatus.Message
	apiCheckStatus.Data = reCheckStatus.Data
	apiCheckStatus.DatetimeCreated = reCheckStatus.DatetimeCreated

	return apiCheckStatus, nil
}
