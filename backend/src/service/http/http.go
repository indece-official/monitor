package http

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"gopkg.in/guregu/null.v4"
)

const ServiceName = "http"

type IService interface {
	gousu.IService

	Query(ctx context.Context, method string, url string, headers map[string][]string, body null.String) error
}

type Service struct {
	log    *logger.Log
	client *resty.Client
}

var _ IService = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	s.client = resty.New()

	return nil
}

func (s *Service) Health() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
