package microsoftteams

import (
	"context"
	"fmt"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
)

const ServiceName = "microsoftteams"

type IService interface {
	Send(ctx context.Context, webhookURL string, msg *adaptivecard.Message) error
}

type Service struct {
	log    *logger.Log
	client *goteamsnotify.TeamsClient
}

var _ IService = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	s.client = goteamsnotify.NewTeamsClient()

	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Health() error {
	return nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ gousu.ServiceFactory = NewService
