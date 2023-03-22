package main

import (
	"backend/internal/config"
	"backend/internal/resources"
	"backend/internal/version"
	"backend/pkg/logging"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Application logger initialized.")
	logger.Info("Start application...", logger.String("version", version.Version), logger.String("build_time", version.BuildTime), logger.String("commit", version.Commit))

	appConf := config.GetAppConfig(&logger)

	logger.Info("Application config read successful...")

	router := mux.NewRouter()

	resourceService, err := resources.NewService(&logger)

	if err != nil {
		fatalServer(err, &logger)
	}

	resourceHandler := resources.GetHandler(resourceService)
	resourceHandler.Register(router)

	logger.Info("Starting server:", logger.String("bind_addr", appConf.BindAddr))

	srv := &http.Server{
		Handler:      router,
		Addr:         appConf.BindAddr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fatalServer(err, &logger)
		}
	}(srv)

	//gracefull shutdown init here

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-c

	shutdown(&logger, srv)

}

func fatalServer(err error, l *logging.Logger) {
	l.Fatal(err.Error())
}

func shutdown(l *logging.Logger, srv *http.Server) {

	l.Info("Shutdown Application...")
	ctx, serverCancel := context.WithTimeout(context.Background(), 15*time.Second)

	err := srv.Shutdown(ctx)
	if err != nil {
		fatalServer(err, l)
	}
	serverCancel()
	l.Info("Application successful shutdown.")

}
