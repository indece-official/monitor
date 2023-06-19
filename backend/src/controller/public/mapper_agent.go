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
	"time"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgAgentV1ToAPIAgentV1(pgAgent *model.PgAgentV1, reStatus *model.ReAgentStatusV1) (*apipublic.AgentV1, error) {
	apiAgent := &apipublic.AgentV1{}

	apiAgent.Uid = pgAgent.UID
	apiAgent.HostUid = pgAgent.HostUID
	apiAgent.Type = pgAgent.Type.Ptr()
	apiAgent.Version = pgAgent.Version.Ptr()

	if !pgAgent.DatetimeRegistered.Valid || (reStatus != nil && reStatus.Status == model.ReAgentStatusV1StatusUnregistered) {
		apiAgent.Status = apipublic.AgentV1StatusUNREGISTERED
	} else if reStatus == nil || reStatus.Status == model.ReAgentStatusV1StatusReady {
		apiAgent.Status = apipublic.AgentV1StatusREADY
	} else if reStatus != nil && reStatus.Status == model.ReAgentStatusV1StatusError {
		apiAgent.Status = apipublic.AgentV1StatusERROR
	} else {
		apiAgent.Status = apipublic.AgentV1StatusUNKNOWN
	}

	if reStatus != nil {
		apiAgent.Connected = reStatus.DatetimeLastPing.Valid && time.Since(reStatus.DatetimeLastPing.Time) < 30*time.Second
		apiAgent.LastPing = reStatus.DatetimeLastPing.Ptr()
		apiAgent.Error = reStatus.Error.Ptr()
	} else {
		apiAgent.Connected = false
	}

	return apiAgent, nil
}

func (c *Controller) mapAPIAddAgentV1RequestBodyToPgAgentV1(requestBody *apipublic.V1AddAgentJSONRequestBody) (*model.PgAgentV1, error) {
	pgAgent := &model.PgAgentV1{}

	pgAgent.HostUID = requestBody.HostUid

	return pgAgent, nil
}

func (c *Controller) mapPgAgentV1ToAPIGetAgentsV1ResponseBody(pgAgents []*model.PgAgentV1, reStatuses map[string]*model.ReAgentStatusV1) (*apipublic.V1GetAgentsJSONResponseBody, error) {
	resp := &apipublic.V1GetAgentsJSONResponseBody{}

	resp.Agents = []apipublic.AgentV1{}

	for _, pgAgent := range pgAgents {
		apiAgent, err := c.mapPgAgentV1ToAPIAgentV1(pgAgent, reStatuses[pgAgent.UID])
		if err != nil {
			return nil, err
		}

		resp.Agents = append(resp.Agents, *apiAgent)
	}

	return resp, nil
}

func (c *Controller) mapPgAgentV1ToAPIGetAgentV1ResponseBody(pgAgent *model.PgAgentV1, reStatus *model.ReAgentStatusV1) (*apipublic.V1GetAgentJSONResponseBody, error) {
	resp := &apipublic.V1GetAgentJSONResponseBody{}

	apiAgent, err := c.mapPgAgentV1ToAPIAgentV1(pgAgent, reStatus)
	if err != nil {
		return nil, err
	}

	resp.Agent = *apiAgent

	return resp, nil
}

func (c *Controller) mapPgAgentV1ToAPIAddAgentV1ResponseBody(pgAgent *model.PgAgentV1, configFile string) (*apipublic.V1AddAgentJSONResponseBody, error) {
	resp := &apipublic.V1AddAgentJSONResponseBody{}

	resp.AgentUid = pgAgent.UID
	resp.ConfigFile = configFile

	return resp, nil
}
