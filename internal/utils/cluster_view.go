package utils

import (
	"github.com/berops/claudie/proto/pb"
	"google.golang.org/protobuf/proto"
)

// ClusterView contains the per-cluster view on a given config.
// No mutex is needed when processing concurrently as long as each cluster only
// works with related values.
type ClusterView struct {
	// CurrentClusters are the individual clusters defined in the kubernetes section of the config of the current state.
	CurrentClusters map[string]*pb.K8Scluster
	// DesiredClusters are the individual clusters defined in the kubernetes section of the config of the desired state.
	DesiredClusters map[string]*pb.K8Scluster

	// Loadbalancers are the loadbalancers attach to a given kubernetes cluster in the current state.
	Loadbalancers map[string][]*pb.LBcluster
	// DesiredLoadbalancers are the loadbalancers attach to a given kubernetes cluster in the desired state.
	DesiredLoadbalancers map[string][]*pb.LBcluster

	// DeletedLoadbalancers are the loadbalancers that will be deleted (present in the current state but missing in the desired state)
	DeletedLoadbalancers map[string][]*pb.LBcluster

	// ClusterWorkflows is additional information per-cluster workflow (current stage of execution, if any error occurred etc..)
	ClusterWorkflows map[string]*pb.Workflow
}

func NewClusterView(config *pb.Config) *ClusterView {
	var (
		clusterWorkflows     = make(map[string]*pb.Workflow)
		clusters             = make(map[string]*pb.K8Scluster)
		desiredClusters      = make(map[string]*pb.K8Scluster)
		loadbalancers        = make(map[string][]*pb.LBcluster)
		desiredLoadbalancers = make(map[string][]*pb.LBcluster)
		deletedLoadbalancers = make(map[string][]*pb.LBcluster)
	)

	for _, current := range config.CurrentState.Clusters {
		clusters[current.ClusterInfo.Name] = current

		// store the cluster name with default workflow state.
		clusterWorkflows[current.ClusterInfo.Name] = &pb.Workflow{
			Stage:  pb.Workflow_NONE,
			Status: pb.Workflow_IN_PROGRESS,
		}
	}

	for _, desired := range config.DesiredState.Clusters {
		desiredClusters[desired.ClusterInfo.Name] = desired

		// store the cluster name with default workflow state.
		clusterWorkflows[desired.ClusterInfo.Name] = &pb.Workflow{
			Stage:  pb.Workflow_NONE,
			Status: pb.Workflow_IN_PROGRESS,
		}
	}

	for _, current := range config.CurrentState.LoadBalancerClusters {
		loadbalancers[current.TargetedK8S] = append(loadbalancers[current.TargetedK8S], current)
	}

	for _, desired := range config.DesiredState.LoadBalancerClusters {
		desiredLoadbalancers[desired.TargetedK8S] = append(desiredLoadbalancers[desired.TargetedK8S], desired)
	}

Lb:
	for _, current := range config.CurrentState.LoadBalancerClusters {
		for _, desired := range config.DesiredState.LoadBalancerClusters {
			if desired.ClusterInfo.Name == current.ClusterInfo.Name && desired.ClusterInfo.Hash == current.ClusterInfo.Hash {
				continue Lb
			}
		}
		deletedLoadbalancers[current.TargetedK8S] = append(deletedLoadbalancers[current.TargetedK8S], proto.Clone(current).(*pb.LBcluster))
	}

	return &ClusterView{
		CurrentClusters:      clusters,
		DesiredClusters:      desiredClusters,
		Loadbalancers:        loadbalancers,
		DesiredLoadbalancers: desiredLoadbalancers,
		DeletedLoadbalancers: deletedLoadbalancers,
		ClusterWorkflows:     clusterWorkflows,
	}
}

// MergeChanges propagates the changes made back to the config.
func (view *ClusterView) MergeChanges(config *pb.Config) {
	config.State = view.ClusterWorkflows

	for name, updated := range view.CurrentClusters {
		idx := GetClusterByName(name, config.CurrentState.Clusters)
		if idx < 0 {
			config.CurrentState.Clusters = append(config.CurrentState.Clusters, updated)
			continue
		}
		config.CurrentState.Clusters[idx] = updated
	}

	for name, updated := range view.DesiredClusters {
		idx := GetClusterByName(name, config.DesiredState.Clusters)
		if idx < 0 {
			config.DesiredState.Clusters = append(config.DesiredState.Clusters, updated)
			continue
		}
		config.DesiredState.Clusters[idx] = updated
	}

	for _, updated := range view.Loadbalancers {
		for _, lb := range updated {
			idx := GetLBClusterByName(lb.ClusterInfo.Name, config.CurrentState.LoadBalancerClusters)
			if idx < 0 {
				config.CurrentState.LoadBalancerClusters = append(config.CurrentState.LoadBalancerClusters, lb)
				continue
			}
			config.CurrentState.LoadBalancerClusters[idx] = lb
		}
	}

	for _, updated := range view.DesiredLoadbalancers {
		for _, lb := range updated {
			idx := GetLBClusterByName(lb.ClusterInfo.Name, config.DesiredState.LoadBalancerClusters)
			if idx < 0 {
				config.DesiredState.LoadBalancerClusters = append(config.DesiredState.LoadBalancerClusters, lb)
				continue
			}
			config.DesiredState.LoadBalancerClusters[idx] = lb
		}
	}
}

// AllClusters returns a slice of cluster all cluster names, from both the current state and desired state.
// This is useful to be abe to distinguish which clusters were deleted and which were not.
func (view *ClusterView) AllClusters() []string {
	clusters := make(map[string]struct{})

	for _, current := range view.CurrentClusters {
		clusters[current.ClusterInfo.Name] = struct{}{}
	}

	for _, desired := range view.DesiredClusters {
		clusters[desired.ClusterInfo.Name] = struct{}{}
	}

	c := make([]string, 0, len(clusters))
	for k := range clusters {
		c = append(c, k)
	}

	return c
}

func (view *ClusterView) SetWorkflowError(clusterName string, err error) {
	view.ClusterWorkflows[clusterName].Status = pb.Workflow_ERROR
	view.ClusterWorkflows[clusterName].Description = err.Error()
}

func (view *ClusterView) SetWorkflowDone(clusterName string) {
	view.ClusterWorkflows[clusterName].Status = pb.Workflow_DONE
	view.ClusterWorkflows[clusterName].Stage = pb.Workflow_NONE
	view.ClusterWorkflows[clusterName].Description = ""
}
