package kuber

import (
	"context"
	"fmt"

	"github.com/Berops/platform/proto/pb"
	"github.com/rs/zerolog/log"
)

func SetUpStorage(c pb.KuberServiceClient, req *pb.SetUpStorageRequest) (*pb.SetUpStorageResponse, error) {
	res, err := c.SetUpStorage(context.Background(), req) //sending request to the server and receiving response
	if err != nil {
		return nil, fmt.Errorf("error while calling SetUpStorage on Kuber: %v", err)
	}
	log.Info().Msg("Storage was successfully set up")
	return res, nil
}

func StoreKubeconfig(c pb.KuberServiceClient, req *pb.StoreKubeconfigRequest) (*pb.StoreKubeconfigResponse, error) {
	res, err := c.StoreKubeconfig(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("Error while calling StoreKubeconfig on kuber: %v", err)
	}
	log.Info().Msg("Saved Kubeconfig Secret successfully")
	return res, nil
}
