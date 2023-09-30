package commands

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"testing-demo/internal/image"
	"testing-demo/internal/logging"
	"testing-demo/internal/middleware"
	"testing-demo/internal/user"
)

type ServerCmd struct {
	LogLevel string `long:"log-level" description:"The log level to run at" default:"debug"`
	ListenOn int    `long:"listen" short:"l" description:"The port you wish the service to listen on" default:"8080"`
}

func (s *ServerCmd) Execute(args []string) error {
	var srv http.Server

	logging.Configure(s.LogLevel, "eval.server", []string{"stdout"}, []string{"stdout"})

	logger := logging.GetLogger()

	idleConnsClosed := make(chan struct{})

	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
	)

	go func() {
		sig := <-signals
		logger.Info("Exiting application!")
		logger.Infof("Signal: %v", sig)
		if sig == syscall.SIGHUP || sig == syscall.SIGTERM || sig == syscall.SIGINT {
			if err := srv.Shutdown(context.Background()); err != nil {
				logger.Errorw("HTTP server Shutdown", zap.Error(err))
			}
			close(idleConnsClosed)
		}
		os.Exit(1)
	}()
	s.setupRoutes(&srv)
	<-idleConnsClosed

	return nil
}

func (s *ServerCmd) setupRoutes(srv *http.Server) http.Handler {
	logger := logging.GetLogger()
	r := gin.New()

	r.Use(ginzap.RecoveryWithZap(logging.GetRootLogger(), true))
	r.Use(middleware.ContextLogger(logging.GetLogger()))
	r.Use(middleware.RequestLogger(time.RFC3339, true))

	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.ErrorHandler())
	apiRouter.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Vela-Request-Id",
			"Accept",
			"Accept-Encoding",
			"DNT",
			"Origin",
			"User-Agent",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
		},
		MaxAge: 12 * time.Hour,
	}))
	user.RegisterHandlers(apiRouter)
	image.RegisterHandlers(apiRouter)

	logger.Infof("Starting on port: %d", s.ListenOn)

	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.ListenOn),
		Handler:      r,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	srv.SetKeepAlivesEnabled(false)
	logger.Fatalw("Can't start Web server", zap.Error(srv.ListenAndServe()))
	return r
}
