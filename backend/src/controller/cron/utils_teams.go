package cron

import (
	"context"

	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/template"
)

func (c *Controller) sendTeamsNotification(
	qctx context.Context,
	locale model.Locale,
	webhookURL string,
	templateType template.TemplateType,
	params map[string]interface{},
) error {
	title, err := c.templateService.Generate(
		locale,
		templateType,
		template.TemplatePartMicrosoftTeamsTitle,
		params,
	)
	if err != nil {
		return err
	}

	body, err := c.templateService.Generate(
		locale,
		templateType,
		template.TemplatePartMicrosoftTeamsBody,
		params,
	)
	if err != nil {
		return err
	}

	msg, err := adaptivecard.NewSimpleMessage(body, title, true)
	if err != nil {
		return err
	}

	err = c.microsoftteamsService.Send(qctx, webhookURL, msg)
	if err != nil {
		return err
	}

	c.log.Infof("Sent microsoft teams notification %s", templateType)

	return nil
}
