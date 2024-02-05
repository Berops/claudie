package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berops/claudie/internal/utils/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpc2 "google.golang.org/grpc"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/services/terraformer/server/adapters/inbound/grpc"
	outboundAdapters "github.com/berops/claudie/services/terraformer/server/adapters/outbound"
	"github.com/berops/claudie/services/terraformer/server/domain/usecases"
)

const (
	defaultPrometheusPort = "9097"
)

func main() {
	// Initialize logger
	utils.InitLog("terraformer")

	dynamoDBAdapter := outboundAdapters.CreateDynamoDBAdapter()
	stateAdapter := outboundAdapters.CreateS3Adapter()

	usecases := &usecases.Usecases{
		DynamoDB:          dynamoDBAdapter,
		StateStorage:      stateAdapter,
		SpawnProcessLimit: make(chan struct{}, usecases.SpawnProcessLimit),
	}

	metricsServer := &http.Server{Addr: fmt.Sprintf(":%s", utils.GetEnvDefault("PROMETHEUS_PORT", defaultPrometheusPort))}
	metrics.MustRegisterCounters()

	grpcAdapter := &grpc.GrpcAdapter{}
	grpcAdapter.Init(
		usecases,
		grpc2.UnaryInterceptor(metrics.MetricsMiddleware),
	)

	errGroup, errGroupContext := errgroup.WithContext(context.Background())
	errGroup.Go(grpcAdapter.Serve)

	// Check if terraformer microservice is in ready state every 30s
	errGroup.Go(func() error {
		ticker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-errGroupContext.Done():
				ticker.Stop()
				return nil

			case <-ticker.C:
				// If healthcheck result is positive then set the microservice as ready otherwise not ready
				if err := healthCheck(stateAdapter.Healthcheck, dynamoDBAdapter.Healthcheck); err != nil {
					grpcAdapter.HealthServer.SetServingStatus("terraformer-readiness", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
					log.Debug().Msgf("Failed to verify healthcheck: %v", err)
				} else {
					grpcAdapter.HealthServer.SetServingStatus("terraformer-readiness", grpc_health_v1.HealthCheckResponse_SERVING)
				}
			}
		}
	})

	errGroup.Go(func() error {
		shutdownSignalChan := make(chan os.Signal, 1)
		signal.Notify(shutdownSignalChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(shutdownSignalChan)

		var err error

		// Wait for either the received program termination signal or
		// check if an error occurred in other go-routines.
		select {
		case <-errGroupContext.Done():
			err = errGroupContext.Err()

		case interruptionSignal := <-shutdownSignalChan:
			log.Info().Msgf("Received program interruption signal %v", interruptionSignal)
			err = errors.New("interrupt signal")
		}

		if err := metricsServer.Shutdown(errGroupContext); err != nil {
			log.Err(err).Msgf("Failed to shutdown metrics server")
		}

		// Gracefully shutdown the gRPC adapter
		grpcAdapter.Stop()

		// Sometimes when the container terminates gRPC logs the following message:
		// rpc error: code = Unknown desc = Error: No such container: hash of the container...
		// It does not affect anything as everything will get terminated gracefully
		// this time.Sleep fixes it so that the message won't be logged.
		time.Sleep(1 * time.Second)

		return err
	})

	errGroup.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return metricsServer.ListenAndServe()
	})

	log.Info().Msgf("Stopping Terraformer: %v", errGroup.Wait())
}

// healthCheck function is a readiness function defined by terraformer
// it checks whether MinIO bucket exists and if dynamoDB table exists.
// If true, returns nil, error otherwise.
func healthCheck(minIOHealthcheck, dynamoDBHealthcheck func() error) error {
	err := minIOHealthcheck()
	if err != nil {
		return err
	}

	return dynamoDBHealthcheck()
}
