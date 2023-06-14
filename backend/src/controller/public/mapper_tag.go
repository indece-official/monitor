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
	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgTagV1ToAPITagV1(pgTag *model.PgTagV1) (*apipublic.TagV1, error) {
	apiTag := &apipublic.TagV1{}

	apiTag.Uid = pgTag.UID
	apiTag.Name = pgTag.Name
	apiTag.Color = pgTag.Color

	return apiTag, nil
}

func (c *Controller) mapAPIAddTagV1RequestBodyToPgTagV1(requestBody *apipublic.V1AddTagJSONRequestBody) (*model.PgTagV1, error) {
	pgTag := &model.PgTagV1{}

	pgTag.Name = requestBody.Name
	pgTag.Color = requestBody.Color

	return pgTag, nil
}

func (c *Controller) mapAPIUpdateTagV1RequestBodyToPgTagV1(requestBody *apipublic.V1UpdateTagJSONRequestBody, oldPgTag *model.PgTagV1) (*model.PgTagV1, error) {
	tmp := *oldPgTag
	pgTag := tmp

	pgTag.Name = requestBody.Name
	pgTag.Color = requestBody.Color

	return &pgTag, nil
}

func (c *Controller) mapPgTagV1ToAPIGetTagsV1ResponseBody(pgTags []*model.PgTagV1) (*apipublic.V1GetTagsJSONResponseBody, error) {
	resp := &apipublic.V1GetTagsJSONResponseBody{}

	resp.Tags = []apipublic.TagV1{}

	for _, pgTag := range pgTags {
		apiTag, err := c.mapPgTagV1ToAPITagV1(pgTag)
		if err != nil {
			return nil, err
		}

		resp.Tags = append(resp.Tags, *apiTag)
	}

	return resp, nil
}

func (c *Controller) mapPgTagV1ToAPIGetTagV1ResponseBody(pgTag *model.PgTagV1) (*apipublic.V1GetTagJSONResponseBody, error) {
	resp := &apipublic.V1GetTagJSONResponseBody{}

	apiTag, err := c.mapPgTagV1ToAPITagV1(pgTag)
	if err != nil {
		return nil, err
	}

	resp.Tag = *apiTag

	return resp, nil
}

func (c *Controller) mapPgTagV1ToAPIAddTagV1ResponseBody(pgTag *model.PgTagV1) (*apipublic.V1AddTagJSONResponseBody, error) {
	resp := &apipublic.V1AddTagJSONResponseBody{}

	resp.TagUid = pgTag.UID

	return resp, nil
}
