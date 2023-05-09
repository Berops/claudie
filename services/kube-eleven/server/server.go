package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	kubeEleven "github.com/berops/claudie/services/kube-eleven/server/kube-eleven"
)

type server struct {
	pb.UnimplementedKubeElevenServiceServer
}

const (
	defaultKubeElevenPort = 50054 // default port for kube-eleven
)

// BuildCluster builds all cluster defined in the desired state
func (*server) BuildCluster(_ context.Context, req *pb.BuildClusterRequest) (*pb.BuildClusterResponse, error) {
	logger := utils.CreateLoggerWithProjectAndClusterName(req.ProjectName, req.Desired.ClusterInfo.Name)

	logger.Info().Msgf("Building kubernetes cluster")
	ke := kubeEleven.KubeEleven{
		K8sCluster: req.Desired,
		LBClusters: req.DesiredLbs,
	}

	if err := ke.BuildCluster(); err != nil {
		return nil, fmt.Errorf("error while building cluster %s project %s : %w", req.Desired.ClusterInfo.Name, req.ProjectName, err)
	}

	logger.Info().Msgf("Kubernetes cluster was successfully build")
	return &pb.BuildClusterResponse{Desired: req.Desired, DesiredLbs: req.DesiredLbs}, nil
}

func main() {
	// initialize logger
	utils.InitLog("kube-eleven")
	// Set KubeEleven port
	kubeElevenPort := utils.GetenvOr("KUBE_ELEVEN_PORT", fmt.Sprint(defaultKubeElevenPort))
	kubeElevenAddr := net.JoinHostPort("0.0.0.0", kubeElevenPort)
	lis, err := net.Listen("tcp", kubeElevenAddr)
	if err != nil {
		log.Fatal().Msgf("Failed to listen on %s : %v", kubeElevenAddr, err)
	}
	log.Info().Msgf("Kube-eleven service is listening on %s", kubeElevenAddr)

	s := grpc.NewServer()
	pb.RegisterKubeElevenServiceServer(s, &server{})

	// Add health service to gRPC
	healthServer := health.NewServer()
	// Kube-eleven does not have any custom health check functions, thus always serving.
	healthServer.SetServingStatus("kube-eleven-liveness", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("kube-eleven-readiness", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	g, ctx := errgroup.WithContext(context.Background())

	//goroutine for interrupt
	g.Go(func() error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		// wait for either the received signal or
		// check if an error occurred in other
		// go-routines.
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case sig := <-ch:
			log.Info().Msgf("Received signal %v", sig)
			err = errors.New("interrupt signal")
		}

		log.Info().Msg("Gracefully shutting down gRPC server")
		s.GracefulStop()
		healthServer.Shutdown()

		// Sometimes when the container terminates gRPC logs the following message:
		// rpc error: code = Unknown desc = Error: No such container: hash of the container...
		// It does not affect anything as everything will get terminated gracefully
		// this time.Sleep fixes it so that the message won't be logged.
		time.Sleep(1 * time.Second)

		return err
	})

	//server goroutine
	g.Go(func() error {
		// s.Serve() will create a service goroutine for each connection
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("kube-eleven failed to serve: %w", err)
		}
		log.Info().Msg("Finished listening for incoming connections")
		return nil
	})

	log.Info().Msgf("Stopping Kube-eleven: %s", g.Wait())
}
