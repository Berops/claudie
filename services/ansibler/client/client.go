package ansibler

import (
	"context"
	"fmt"

	"github.com/berops/claudie/proto/pb"
)

func RemoveClaudieUtilities(c pb.AnsiblerServiceClient, req *pb.RemoveClaudieUtilitiesRequest) (*pb.RemoveClaudieUtilitiesResponse, error) {
	res, err := c.RemoveClaudieUtilities(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling RemoveClaudieUtilities on Ansibler: %w", err)
	}
	return res, nil
}

// InstallVPN installs a Wireguard VPN on the nodes in the cluster and loadbalancers
func InstallVPN(c pb.AnsiblerServiceClient, req *pb.InstallRequest) (*pb.InstallResponse, error) {
	res, err := c.InstallVPN(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling InstallVPN on Ansibler: %w", err)
	}
	return res, nil
}

// InstallNodeRequirements install any additional requirements nodes in the cluster have (e.g. longhorn req, ...)
func InstallNodeRequirements(c pb.AnsiblerServiceClient, req *pb.InstallRequest) (*pb.InstallResponse, error) {
	res, err := c.InstallNodeRequirements(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling InstallNodeRequirements on Ansibler: %w", err)
	}
	return res, nil
}

// SetUpLoadbalancers ensures the nginx loadbalancer is set up correctly, with a correct DNS records
func SetUpLoadbalancers(c pb.AnsiblerServiceClient, req *pb.SetUpLBRequest) (*pb.SetUpLBResponse, error) {
	res, err := c.SetUpLoadbalancers(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling SetUpLoadbalancers on Ansibler: %w", err)
	}
	return res, nil
}

// TeardownLoadBalancers correctly destroys loadbalancers including selecting the new ApiServer endpoint
func TeardownLoadBalancers(c pb.AnsiblerServiceClient, req *pb.TeardownLBRequest) (*pb.TeardownLBResponse, error) {
	res, err := c.TeardownLoadBalancers(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling TeardownLoadBalancers on Ansibler: %w", err)
	}
	return res, nil
}

func UpdateAPIEndpoint(c pb.AnsiblerServiceClient, req *pb.UpdateAPIEndpointRequest) (*pb.UpdateAPIEndpointResponse, error) {
	res, err := c.UpdateAPIEndpoint(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling UpdateAPIEndpoint on Ansibler: %w", err)
	}
	return res, nil
}

func UpdateNoProxyEnvsInK8s(c pb.AnsiblerServiceClient, req *pb.UpdateNoProxyEnvsInK8SRequest) (*pb.UpdateNoProxyEnvsInK8SResponse, error) {
	res, err := c.UpdateNoProxyEnvsInK8S(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling UpdateNoProxyEnvsInK8s on Ansibler: %w", err)
	}
	return res, nil
}

func UpdateProxyEnvsOnNodes(c pb.AnsiblerServiceClient, req *pb.UpdateProxyEnvsOnNodesRequest) (*pb.UpdateProxyEnvsOnNodesResponse, error) {
	res, err := c.UpdateProxyEnvsOnNodes(context.Background(), req)
	if err != nil {
		return res, fmt.Errorf("error while calling UpdateProxyEnvsOnNodes on Ansibler: %w", err)
	}
	return res, nil
}
