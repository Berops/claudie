package service

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/berops/claudie/services/manager/server/internal/store"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO: verify all paths and write tests.

func scheduleTasks(scheduled *store.Config) error {
	scheduledGRPC, err := store.ConvertToGRPC(scheduled)
	if err != nil {
		return fmt.Errorf("failed to convert database representation to GRPC for %q: %w", scheduled.Name, err)
	}

	for _, state := range scheduledGRPC.Clusters {
		var events []*spec.TaskEvent
		switch {
		// create
		case state.Current == nil && state.Desired != nil:
			events = append(events, &spec.TaskEvent{
				Id:        uuid.New().String(),
				Timestamp: timestamppb.New(time.Now().UTC()),
				Event:     spec.Event_CREATE,
				Task: &spec.Task{
					CreateState: &spec.CreateState{
						K8S: state.Desired.GetK8S(),
						Lbs: state.Desired.GetLoadBalancers(),
					},
				},
			})
		// delete
		case state.Desired == nil && state.Current != nil:
			events = append(events, &spec.TaskEvent{
				Id:        uuid.New().String(),
				Timestamp: timestamppb.New(time.Now().UTC()),
				Event:     spec.Event_DELETE,
				Task: &spec.Task{
					DeleteState: &spec.DeleteState{
						K8S: state.Current.GetK8S(),
						Lbs: state.Current.GetLoadBalancers(),
					},
				},
			})
		// update
		default:
			events = append(events, Diff(
				state.Current.K8S,
				state.Desired.K8S,
				state.Current.GetLoadBalancers().GetClusters(),
				state.Desired.GetLoadBalancers().GetClusters(),
			)...)
		}

		state.Events = &spec.Events{Events: events, Ttl: 0}
		state.State = &spec.Workflow{
			Stage:  spec.Workflow_NONE,
			Status: spec.Workflow_DONE,
		}
		// TODO: verify if not needed.
		//if e := state.Events; e == nil {
		//	state.Events = &spec.Events{Events: events, Ttl: 0}
		//} else {
		//	if len(e.Events) != 0 {
		//		return fmt.Errorf("failed to schedule tasks for cluster %q config %q. Cannot schedule tasks for a config which previous scheduled tasks have not been finished", cluster, scheduled.Name)
		//	}
		//	state.Events.Events = append(state.Events.Events, events...)
		//}
	}

	db, err := store.ConvertFromGRPC(scheduledGRPC)
	if err != nil {
		return fmt.Errorf("failed to convert GRPC representation to database for %q: %w", scheduled.Name, err)
	}

	*scheduled = *db
	return nil
}

