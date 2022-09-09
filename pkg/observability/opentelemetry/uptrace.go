package opentelemetry

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
)

type ShutdownUptraceFunc func(ctx context.Context) error

//ForceFlush

//InitUptrace initialize uptrace
func InitUptrace(dsn, serviceName, version string) ShutdownUptraceFunc {
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		uptrace.WithDSN(dsn),
		//uptrace.WithDSN(""),
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(version),
		//opt...,
	)

	//TODO don't show the full dsn
	log.WithField("dsn", dsn).Infof("uptrace observability has been initialized")
	return func(ctx context.Context) error {
		log.Info("call shutdown")
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := uptrace.Shutdown(ctx); err != nil {
			log.WithError(err).Error("failed to shutdown TracerProvider")
		}
		return nil
	}
	//return func(ctx context.Context) error {
	//	cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
	//	defer cancel()
	//	if err := exporter.Shutdown(cxt); err != nil {
	//		//log.Info().Err(err).Msg("failed to shutdown TracerProvider")
	//		return err
	//	}
	//
	//	log.Info("tracer exporter has been shutdown")
	//	return nil
	//}
}
