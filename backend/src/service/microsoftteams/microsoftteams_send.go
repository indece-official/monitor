package microsoftteams

import (
	"context"

	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func (s *Service) Send(ctx context.Context, webhookURL string, msg *adaptivecard.Message) error {
	return s.client.SendWithContext(ctx, webhookURL, msg)
}