// Diff takes the desired and current state to calculate difference between them to determine the difference and returns
// a number of tasks to be performed in specific order.
func Diff(current, desired *spec.K8Scluster, currentLbs, desiredLbs []*spec.LBcluster) []*spec.TaskEvent {
	k8sDynamic, k8sStatic := NodePoolNodes(current)
	lbsDynamic, lbsStatic := LbsNodePoolNodes(currentLbs)

	k8sDiffResult := k8sNodePoolDiff(k8sDynamic, k8sStatic, desired)
	lbsDiffResult := lbsNodePoolDiff(lbsDynamic, lbsStatic, desiredLbs)

	k8sAllDeletedNodes := make(map[string][]string)
	maps.Insert(k8sAllDeletedNodes, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(k8sAllDeletedNodes, maps.All(k8sDiffResult.deletedStatic))
	maps.Insert(k8sAllDeletedNodes, maps.All(k8sDiffResult.partialDeletedDynamic))
	maps.Insert(k8sAllDeletedNodes, maps.All(k8sDiffResult.partialDeletedStatic))

	var deletedLoadbalancers []*spec.LBcluster
	for _, current := range currentLbs {
		found := slices.ContainsFunc(desiredLbs, func(bcluster *spec.LBcluster) bool {
			return current.ClusterInfo.Name == bcluster.ClusterInfo.Name
		})
		if !found {
			deletedLoadbalancers = append(deletedLoadbalancers, current)
		}
	}

	var addedLoadBalancers []*spec.LBcluster
	for _, desired := range desiredLbs {
		found := slices.ContainsFunc(currentLbs, func(bcluster *spec.LBcluster) bool {
			return desired.ClusterInfo.Name == bcluster.ClusterInfo.Name
		})
		if !found {
			addedLoadBalancers = append(addedLoadBalancers, desired)
		}
	}

	var events []*spec.TaskEvent

	// will contain also the deleted nodes / nodepools if any.
	ir := craftK8sIR(k8sDiffResult, current, desired)

	// since lbs are not part of the k8s cluster no need to keep track
	// of any deletions, we simply just delete the infra.
	irLbs := lbClone(currentLbs)

	if k8sDiffResult.adding {
		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_UPDATE,
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: ir,
					Lbs: &spec.LoadBalancers{Clusters: irLbs}, // keep current lbs
				},
			},
		})
	}

	if target, deleted := deletedTargetApiNodePool(k8sDiffResult, current, currentLbs); deleted {
		lb := utils.FindLbAPIEndpointCluster(irLbs)

		var nextControlNodePool *spec.NodePool
		for _, np := range utils.FindControlNodepools(desired.ClusterInfo.NodePools) {
			if np.Name != target {
				nextControlNodePool = np
				break
			}
		}
		// No need to check if nextControlNodepool is nil. Validation of the inputmanifest
		// does not allow for the user to specify an empty list of control nodes
		nameWithoutHash := nextControlNodePool.Name

		// Each dynamic nodepool after the scheduler stage has a hash appended to it.
		// to get the original nodepool name as defined in the input manifest
		// we need to strip the hash.
		if nextControlNodePool.GetDynamicNodePool() != nil {
			nameWithoutHash = nextControlNodePool.Name[:len(nextControlNodePool.Name)-(utils.HashLength+1)] // +1 for '-'
		}

		for _, role := range lb.GetRoles() {
			if role.RoleType == spec.RoleType_ApiServer {
				role.TargetPools = slices.DeleteFunc(role.TargetPools, func(s string) bool { return s == target })
				role.TargetPools = append(role.TargetPools, nameWithoutHash)
				break
			}
		}

		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_UPDATE,
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: ir,
					Lbs: &spec.LoadBalancers{Clusters: irLbs},
				},
			},
		})
	}

	if k8sNodeApiEndpointDeleted(k8sDiffResult, current) {
		newApiNodePool := findNewAPIEndpointCandidate(desired.ClusterInfo.NodePools)

		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_UPDATE,
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{ApiNodePool: newApiNodePool},
			},
		})
	}

	if k8sDiffResult.deleting {
		dn := make(map[string]*spec.DeletedNodes)
		for k, v := range k8sAllDeletedNodes {
			dn[k] = &spec.DeletedNodes{Nodes: v}
		}
		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_DELETE,
			Task:      &spec.Task{DeleteState: &spec.DeleteState{Nodepools: dn}},
		})

		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_UPDATE,
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: desired,                              // since we don't work with the IR anymore will trigger the deletion of the infra.
					Lbs: &spec.LoadBalancers{Clusters: irLbs}, // no changes to the Lbs yet.
				},
			},
		})
	}

	// at last handle lb changes.
	// This will move the current state from an intermediate representation (if any)
	// to the desired as given in the manifest.
	lbsChanges := lbsDiffResult.adding || lbsDiffResult.deleting
	lbsChanges = lbsChanges || !proto.Equal(&spec.LoadBalancers{Clusters: currentLbs}, &spec.LoadBalancers{Clusters: desiredLbs})
	if lbsChanges || len(deletedLoadbalancers) > 0 || len(addedLoadBalancers) > 0 {
		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_UPDATE,
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: desired,
					Lbs: &spec.LoadBalancers{Clusters: desiredLbs},
				},
				DeleteState: func() *spec.DeleteState {
					if len(deletedLoadbalancers) > 0 {
						return &spec.DeleteState{Lbs: &spec.LoadBalancers{Clusters: deletedLoadbalancers}}
					}
					return nil
				}(),
			},
		})
	}

	if len(deletedLoadbalancers) > 0 {
		events = append(events, &spec.TaskEvent{
			Id:        uuid.New().String(),
			Timestamp: timestamppb.New(time.Now().UTC()),
			Event:     spec.Event_DELETE,
			Task: &spec.Task{
				DeleteState: &spec.DeleteState{Lbs: &spec.LoadBalancers{Clusters: deletedLoadbalancers}},
			},
		})
	}

	return events
}

type lbsNodePoolDiffResult struct {
	adding   bool
	deleting bool
}

