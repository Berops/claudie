package outboundAdapters

import (
	"errors"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/berops/claudie/proto/pb"
	cbox "github.com/berops/claudie/services/context-box/client"
)

// Communicates with the gRPC server of context-box microservice
type ContextBoxConnector struct {
	connectionUri string

	// grpcConnection is the underlying gRPC connection to context-box microservice.
	grpcConnection *grpc.ClientConn

	// GrpcClient is a gRPC client connection to context-box microservice.
	GrpcClient pb.ContextBoxServiceClient
}

// Creates and returns an instance of the ContextBoxConnector struct
func NewContextBoxConnector(connectionUri string) *ContextBoxConnector {
	return &ContextBoxConnector{
		connectionUri: connectionUri,
	}
}

// Creates a gRPC connection to the context-box microservice.
// If the connection is established, then performs a healthcheck.
func (c *ContextBoxConnector) Connect() error {

	// Since the k8sSidecarNotificationsReceiver will be responding to incoming notifications we can't
	// use a blocking gRPC dial to the context-box service. Thus we default to a non-blocking
	// connection with a retry policy of ~4 seconds instead.
	interceptorOptions := []grpc_retry.CallOption{

		grpc_retry.WithBackoff(grpc_retry.BackoffExponentialWithJitter(4*time.Second, 0.2)),
		grpc_retry.WithMax(7),
		grpc_retry.WithCodes(codes.Unavailable),
	}

	grpcConnection, err := grpc.Dial(
		c.connectionUri,
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(interceptorOptions...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(interceptorOptions...)),
	)
	if err != nil {
		return err
	}

	c.grpcConnection = grpcConnection
	c.GrpcClient = pb.NewContextBoxServiceClient(grpcConnection)

	return c.PerformHealthCheck()

}

// PerformHealthCheck checks health of the underlying gRPC connection to context-box microservice
func (c *ContextBoxConnector) PerformHealthCheck() error {
	if c.grpcConnection.GetState() == connectivity.Shutdown {
		return errors.New("Unhealthy gRPC connection to context-box microservice")
	}

	return nil
}

// Fetches all configs present in context-box DB
func (c *ContextBoxConnector) GetAllConfigs() ([]*pb.Config, error) {

	response, err := cbox.GetAllConfigs(c.GrpcClient)
	if err != nil {
		return []*pb.Config{}, err
	}

	return response.GetConfigs(), nil
}

// Sends request to the context-box microservice, to save a config in context-box DB.
func (c *ContextBoxConnector) SaveConfig(config *pb.Config) error {

	_, err := cbox.SaveConfigFrontEnd(c.GrpcClient, &pb.SaveConfigRequest{Config: config})
	return err
}

// Sends request to the context-box microservice, to delete a config with the given id, from context-box DB.
func (c *ContextBoxConnector) DeleteConfig(id string) error {

	err := cbox.DeleteConfig(c.GrpcClient,
		&pb.DeleteConfigRequest{
			Id:   id,
			Type: pb.IdType_HASH,
		},
	)
	return err
}

// Closes the gRPC connection to context-box microservice
func (c *ContextBoxConnector) Disconnect() error {
	return c.grpcConnection.Close()
}
