package smtp

import (
	"fmt"
	"sync"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/broadcaster"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/model"
)

// ServiceName defines the name of smtp service used for dependency injection
const ServiceName = "smtp"

// EmailAttachement defines the base model of an email attachemet
type EmailAttachement struct {
	Filename string
	Mimetype string
	Embedded bool
	Content  []byte
}

// Email defines the base model of an email
type Email struct {
	To           string
	Subject      string
	BodyPlain    string
	BodyHTML     string
	Attachements []EmailAttachement
}

// IService defines the interface of the smtp service
type IService interface {
	gousu.IService

	Open(config *model.PgNotifierV1ConfigParamsEmailSmtp) *Sender
}

// Service provides an smtp sender running in a separate thread
type Service struct {
	log             *logger.Log
	stopBroadcaster *broadcaster.Bool
	runningFuncs    sync.WaitGroup
}

var _ IService = (*Service)(nil)

// Name returns the name of the smtp service from ServiceName
func (s *Service) Name() string {
	return ServiceName
}

// Start starts the SMTP-Sender in a separate thread
func (s *Service) Start() error {
	return nil
}

// Stop stops the SMTP-Sender thread
func (s *Service) Stop() error {
	s.stopBroadcaster.Next(true)

	s.runningFuncs.Wait()

	return nil
}

// Health checks if the MailService is healthy
func (s *Service) Health() error {
	return nil
}

// NewService if the ServiceFactory for an initialized Service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log:             logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
		stopBroadcaster: broadcaster.NewBool(false),
	}
}

// Assert NewService fullfills gousu.ServiceFactory
var _ (gousu.ServiceFactory) = NewService
