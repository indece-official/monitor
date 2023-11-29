package cron

import (
	"context"
	"fmt"

	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/template"
	"gopkg.in/guregu/null.v4"
)

func (c *Controller) sendHttpQuery(
	qctx context.Context,
	locale model.Locale,
	templateType template.TemplateType,
	params map[string]interface{},
	method model.PgNotifierV1ConfigParamsHttpMethod,
	url string,
	headers []*model.PgNotifierV1ConfigParamsHttpHeader,
	body null.String,
) error {
	methodStr := ""

	switch method {
	case model.PgNotifierV1ConfigParamsHttpMethodGet:
		methodStr = "GET"
	case model.PgNotifierV1ConfigParamsHttpMethodPost:
		methodStr = "POST"
	case model.PgNotifierV1ConfigParamsHttpMethodPut:
		methodStr = "PUT"
	default:
		return fmt.Errorf("unsupported http method: %s", method)
	}

	headersMap := map[string][]string{}
	for _, header := range headers {
		if len(headersMap[header.Name]) == 0 {
			headersMap[header.Name] = []string{}
		}

		headersMap[header.Name] = append(headersMap[header.Name], header.Value)
	}

	err := c.httpService.Query(
		qctx,
		methodStr,
		url,
		headersMap,
		body,
	)
	if err != nil {
		return err
	}

	return nil
}