func lbsNodePoolDiff(dynamic, static map[string]map[string][]string, desiredLbs []*spec.LBcluster) lbsNodePoolDiffResult {
	result := lbsNodePoolDiffResult{
		adding:   false,
		deleting: false,
	}

	for _, desired := range desiredLbs {
		for current := range dynamic[desired.GetClusterInfo().GetName()] {
			found := slices.ContainsFunc(desired.GetClusterInfo().GetNodePools(), func(pool *spec.NodePool) bool {
				return (pool.GetDynamicNodePool() != nil) && pool.Name == current
			})
			if !found {
				result.deleting = true
			}
		}

		for current := range static[desired.GetClusterInfo().GetName()] {
			found := slices.ContainsFunc(desired.GetClusterInfo().GetNodePools(), func(pool *spec.NodePool) bool {
				return (pool.GetStaticNodePool() != nil) && pool.Name == current
			})
			if !found {
				result.deleting = true
			}
		}

		for _, desiredNps := range desired.GetClusterInfo().GetNodePools() {
			if desiredNps.GetDynamicNodePool() != nil {
				current, ok := dynamic[desired.GetClusterInfo().GetName()][desiredNps.Name]
				if !ok {
					result.adding = true
					continue
				}

				if desiredNps.GetDynamicNodePool().Count > int32(len(current)) {
					result.adding = true
					continue
				}

				if desiredNps.GetDynamicNodePool().Count < int32(len(current)) {
					result.deleting = true
					// we don't need to keep track of which nodes are being deleted
					// as lbs are not part of k8s cluster.
				}
			} else {
				current, ok := static[desired.GetClusterInfo().GetName()][desiredNps.Name]
				if !ok {
					result.adding = true
					continue
				}

				if len(desiredNps.Nodes) > len(current) {
					result.adding = true
					continue
				}

				// TODO: this does not correctly handle the static node deletion.
				// For example if we would replace a ndoe with a different one it would fail.
				// TODO: issue for static nodes diff.
				// check for deletion of static nodes (also include replacement).
				for _, cnode := range current {
					var found bool
					for _, dnode := range desiredNps.Nodes {
						found = found || cnode == dnode.Name
					}
					if !found {
						result.deleting = true
						// we don't need to keep track of which nodes are being deleted
						// as lbs are not part of k8s cluster.
					}
				}
			}
		}
	}

	return result
}

type nodePoolDiffResult struct {
	partialDeletedDynamic map[string][]string
	partialDeletedStatic  map[string][]string
	deletedDynamic        map[string][]string
	deletedStatic         map[string][]string
	adding                bool
	deleting              bool
}

// k8sNodePoolDiff calculates difference between desired nodepools and current nodepools in a k8s cluster.
func k8sNodePoolDiff(dynamic, static map[string][]string, desiredCluster *spec.K8Scluster) nodePoolDiffResult {
	result := nodePoolDiffResult{
		partialDeletedDynamic: map[string][]string{},
		partialDeletedStatic:  map[string][]string{},
		deletedStatic:         map[string][]string{},
		deletedDynamic:        map[string][]string{},
		adding:                false,
		deleting:              false,
	}

	// look for whole dynamic nodepools deleted
	for currentNodePool := range dynamic {
		found := slices.ContainsFunc(desiredCluster.GetClusterInfo().GetNodePools(), func(pool *spec.NodePool) bool {
			return (pool.GetDynamicNodePool() != nil) && pool.Name == currentNodePool
		})
		if !found {
			result.deleting = true
			result.deletedDynamic[currentNodePool] = dynamic[currentNodePool]
		}
	}

	// look for whole static nodepools deleted
	for currentNodePool := range static {
		found := slices.ContainsFunc(desiredCluster.GetClusterInfo().GetNodePools(), func(pool *spec.NodePool) bool {
			return (pool.GetStaticNodePool() != nil) && pool.Name == currentNodePool
		})
		if !found {
			result.deleting = true
			result.deletedStatic[currentNodePool] = static[currentNodePool]
		}
	}

	// either both in current/desired but counts may differ or only in desired.
	for _, desired := range desiredCluster.GetClusterInfo().GetNodePools() {
		if desired.GetDynamicNodePool() != nil {
			current, ok := dynamic[desired.Name]
			if !ok {
				// not in current state, adding.
				result.adding = true
				continue
			}

			if desired.GetDynamicNodePool().Count > int32(len(current)) {
				result.adding = true
				continue
			}

			if desired.GetDynamicNodePool().Count < int32(len(current)) {
				result.deleting = true

				// chose nodes to delete.
				toDelete := int(int32(len(current)) - desired.GetDynamicNodePool().Count)
				for i := len(current) - 1; i >= len(current)-toDelete; i-- {
					result.partialDeletedDynamic[desired.Name] = append(result.partialDeletedDynamic[desired.Name], current[i])
				}
			}
		} else {
			current, ok := static[desired.Name]
			if !ok {
				// not in current state, adding.
				result.adding = true
				continue
			}
			if len(desired.Nodes) > len(current) {
				result.adding = true
				continue
			}

			// TODO: this does not correctly handle the static node deletion.
			// For example if we would replace a ndoe with a different one it would fail.
			// TODO: issue for static nodes diff.
			// check for deletion of static nodes (also include replacement).
			for _, cnode := range current {
				var found bool
				for _, dnode := range desired.Nodes {
					found = found || cnode == dnode.Name
				}
				if !found {
					result.deleting = true
					result.partialDeletedStatic[desired.Name] = append(result.partialDeletedStatic[desired.Name], cnode)
				}
			}
		}
	}
	return result
}

