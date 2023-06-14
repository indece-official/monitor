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

package connector

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/monitor/backend/src/generated/model/apiconnector"
	"github.com/indece-official/monitor/backend/src/model"
	"github.com/indece-official/monitor/backend/src/service/cache"
	"github.com/indece-official/monitor/backend/src/service/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/namsral/flag"
)

// ControllerName defines the name of the api controller used for dependency exception
const ControllerName = "connector"

var (
	serverHost = flag.String("server_connector_host", "0.0.0.0", "")
	serverPort = flag.Int("server_connector_port", 9440, "")
)

// IController is the interface of the api controller
type IController interface {
	gousu.IController
}

// Controller is the admin api controller
type Controller struct {
	apiconnector.UnimplementedConnectorServer
	log             *logger.Log
	server          *grpc.Server
	postgresService postgres.IService
	cacheService    cache.IService
	listener        net.Listener
	error           error
	stop            bool
	waitGroupStop   sync.WaitGroup
}

var _ IController = (*Controller)(nil)

func (c *Controller) Name() string {
	return ControllerName
}

func (c *Controller) startServer() error {
	var err error

	ctx := context.Background()

	c.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", *serverHost, *serverPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	pgConfigProperties, err := c.postgresService.GetConfigProperties(
		ctx,
		&postgres.GetConfigPropertiesFilter{
			Keys: []model.PgConfigPropertyV1Key{
				model.PgConfigPropertyV1KeyTLSCaCrt,
				model.PgConfigPropertyV1KeyTLSServerCrt,
				model.PgConfigPropertyV1KeyTLSServerKey,
			},
		},
	)
	if err != nil {
		_ = c.listener.Close()
		c.listener = nil

		return fmt.Errorf("error loading config properties: %s", err)
	}

	if pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value == "" ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerCrt] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerCrt].Value == "" ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerKey] == nil ||
		pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerKey].Value == "" {
		_ = c.listener.Close()
		c.listener = nil

		c.log.Infof("Not starting connector server: no tls certs")

		return nil
	}

	serverCert, err := tls.X509KeyPair(
		[]byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerCrt].Value),
		[]byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSServerKey].Value),
	)
	if err != nil {
		return err
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM([]byte(pgConfigProperties[model.PgConfigPropertyV1KeyTLSCaCrt].Value)) {
		return fmt.Errorf("credentials: failed to append certificates")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{
			serverCert,
		},
		ClientCAs:  rootCAs,
		RootCAs:    rootCAs,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
	})
	if err != nil {
		_ = c.listener.Close()
		c.listener = nil

		return fmt.Errorf("failed to load tls certs: %s", err)
	}

	c.server = grpc.NewServer(grpc.Creds(creds))
	apiconnector.RegisterConnectorServer(c.server, c)
	go c.server.Serve(c.listener)

	return nil
}

func (c *Controller) stopServer() error {
	if c.server != nil {
		c.server.GracefulStop()
		c.server = nil
	}

	if c.listener != nil {
		err := c.listener.Close()
		if err != nil {
			return err
		}

		c.listener = nil
	}

	return nil
}

func (c *Controller) Start() error {
	err := c.startServer()
	if err != nil {
		return err
	}

	c.waitGroupStop.Add(1)
	go func() {
		for !c.stop {
			c.error = nil

			err := c.restartLoop()
			if err != nil {
				c.error = err

				c.log.Errorf("Error in restart loop: %s", err)
			}

			time.Sleep(5 * time.Second)
		}

		c.waitGroupStop.Done()
	}()

	return nil
}

// Health checks if the api server has thrown unresolvable internal errors
func (c *Controller) Health() error {
	return c.error
}

func (c *Controller) Stop() error {
	c.stop = true

	c.waitGroupStop.Wait()

	return c.stopServer()
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	log := logger.GetLogger(fmt.Sprintf("controller.%s", ControllerName))

	return &Controller{
		log:             log,
		postgresService: ctx.GetService(postgres.ServiceName).(postgres.IService),
		cacheService:    ctx.GetService(cache.ServiceName).(cache.IService),
	}
}

var _ gousu.ControllerFactory = NewController
