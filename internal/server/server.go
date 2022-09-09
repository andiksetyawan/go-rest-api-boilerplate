package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/db"
	httpTransport "go-rest-api-boilerplate/internal/transport/http"
	"go-rest-api-boilerplate/internal/usecase/repository"
	"go-rest-api-boilerplate/internal/usecase/service"
	"go-rest-api-boilerplate/pkg/logger"
	"go-rest-api-boilerplate/pkg/opentelemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type server struct {
	router  http.Handler
	address string

	shutdownUptrace opentelemetry.ShutdownUptraceFunc
	//stopMetricPush   opentelemetry.StopMetricPushFunc
	//shutdownTracerEx opentelemetry.ShutDownExFunc
}

func NewServer() *server {
	config.Init()
	logger.InitLogger()

	//init observability with uptrace only on production environment
	var shutdown opentelemetry.ShutdownUptraceFunc
	if config.App.ServiceEnvironment == "production" {
		shutdown = opentelemetry.InitUptrace(config.App.OtelUptraceDsn, config.App.ServiceName, "v1.0.0")
	}

	conf := config.App
	dbConn, _ := db.NewPostgreeDb(conf.DbHost, conf.DbPort, conf.DbName, conf.DbUser, conf.DbPass).Connect()

	userRepo := repository.NewUserRepository(dbConn)
	userSvc := service.NewUserService(userRepo)

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(conf.ServiceName))
	httpTransport.NewUserHandlerRegister(router, userSvc)
	log.Infoln(conf.ServiceAddress)
	fmt.Println("conf.ServiceAddress", conf.ServiceAddress)
	return &server{
		address:         config.App.ServiceAddress,
		router:          router,
		shutdownUptrace: shutdown,
	}
}

func (s *server) Run() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	srv := http.Server{Handler: s.router, Addr: s.address}

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