// NodePoolNodes returns the current nodes for the dynamic and static nodepools.
func NodePoolNodes(cluster *spec.K8Scluster) (map[string][]string, map[string][]string) {
	dynamic, static := make(map[string][]string), make(map[string][]string)

	for _, nodePool := range cluster.GetClusterInfo().GetNodePools() {
		if nodePool.GetDynamicNodePool() != nil {
			for _, node := range nodePool.Nodes {
				dynamic[nodePool.Name] = append(dynamic[nodePool.Name], node.Name)
			}
		}
		if nodePool.GetStaticNodePool() != nil {
			for _, node := range nodePool.Nodes {
				static[nodePool.Name] = append(static[nodePool.Name], node.Name)
			}
		}
	}

	return dynamic, static
}

func LbsNodePoolNodes(clusters []*spec.LBcluster) (map[string]map[string][]string, map[string]map[string][]string) {
	dynamic, static := make(map[string]map[string][]string), make(map[string]map[string][]string)

	for _, cluster := range clusters {
		dynamic[cluster.ClusterInfo.Name] = make(map[string][]string)
		static[cluster.ClusterInfo.Name] = make(map[string][]string)

		for _, nodepool := range cluster.GetClusterInfo().GetNodePools() {
			if nodepool.GetDynamicNodePool() != nil {
				for _, node := range nodepool.Nodes {
					dynamic[cluster.ClusterInfo.Name][nodepool.Name] = append(dynamic[cluster.ClusterInfo.Name][nodepool.Name], node.Name)
				}
			}
			if nodepool.GetStaticNodePool() != nil {
				for _, node := range nodepool.Nodes {
					static[cluster.ClusterInfo.Name][nodepool.Name] = append(static[cluster.ClusterInfo.Name][nodepool.Name], node.Name)
				}
			}
		}
	}

	return dynamic, static
}

func lbClone(desiredLbs []*spec.LBcluster) []*spec.LBcluster {
	var result []*spec.LBcluster
	for _, lb := range desiredLbs {
		result = append(result, proto.Clone(lb).(*spec.LBcluster))
	}
	return result
}

func craftK8sIR(k8sDiffResult nodePoolDiffResult, current, desired *spec.K8Scluster) *spec.K8Scluster {
	// Build the Intermediate Representation such that no deletion occurs in desired cluster.
	ir := proto.Clone(desired).(*spec.K8Scluster)

	for nodepool := range k8sDiffResult.partialDeletedDynamic {
		np := utils.GetNodePoolByName(nodepool, ir.ClusterInfo.NodePools)
		np.GetDynamicNodePool().Count = utils.GetNodePoolByName(nodepool, current.ClusterInfo.NodePools).GetDynamicNodePool().Count
	}

	for nodepool := range k8sDiffResult.partialDeletedStatic {
		np := utils.GetNodePoolByName(nodepool, ir.ClusterInfo.NodePools)
		np.Nodes = utils.GetNodePoolByName(nodepool, current.ClusterInfo.NodePools).Nodes
	}

	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	for nodepool := range deletedNodePools {
		np := utils.GetNodePoolByName(nodepool, current.ClusterInfo.NodePools)
		ir.ClusterInfo.NodePools = append(ir.ClusterInfo.NodePools, np)
	}

	return ir
}

func k8sNodeApiEndpointDeleted(k8sDiffResult nodePoolDiffResult, current *spec.K8Scluster) bool {
	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	for nodepool := range deletedNodePools {
		np := utils.GetNodePoolByName(nodepool, current.ClusterInfo.NodePools)
		if _, err := utils.FindAPIEndpointNode([]*spec.NodePool{np}); err == nil {
			return true
		}
	}
	return false
}

func deletedTargetApiNodePool(k8sDiffResult nodePoolDiffResult, current *spec.K8Scluster, currentLbs []*spec.LBcluster) (string, bool) {
	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	for nodepool := range deletedNodePools {
		np := utils.GetNodePoolByName(nodepool, current.ClusterInfo.NodePools)
		if utils.IsNodepoolOnlyTargetOfLbAPI(currentLbs, np) {
			return nodepool, true
		}
	}
	return "", false
}

func findNewAPIEndpointCandidate(desired []*spec.NodePool) string {
	for _, np := range desired {
		if np.IsControl {
			return np.Name
		}
	}
	panic("no suitable api endpoint replacement candidate found, malformed state.")
}
