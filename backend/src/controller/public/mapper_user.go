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

	"golang.org/x/crypto/bcrypt"

	"github.com/indece-official/monitor/backend/src/generated/model/apipublic"
	"github.com/indece-official/monitor/backend/src/model"
)

func (c *Controller) mapPgUserV1SourceToAPIUserV1Source(pgUserSource model.PgUserV1Source) (apipublic.UserV1Source, error) {
	switch pgUserSource {
	case model.PgUserV1SourceLocal:
		return apipublic.LOCAL, nil
	default:
		return "", fmt.Errorf("invalid user source: %s", pgUserSource)
	}
}

func (c *Controller) mapPgUserV1RoleToAPIUserV1Role(pgUserRole model.PgUserV1Role) (apipublic.UserV1Role, error) {
	switch pgUserRole {
	case model.PgUserV1RoleShow:
		return apipublic.SHOW, nil
	case model.PgUserV1RoleAdmin:
		return apipublic.ADMIN, nil
	default:
		return "", fmt.Errorf("invalid user role: %s", pgUserRole)
	}
}

func (c *Controller) mapAPIUserV1RoleToPgUserV1Role(apiUserRole apipublic.UserV1Role) (model.PgUserV1Role, error) {
	switch apiUserRole {
	case apipublic.SHOW:
		return model.PgUserV1RoleShow, nil
	case apipublic.ADMIN:
		return model.PgUserV1RoleAdmin, nil
	default:
		return "", fmt.Errorf("invalid user role: %s", apiUserRole)
	}
}

func (c *Controller) mapPgUserV1ToAPIUserV1(pgUser *model.PgUserV1) (*apipublic.UserV1, error) {
	var err error

	apiUser := &apipublic.UserV1{}

	apiUser.Uid = pgUser.UID
	apiUser.Username = pgUser.Username
	apiUser.Name = pgUser.Name.Ptr()
	apiUser.Email = pgUser.Email.Ptr()
	apiUser.Source, err = c.mapPgUserV1SourceToAPIUserV1Source(pgUser.Source)
	if err != nil {
		return nil, fmt.Errorf("error mapping source: %s", err)
	}

	apiUser.Roles = []apipublic.UserV1Role{}
	for _, pgRole := range pgUser.LocalRoles {
		apiRole, err := c.mapPgUserV1RoleToAPIUserV1Role(pgRole)
		if err != nil {
			return nil, fmt.Errorf("error mapping role: %s", err)
		}

		apiUser.Roles = append(apiUser.Roles, apiRole)
	}

	return apiUser, nil
}

func (c *Controller) mapAPIAddUserV1RequestBodyToPgUserV1(requestBody *apipublic.V1AddUserJSONRequestBody) (*model.PgUserV1, error) {
	pgUser := &model.PgUserV1{}

	pgUser.Source = model.PgUserV1SourceLocal
	pgUser.Username = requestBody.Username

	if requestBody.Name != nil {
		pgUser.Name.Scan(*requestBody.Name)
	}

	if requestBody.Email != nil {
		pgUser.Email.Scan(*requestBody.Email)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("can't hash password: %s", err)
	}

	pgUser.PasswordHash.Scan(string(passwordHash))

	pgUser.LocalRoles = []model.PgUserV1Role{}
	for _, apiRole := range requestBody.Roles {
		pgRole, err := c.mapAPIUserV1RoleToPgUserV1Role(apiRole)
		if err != nil {
			return nil, fmt.Errorf("error mapping role: %s", err)
		}

		pgUser.LocalRoles = append(pgUser.LocalRoles, pgRole)
	}

	return pgUser, nil
}

func (c *Controller) mapAPIUpdateUserV1RequestBodyToPgUserV1(requestBody *apipublic.V1UpdateUserJSONRequestBody, oldPgUser *model.PgUserV1) (*model.PgUserV1, error) {
	tmp := *oldPgUser
	pgUser := tmp

	pgUser.Name.Scan(requestBody.Name)

	if requestBody.Email != nil {
		pgUser.Email.Scan(*requestBody.Email)
	}

	pgUser.LocalRoles = []model.PgUserV1Role{}
	for _, apiRole := range requestBody.Roles {
		pgRole, err := c.mapAPIUserV1RoleToPgUserV1Role(apiRole)
		if err != nil {
			return nil, fmt.Errorf("error mapping role: %s", err)
		}

		pgUser.LocalRoles = append(pgUser.LocalRoles, pgRole)
	}

	return &pgUser, nil
}

func (c *Controller) mapPgUserV1ToAPIGetUsersV1ResponseBody(pgUsers []*model.PgUserV1) (*apipublic.V1GetUsersJSONResponseBody, error) {
	resp := &apipublic.V1GetUsersJSONResponseBody{}

	resp.Users = []apipublic.UserV1{}

	for _, pgUser := range pgUsers {
		apiUser, err := c.mapPgUserV1ToAPIUserV1(pgUser)
		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, *apiUser)
	}

	return resp, nil
}

func (c *Controller) mapPgUserV1ToAPIGetOwnUserV1ResponseBody(pgUser *model.PgUserV1) (*apipublic.V1GetOwnUserJSONResponseBody, error) {
	resp := &apipublic.V1GetOwnUserJSONResponseBody{}

	apiUser, err := c.mapPgUserV1ToAPIUserV1(pgUser)
	if err != nil {
		return nil, err
	}

	resp.User = *apiUser

	return resp, nil
}

func (c *Controller) mapPgUserV1ToAPIGetUserV1ResponseBody(pgUser *model.PgUserV1) (*apipublic.V1GetUserJSONResponseBody, error) {
	resp := &apipublic.V1GetUserJSONResponseBody{}

	apiUser, err := c.mapPgUserV1ToAPIUserV1(pgUser)
	if err != nil {
		return nil, err
	}

	resp.User = *apiUser

	return resp, nil
}

func (c *Controller) mapPgUserV1ToAPIAddUserV1ResponseBody(pgUser *model.PgUserV1) (*apipublic.V1AddUserJSONResponseBody, error) {
	resp := &apipublic.V1AddUserJSONResponseBody{}

	resp.UserUid = pgUser.UID

	return resp, nil
}
