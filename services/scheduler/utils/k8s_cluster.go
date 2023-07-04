package utils

import (
	"fmt"
	"strings"

	"github.com/berops/claudie/internal/manifest"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
)

// CreateK8sCluster reads the unmarshalled manifest and creates desired state for Kubernetes clusters.
// Returns slice of *pb.K8Scluster if successful, nil otherwise
func CreateK8sCluster(unmarshalledManifest *manifest.Manifest) ([]*pb.K8Scluster, error) {
	var clusters []*pb.K8Scluster
	// Loop through clusters mentioned in the manifest
	for _, cluster := range unmarshalledManifest.Kubernetes.Clusters {
		// Generate variables
		newCluster := &pb.K8Scluster{
			ClusterInfo: &pb.ClusterInfo{
				Name: strings.ToLower(cluster.Name),
				Hash: utils.CreateHash(utils.HashLength),
			},
			Kubernetes: cluster.Version,
			Network:    cluster.Network,
		}

		// create node-pools
		controlNodePools, err := unmarshalledManifest.CreateNodepools(cluster.Pools.Control, true)
		if err != nil {
			return nil, fmt.Errorf("error while creating control nodepool for %s : %w", cluster.Name, err)
		}
		computeNodePools, err := unmarshalledManifest.CreateNodepools(cluster.Pools.Compute, false)
		if err != nil {
			return nil, fmt.Errorf("error while creating compute nodepool for %s : %w", cluster.Name, err)
		}

		mergedPools := utils.MergeMapsFunc(func(merged map[string]*pb.NodePool, key string, val *pb.NodePool) {
			if _, ok := merged[key]; !ok {
				merged[key] = val
				return
			}

			// if the same nodepool is used in both compute and control just append a string to distinguish between them.
			// In the validation the same nodepool cannot be re-used multiple times in either control or compute planes.
			// thus a simple check for `isControl` is suffucient.
			if val.IsControl {
				val.Name += "-control"
				merged[key].Name += "-compute"
			} else {
				merged[key].Name += "-control"
				val.Name += "-compute"
			}

			merged[val.Name] = val
		}, computeNodePools, controlNodePools)

		newCluster.ClusterInfo.NodePools = utils.MapValues(mergedPools)
		clusters = append(clusters, newCluster)
	}
	return clusters, nil
}

// UpdateK8sClusters updates the desired state of the kubernetes clusters based on the current state
// returns error if failed, nil otherwise
func UpdateK8sClusters(newConfig *pb.Config) error {
clusterDesired:
	for _, clusterDesired := range newConfig.DesiredState.Clusters {
		for _, clusterCurrent := range newConfig.CurrentState.Clusters {
			// Found current cluster with matching name
			if clusterDesired.ClusterInfo.Name == clusterCurrent.ClusterInfo.Name {
				updateClusterInfo(clusterDesired.ClusterInfo, clusterCurrent.ClusterInfo)
				if clusterCurrent.Kubeconfig != "" {
					clusterDesired.Kubeconfig = clusterCurrent.Kubeconfig
				}
				// Skip the checks bellow
				continue clusterDesired
			}
		}

		// No current cluster found with matching name, create keys
		if clusterDesired.ClusterInfo.PublicKey == "" {
			err := createSSHKeyPair(clusterDesired.ClusterInfo)
			if err != nil {
				return fmt.Errorf("error encountered while creating desired state for %s : %w", clusterDesired.ClusterInfo.Name, err)
			}
		}
	}
	return nil
}
