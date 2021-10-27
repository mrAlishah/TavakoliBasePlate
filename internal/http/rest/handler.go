package rest

import (
	"GolangTraining/internal/subscription"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	VersionInfo struct {
		GitCommit       string
		BuildTime       string
		ContainerName   string
		BinaryStartTime time.Time
	}
	HTTPServer          *http.Server
	SubscriptionService subscription.Service
}

// CreateHandler Creates a new instance of REST handler
func CreateHandler(ss subscription.Service) *Handler {
	return &Handler{
		SubscriptionService: ss,
	}
}

// Start starts the http server
func (h *Handler) Start(ctx context.Context, port int32, r *gin.Engine) {
	const op = "http.rest.start"

	addr := fmt.Sprintf(":%d", port)

	h.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: r,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	logrus.Infof("[OK] Starting HTTP REST Server on %s ", addr)
	err := h.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.Fatal(errors.WithMessage(err, op))
	}
	// Code Reach Here after HTTP Server Shutdown!
	logrus.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop handles the http server in graceful shutdown
func (h *Handler) Stop() {
	const op = "http.rest.stop"

	// Create an 5s timeout context or waiting for app to shutdown after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	h.HTTPServer.SetKeepAlivesEnabled(false)
	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		logrus.Error(errors.WithMessage(err, op))
	}
	logrus.Info("HTTP REST Server graceful shutdown completed")
}
