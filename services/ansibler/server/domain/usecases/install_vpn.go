package usecases

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	commonUtils "github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	"github.com/berops/claudie/services/ansibler/server/utils"
	"github.com/berops/claudie/services/ansibler/templates"
)

const wireguardPlaybookFilePath = "../../ansible-playbooks/wireguard.yml"

type VPNInfo struct {
	ClusterNetwork string
	// NodepoolsInfos is a slice with each element of type *DesiredClusterNodepoolsInfo.
	// Each element corresponds to a cluster (either a Kubernetes cluster or attached LB clusters).
	NodepoolsInfos []*NodepoolsInfo
}

// InstallVPN installs VPN between nodes in the k8s cluster and lb clusters
func (u *Usecases) InstallVPN(request *pb.InstallRequest) (*pb.InstallResponse, error) {
	logger := log.With().Str("project", request.ProjectName).Str("cluster", request.Desired.ClusterInfo.Name).Logger()
	logger.Info().Msgf("Installing VPN")

	vpnInfo := &VPNInfo{
		ClusterNetwork: request.Desired.Network,
		NodepoolsInfos: []*NodepoolsInfo{
			// Construct and add NodepoolsInfo for the Kubernetes cluster
			{
				Nodepools:      request.Desired.ClusterInfo.NodePools,
				PrivateKey:     request.Desired.ClusterInfo.PrivateKey,
				ClusterID:      fmt.Sprintf("%s-%s", request.Desired.ClusterInfo.Name, request.Desired.ClusterInfo.Hash),
				ClusterNetwork: request.Desired.Network,
			},
		},
	}
	// Construct and add NodepoolsInfo for each of the attached LB clusters
	for _, lbCluster := range request.DesiredLbs {
		vpnInfo.NodepoolsInfos = append(vpnInfo.NodepoolsInfos,
			&NodepoolsInfo{
				Nodepools:      lbCluster.ClusterInfo.NodePools,
				PrivateKey:     lbCluster.ClusterInfo.PrivateKey,
				ClusterID:      fmt.Sprintf("%s-%s", lbCluster.ClusterInfo.Name, lbCluster.ClusterInfo.Hash),
				ClusterNetwork: request.Desired.Network,
			},
		)
	}

	if err := installWireguardVPN(fmt.Sprintf("%s-%s", request.Desired.ClusterInfo.Name, request.Desired.ClusterInfo.Hash), vpnInfo); err != nil {
		logger.Err(err).Msgf("Error encountered while installing VPN")
		return nil, fmt.Errorf("error encountered while installing VPN for cluster %s project %s : %w", request.Desired.ClusterInfo.Name, request.ProjectName, err)
	}

	logger.Info().Msgf("VPN was successfully installed")
	return &pb.InstallResponse{Desired: request.Desired, DesiredLbs: request.DesiredLbs}, nil
}

// installWireguardVPN install wireguard VPN for all nodes in the infrastructure.
func installWireguardVPN(clusterID string, vpnInfo *VPNInfo) error {
	// Directory where files (required by Ansible) will be generated.
	clusterDirectory := filepath.Join(baseDirectory, outputDirectory, fmt.Sprintf("%s-%s", clusterID, commonUtils.CreateHash(commonUtils.HashLength)))
	if err := commonUtils.CreateDirectory(clusterDirectory); err != nil {
		return fmt.Errorf("failed to create directory %s : %w", clusterDirectory, err)
	}

	if err := assignPrivateIPs(getAllNodepools(vpnInfo.NodepoolsInfos), vpnInfo.ClusterNetwork); err != nil {
		return fmt.Errorf("error while setting the private IPs for %s : %w", clusterDirectory, err)
	}

	if err := utils.GenerateInventoryFile(templates.AllNodesInventoryTemplate, clusterDirectory,
		// Value of Ansible template parameters.
		AllNodesInventoryData{
			NodepoolsInfo: vpnInfo.NodepoolsInfos,
		},
	); err != nil {
		return fmt.Errorf("error while creating inventory file for %s : %w", clusterDirectory, err)
	}

	for _, nodepoolInfo := range vpnInfo.NodepoolsInfos {
		if err := commonUtils.CreateKeyFile(nodepoolInfo.PrivateKey, clusterDirectory, fmt.Sprintf("%s.%s", nodepoolInfo.ClusterID, sshPrivateKeyFileExtension)); err != nil {
			return fmt.Errorf("failed to create key file for %s : %w", nodepoolInfo.ClusterID, err)
		}
	}
	ansible := utils.Ansible{
		Playbook:  wireguardPlaybookFilePath,
		Inventory: utils.InventoryFileName,
		Directory: clusterDirectory,
	}
	if err := ansible.RunAnsiblePlaybook(fmt.Sprintf("VPN - %s", clusterID)); err != nil {
		return fmt.Errorf("error while running ansible for %s : %w", clusterDirectory, err)
	}

	return os.RemoveAll(clusterDirectory)
}

// getAllNodepools flattens []*DesiredClusterNodepoolsInfo to []*pb.NodePool.
// Returns a slice of all the nodepools.
func getAllNodepools(nodepoolsInfo []*NodepoolsInfo) []*pb.NodePool {
	var nodepools []*pb.NodePool
	for _, nodepoolInfo := range nodepoolsInfo {
		nodepools = append(nodepools, nodepoolInfo.Nodepools...)
	}

	return nodepools
}

// assignPrivateIPs will assign private IP addresses from the specified cluster network CIDR to all the nodes.
// Nodes which already have private IPs assigned will be ignored.
func assignPrivateIPs(nodepools []*pb.NodePool, cidr string) error {
	network, err := netip.ParsePrefix(cidr)
	if err != nil {
		return err
	}

	var (
		assignedPrivateIPs    = make(map[string]struct{})
		nodesWithoutPrivateIP []*pb.Node
	)

	// Construct nodesWithoutPrivateIP.
	for _, nodepool := range nodepools {
		for _, node := range nodepool.Nodes {
			if node.Private != "" {
				assignedPrivateIPs[node.Private] = struct{}{}
			} else {
				nodesWithoutPrivateIP = append(nodesWithoutPrivateIP, node)
			}
		}
	}

	for address := network.Addr().Next(); network.Contains(address) && len(nodesWithoutPrivateIP) > 0; address = address.Next() {
		// If private IP is already assigned to some node
		// then skip that IP.
		if _, ok := assignedPrivateIPs[address.String()]; ok {
			continue
		}

		// Otherwise assign it to the node.
		nodesWithoutPrivateIP[len(nodesWithoutPrivateIP)-1].Private = address.String()
		nodesWithoutPrivateIP = nodesWithoutPrivateIP[:len(nodesWithoutPrivateIP)-1]
	}

	if len(nodesWithoutPrivateIP) > 0 {
		return errors.New("failed to assign private IPs to all nodes")
	}

	return nil
}
