package main

import (
	"GolangTraining/internal/http/rest"
	"GolangTraining/internal/logger"
	"GolangTraining/internal/metrics"
	"GolangTraining/internal/subscription"
	"context"
	"github.com/go-playground/validator/v10"
	"os"
	"sync"
)

type Server struct {
	sync.WaitGroup
	Config      *MainConfig
	RESTHandler *rest.Handler
	Prometheus  *metrics.Prometheus
	Logger      *logger.StandardLogger
}

func NewServer(cfg *MainConfig, logger *logger.StandardLogger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
	}
}

// Initialize is responsible for app initialization and wrapping required dependencies
func (s *Server) Initialize(ctx context.Context) error {
	v := validator.New()
	prometheus := metrics.NewPrometheus("gin")

	//rc, err := rds.CreateConnection(&s.Config.Redis)
	//if err != nil {
	//	return err
	//}
	//redisRep, err := redis.CreateRepository(rc)
	//if err != nil {
	//	return err
	//}

	//mysqlConn, err := msql.NewConnection(ctx, &s.Config.MySQL)
	//if err != nil {
	//	return err
	//}
	//mysqlRep, err := mysql.CreateRepository(mysqlConn)
	//if err != nil {
	//	return err
	//}

	service := subscription.CreateService(&s.Config.Service, s.Logger, nil, nil, prometheus, v)

	handler := rest.CreateHandler(service)

	s.Prometheus = prometheus
	s.RESTHandler = handler

	return nil
}

// Start starts the application in blocking mode
func (s *Server) Start(ctx context.Context) {
	// Start TimerBased CRON Jobs
	//go app.StartCronJobs(ctx)

	// Create Router for HTTP Server
	router := SetupRouter(s.RESTHandler, s.Config, s.Prometheus)
	//s.RESTHandler.Prometheus = prometheus
	//logrus.Info(prometheus)

	// Start REST Server in Blocking mode
	s.RESTHandler.Start(ctx, s.Config.Server.Port, router)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *Server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// const op = "app.gacefulshutdown"
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.RESTHandler.Stop()

	close(done)
}
