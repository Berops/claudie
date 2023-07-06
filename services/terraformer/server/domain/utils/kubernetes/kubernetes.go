package kubernetes

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	cluster_builder "github.com/berops/claudie/services/terraformer/server/domain/utils/cluster-builder"
)

type K8Scluster struct {
	ProjectName string

	DesiredState *pb.K8Scluster
	CurrentState *pb.K8Scluster

	// AttachedLBClusters are the LB clusters that are
	// attached to this K8s cluster.
	AttachedLBClusters []*pb.LBcluster
}

func (k *K8Scluster) Id() string {
	state := k.DesiredState
	if state == nil {
		state = k.CurrentState
	}

	return utils.GetClusterID(state.ClusterInfo)
}

func (k *K8Scluster) Build(logger zerolog.Logger) error {
	logger.Info().Msgf("Building K8S Cluster %s", k.DesiredState.ClusterInfo.Name)

	var currentClusterInfo *pb.ClusterInfo
	// Check if current cluster was defined, to avoid access of unreferenced memory
	if k.CurrentState != nil {
		currentClusterInfo = k.CurrentState.ClusterInfo
	}

	cluster := cluster_builder.ClusterBuilder{
		DesiredClusterInfo: k.DesiredState.ClusterInfo,
		CurrentClusterInfo: currentClusterInfo,

		ProjectName: k.ProjectName,
		ClusterType: pb.ClusterType_K8s,
		Metadata: map[string]any{
			"loadBalancers": k.AttachedLBClusters,
		},
	}

	err := cluster.CreateNodepools()
	if err != nil {
		return fmt.Errorf("error while creating the K8s cluster %s : %w", k.DesiredState.ClusterInfo.Name, err)
	}

	return nil
}

func (k *K8Scluster) Destroy(logger zerolog.Logger) error {
	logger.Info().Msgf("Destroying K8S Cluster %s", k.CurrentState.ClusterInfo.Name)
	cluster := cluster_builder.ClusterBuilder{
		CurrentClusterInfo: k.CurrentState.ClusterInfo,

		ProjectName: k.ProjectName,
		ClusterType: pb.ClusterType_K8s,
	}

	err := cluster.DestroyNodepools()
	if err != nil {
		return fmt.Errorf("error while destroying the K8s cluster %s : %w", k.CurrentState.ClusterInfo.Name, err)
	}

	return nil
}

func (k *K8Scluster) UpdateCurrentState() { k.CurrentState = k.DesiredState }
