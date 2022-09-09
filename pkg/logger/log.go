package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go-rest-api-boilerplate/config"
)

func InitLogger() {
	log.Infof("environment: %s", config.App.ServiceEnvironment)
	if config.App.ServiceEnvironment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		)))
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}
}
