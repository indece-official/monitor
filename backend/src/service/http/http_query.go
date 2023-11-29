package http

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

func (s *Service) Query(ctx context.Context, method string, url string, headers map[string][]string, body null.String) error {
	req := s.client.R()

	req.SetContext(ctx)
	req.SetHeaderMultiValues(headers)

	if body.Valid {
		req.SetBody(body.String)
	}

	resp, err := req.Execute(method, url)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("request failed %d %s - %s", resp.StatusCode(), resp.Status(), resp.Body())
	}

	return nil
}
