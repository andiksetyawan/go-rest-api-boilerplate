package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/db"
	"go-rest-api-boilerplate/pkg/observability/opentelemetry"
)

type server struct {
	handler http.Handler
	address string

	shutdownUptrace opentelemetry.ShutdownUptraceFunc
	//stopMetricPush   opentelemetry.StopMetricPushFunc
	//shutdownTracerEx opentelemetry.ShutDownExFunc
}

func NewServer() *server {
	config.Init()
	config.InitLogger()

	//init observability with uptrace only on production environment
	var shutdown opentelemetry.ShutdownUptraceFunc
	if config.App.ServiceEnvironment == "production" {
		shutdown = opentelemetry.InitUptrace(config.App.OtelUptraceDsn, config.App.ServiceName, "v1.0.0")
	}

	pgDb := db.NewPostgreeDb(config.App.DbHost, config.App.DbPort, config.App.DbName, config.App.DbUser, config.App.DbPass)
	dbCon := pgDb.Connect().GetConnection()

	handler := InitializedHandlerServer(dbCon)
	return &server{
		address:         config.App.ServiceAddress,
		handler:         handler,
		shutdownUptrace: shutdown,
	}
}

func (s *server) Run() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	srv := http.Server{Handler: s.handler, Addr: s.address}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("unable to create http listener")
			return
		}
	}()

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Error("failed to shutdown web server")
	} else {
		log.Info("http server has been shutdown")
	}

	if s.shutdownUptrace != nil {
		err = s.shutdownUptrace(ctx)
		if err != nil {
			log.WithError(err).Error("failed to shutdown uptrace pusher")
		} else {
			log.Info("uptrace pusher has been shutdown")
		}
	}

	//if s.shutdownTracerEx != nil {
	//	if err := s.shutdownTracerEx(ctxShutDown); err != nil {
	//		log.WithError(err).Error("failed to shutdown tracer exporter")
	//	}
	//	log.Info("tracer exporter has been shutdown")
	//}
	//
	//if s.stopMetricPush != nil {
	//	if err := s.stopMetricPush(ctxShutDown); err != nil {
	//		log.WithError(err).Error("failed to stoping pusher metrics")
	//	}
	//	log.Info("metrics pusher has been stopped")
	//}

	if err == nil {
		log.Info("done, server exited properly :)")
	}
}
