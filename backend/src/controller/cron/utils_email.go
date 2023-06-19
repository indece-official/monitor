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

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/smtp"
	"github.com/indece-official/monitor/backend/src/service/template"
)

func (c *Controller) sendEmail(
	qctx context.Context,
	sender *smtp.Sender,
	locale model.Locale,
	templateType template.TemplateType,
	to string,
	params map[string]interface{},
) error {
	bodyHTML, err := c.templateService.Generate(
		locale,
		templateType,
		template.TemplatePartEmailBodyHTML,
		params,
	)
	if err != nil {
		return err
	}

	bodyText, err := c.templateService.Generate(
		locale,
		templateType,
		template.TemplatePartEmailBodyText,
		params,
	)
	if err != nil {
		return err
	}

	subject, err := c.templateService.Generate(
		locale,
		templateType,
		template.TemplatePartEmailSubject,
		params,
	)
	if err != nil {
		return err
	}

	email := &smtp.Email{}
	email.To = to
	email.Subject = subject
	email.BodyPlain = bodyText
	email.BodyHTML = bodyHTML

	err = sender.SendEmail(email)
	if err != nil {
		return err
	}

	c.log.Infof("Sent email %s to %s", templateType, to)

	return nil
}
