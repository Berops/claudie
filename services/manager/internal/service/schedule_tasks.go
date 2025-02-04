package service

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/berops/claudie/internal/clusters"
	"github.com/berops/claudie/internal/kubectl"
	"github.com/berops/claudie/internal/loggerutils"
	"github.com/berops/claudie/internal/nodepools"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/berops/claudie/services/manager/internal/store"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScheduleResult describes what has happened during the
// scheduling of the tasks.
type ScheduleResult uint8

const (
	// NoReschedule describes the case where the manifest should not be rescheduled again
	// after either error-ing or completing.
	NoReschedule ScheduleResult = iota
	// Reschedule describes the case where the manifest should be rescheduled again
	// after either error-ing or completing.
	Reschedule
	// NotReady describes the case where the manifest is not ready to be scheduled yet,
	// this is mostly related to the retry policies which can vary. For example if
	// an exponential retry policy is used the manifest will not be ready to be scheduled
	// until the specified number of Tick pass.
	NotReady
	// FinalRetry describes the case where a manifest had a retry policy to retry
	// rescheduling the manifest N times before giving up. FinalRetry states that
	// the manifest should be retried one last time before giving up.
	FinalRetry
)

func scheduleTasks(scheduled *store.Config) (ScheduleResult, error) {
	scheduledGRPC, err := store.ConvertToGRPC(scheduled)
	if err != nil {
		return NotReady, fmt.Errorf("failed to convert database representation to GRPC for %q: %w", scheduled.Name, err)
	}

	var result ScheduleResult

	for cluster, state := range scheduledGRPC.Clusters {
		logger := loggerutils.WithProjectAndCluster(scheduledGRPC.Name, cluster)

		var events []*spec.TaskEvent
		switch {
		case state.Current == nil && state.Desired == nil:
			// nothing to do (desired state was not build).
		// create
		case state.Current == nil && state.Desired != nil:
			// Choose initial api endpoint.
		clusters:
			for _, state := range scheduledGRPC.Clusters {
				for _, lb := range state.Desired.GetLoadBalancers().GetClusters() {
					if lb.HasApiRole() {
						lb.UsedApiEndpoint = true
						continue clusters
					}
				}
				nps := state.Desired.K8S.ClusterInfo.NodePools
				nodepools.FirstControlNode(nps).NodeType = spec.NodeType_apiEndpoint
			}

			events = append(events, &spec.TaskEvent{
				Id:          uuid.New().String(),
				Timestamp:   timestamppb.New(time.Now().UTC()),
				Event:       spec.Event_CREATE,
				Description: "creating cluster",
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
				Id:          uuid.New().String(),
				Timestamp:   timestamppb.New(time.Now().UTC()),
				Event:       spec.Event_DELETE,
				Description: "deleting cluster",
				Task: &spec.Task{
					DeleteState: &spec.DeleteState{
						K8S: state.Current.GetK8S(),
						Lbs: state.Current.GetLoadBalancers(),
					},
				},
			})
		// update
		default:
			k8sip, _, err := clusters.PingNodes(logger, state.Current)
			if err != nil {
				apply, e := tryReachK8sNodes(logger, k8sip, state)
				if !apply {
					result = NotReady
					break
				}
				events = append(events, e...)
				result = Reschedule
				break
			}

			if state.State.Status == spec.Workflow_ERROR {
				if tasks := state.Events.Events; len(tasks) != 0 && tasks[0].OnError.Do != nil {
					task := tasks[0]

					switch s := task.OnError.Do.(type) {
					case *spec.Retry_Repeat_:
						events = tasks

						if s.Repeat.Kind == spec.Retry_Repeat_EXPONENTIAL {
							if s.Repeat.RetryAfter > 0 {
								s.Repeat.RetryAfter--
								result = NotReady
								break
							}

							s.Repeat.CurrentTick <<= 1
							if s.Repeat.CurrentTick >= s.Repeat.StopAfter {
								// final retry before error-ing out.
								result = FinalRetry
								task.OnError.Do = nil
								break
							}

							s.Repeat.RetryAfter = s.Repeat.CurrentTick
						}

						result = Reschedule
						logger.Debug().Msgf("rescheduled for a retry of previously failed task with ID %q.", task.Id)
					case *spec.Retry_Rollback_:
						result = Reschedule
						events = s.Rollback.Tasks
						logger.Debug().Msgf("rescheduled for a rollback with task ID %q of previous failed task with ID %q.", events[0].Id, task.Id)
					default:
						result = NoReschedule
						logger.Debug().Msgf("has not been rescheduled for a retry on failure")
					}

					if result == Reschedule || result == NotReady || result == FinalRetry {
						break
					}
				}
			}

			ir, e, err := rollingUpdate(state.Current, state.Desired)
			if err != nil {
				return NotReady, err
			}

			events = append(events, e...)
			if len(events) != 0 {
				logger.Debug().
					Msgf("[%d] rolling updates scheduled for k8s cluster, to be performed before building the actual desired state, starting with task with ID %q.", len(events), events[0].Id)
				// First we will let claudie to work on the rolling update
				// to have the latest versions of the terraform manifests.
				// After that the manifest will be rescheduled again
				// to handle the diff between the new current state (with
				// updated terraform files) and the desired state as specified
				// in the Manifest.
				result = Reschedule
				// We set the desired state to the intermediate desired state which is the same as the
				// current state but with updated templates for k8s cluster. After this state is build
				// by the builder the config will be rescheduled again to actually reflect the changes
				// made. (if any by the user).
				state.Desired = ir
				break
			}

			ir, e, err = rollingUpdateLBs(state.Current, state.Desired)
			if err != nil {
				return NotReady, err
			}

			events = append(events, e...)
			if len(events) > 0 {
				logger.Debug().
					Msgf("[%d] rolling updates scheduled for attached lb clusters, to be performed before building the actual desired state, starting with task with ID %q.", len(events), events[0].Id)
				result = Reschedule
				state.Desired = ir
				break
			}

			events = append(events, Diff(
				state.Current.K8S,
				state.Desired.K8S,
				state.Current.GetLoadBalancers().GetClusters(),
				state.Desired.GetLoadBalancers().GetClusters(),
			)...)

			logger.Debug().Msgf("Scheduled final [%d] tasks to be worked on to build the desired state", len(events))
		}

		switch result {
		case Reschedule, NoReschedule, FinalRetry:
			// Events are going to be worked on, thus clear the Error state, if any.
			state.State = &spec.Workflow{Stage: spec.Workflow_NONE, Status: spec.Workflow_DONE}
		case NotReady:
		}

		state.Events = &spec.Events{Events: events}
	}

	db, err := store.ConvertFromGRPC(scheduledGRPC)
	if err != nil {
		return NotReady, fmt.Errorf("failed to convert GRPC representation to database for %q: %w", scheduled.Name, err)
	}

	*scheduled = *db
	return result, nil
}

// Diff takes the desired and current state to determine the difference and returns
// a number of tasks to be performed in specific order. It is expected that the current state actually represents
// the actual current state of the cluster and the desired state contains relevant data from the current state with
// the requested changes (i.e. deletion, addition of nodes) from the new config changes, (relevant data was transferred
// to desired state).
func Diff(current, desired *spec.K8Scluster, currentLbs, desiredLbs []*spec.LBcluster) []*spec.TaskEvent {
	k8sDynamic, k8sStatic := NodePoolNodes(current)
	lbsDynamic, lbsStatic := LbsNodePoolNodes(currentLbs)

	k8sDiffResult := k8sNodePoolDiff(k8sDynamic, k8sStatic, desired)
	lbsDiffResult := lbsNodePoolDiff(lbsDynamic, lbsStatic, desiredLbs)
	autoscalerConfigUpdated := k8sAutoscalerDiff(current, desired)
	labelsAnnotationsTaintsUpdated := labelsTaintsAnnotationsDiff(current, desired)

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

	if current.Kubernetes != desired.Kubernetes {
		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_UPDATE,
			Description: fmt.Sprintf("changing kubernetes version from %v to %v", current.Kubernetes, desired.Kubernetes),
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: &spec.K8Scluster{
						ClusterInfo:       current.ClusterInfo,
						Network:           current.Network,
						Kubeconfig:        current.Kubeconfig,
						Kubernetes:        desired.Kubernetes,
						InstallationProxy: current.InstallationProxy,
					},
					Lbs: &spec.LoadBalancers{Clusters: currentLbs},
				},
			},
		})
	}

	currProxySettings := &spec.InstallationProxy{
		Mode: "default",
	}
	if current.InstallationProxy != nil {
		currProxySettings = current.InstallationProxy
	}

	desiredProxySettings := &spec.InstallationProxy{
		Mode: "default",
	}
	if desired.InstallationProxy != nil {
		desiredProxySettings = desired.InstallationProxy
	}

	if currProxySettings.Mode != desiredProxySettings.Mode || currProxySettings.Endpoint != desiredProxySettings.Endpoint {
		// Proxy settings have been set or changed.
		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_UPDATE,
			Description: "changing installation proxy settings",
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: &spec.K8Scluster{
						ClusterInfo: current.ClusterInfo,
						Kubernetes:  current.Kubernetes,
						Network:     current.Network,
						InstallationProxy: &spec.InstallationProxy{
							Mode:     desired.InstallationProxy.Mode,
							Endpoint: desired.InstallationProxy.Endpoint,
						},
					},
					Lbs: &spec.LoadBalancers{Clusters: currentLbs},
				},
			},
		})
	}

	// will contain also the deleted nodes / nodepools if any.
	ir := craftK8sIR(k8sDiffResult, current, desired)

	if k8sDiffResult.adding {
		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_UPDATE,
			Description: "adding nodes to k8s cluster",
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					K8S: ir,
					Lbs: &spec.LoadBalancers{Clusters: currentLbs}, // keep current lbs
				},
			},
		})
	}

	// determine any changes to the api endpoint at the K8s level.
	if endpointNodeDeleted(k8sDiffResult, current) {
		nodePool, node := newAPIEndpointNodeCandidate(desired.ClusterInfo.NodePools)

		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_UPDATE,
			Description: fmt.Sprintf("moving endpoint from old control plane node to a new control plane node %s from nodepool %s", node, nodePool),
			Task: &spec.Task{
				UpdateState: &spec.UpdateState{
					EndpointChange: &spec.UpdateState_NewControlEndpoint{
						NewControlEndpoint: &spec.UpdateState_K8SEndpoint{
							Nodepool: nodePool,
							Node:     node,
						},
					},
				},
			},
			OnError: &spec.Retry{Do: &spec.Retry_Repeat_{Repeat: &spec.Retry_Repeat{
				Kind: spec.Retry_Repeat_ENDLESS,
			}}},
		})
	}

	// determine any changes to the api endpoint at the LB level.
	cid, did, change := clusters.DetermineLBApiEndpointChange(currentLbs, desiredLbs)
	applylbIr := deletedTargetApiNodePools(k8sDiffResult, current, currentLbs)
	// Manager can't handle the endpoint renamed case as it requires to be part of the workflow
	// where the new DNS hostname is generated, as it needs to be updated immediately after.
	// Every other case can be handled by the manager as a separate step.
	applylbIr = applylbIr || (change != spec.ApiEndpointChangeState_EndpointRenamed && change != spec.ApiEndpointChangeState_NoChange)
	if applylbIr {
		// will contain merged roles from current/desired state
		// and will include added loadbalancers if any.
		lbsir := craftLbsIR(currentLbs, desiredLbs, addedLoadBalancers)

		// options that adjusts the processing of the task.
		irOptions := uint64(0)
		if change == spec.ApiEndpointChangeState_DetachingLoadBalancer || change == spec.ApiEndpointChangeState_AttachingLoadBalancer {
			irOptions |= spec.ForceExportPort6443OnControlPlane
		}

		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_UPDATE,
			Description: "applying load balancer intermediate representation",
			Task: &spec.Task{
				Options: irOptions,
				UpdateState: &spec.UpdateState{
					K8S: ir,
					Lbs: &spec.LoadBalancers{Clusters: lbsir},
				},
			},
		})

		if change != spec.ApiEndpointChangeState_EndpointRenamed && change != spec.ApiEndpointChangeState_NoChange {
			events = append(events, &spec.TaskEvent{
				Id:          uuid.New().String(),
				Timestamp:   timestamppb.New(time.Now().UTC()),
				Event:       spec.Event_UPDATE,
				Description: fmt.Sprintf("performing API endpoint change, reason: %s", change.String()),
				Task: &spec.Task{
					Options: irOptions,
					UpdateState: &spec.UpdateState{
						EndpointChange: &spec.UpdateState_LbEndpointChange{
							LbEndpointChange: &spec.UpdateState_LbEndpoint{
								State:             change,
								CurrentEndpointId: cid,
								DesiredEndpointId: did,
							},
						},
					},
				},
				OnError: &spec.Retry{Do: &spec.Retry_Repeat_{Repeat: &spec.Retry_Repeat{
					Kind: spec.Retry_Repeat_ENDLESS,
				}}},
			})
		}
	}

	if k8sDiffResult.deleting {
		dn := make(map[string]*spec.DeletedNodes)
		for k, v := range k8sAllDeletedNodes {
			dn[k] = &spec.DeletedNodes{Nodes: v}
		}
		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_DELETE,
			Description: "deleting nodes from k8s cluster",
			Task: &spec.Task{
				DeleteState: &spec.DeleteState{Nodepools: dn},
			},
		})
	}

	if len(deletedLoadbalancers) > 0 {
		events = append(events, &spec.TaskEvent{
			Id:          uuid.New().String(),
			Timestamp:   timestamppb.New(time.Now().UTC()),
			Event:       spec.Event_DELETE,
			Description: "deleting loadbalancer infrastructure",
			Task: &spec.Task{
				DeleteState: &spec.DeleteState{Lbs: &spec.LoadBalancers{Clusters: deletedLoadbalancers}},
			},
		})
	}

	// as the last step commit to the changes requested in the desired state, as edge cases and other
	// manipulations have been done beforehand.
	// we match the infrastructure of the desired state.
	events = append(events, &spec.TaskEvent{
		Id:          uuid.New().String(),
		Timestamp:   timestamppb.New(time.Now().UTC()),
		Event:       spec.Event_UPDATE,
		Description: "commiting requested desired state",
		Task: &spec.Task{
			UpdateState: &spec.UpdateState{
				K8S: desired,
				Lbs: &spec.LoadBalancers{Clusters: desiredLbs},
			},
		},
	})

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

				// Node names are transferred over from current state based on the public IP.
				// Thus, at this point we can figure out based on nodes names which were deleted/added
				// see existing_state.go:transferStaticNodes
				for _, dnode := range desiredNps.Nodes {
					found := slices.ContainsFunc(current, func(s string) bool { return s == dnode.Name })
					if !found {
						result.adding = true
					}
				}

				// Node names are transferred over from current state based on the public IP.
				// Thus, at this point we can figure out based on nodes names which were deleted/added
				// see existing_state.go:transferStaticNodes
				for _, cnode := range current {
					found := slices.ContainsFunc(desiredNps.Nodes, func(dn *spec.Node) bool { return cnode == dn.Name })
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

			// Node names are transferred over from current state based on the public IP.
			// Thus, at this point we can figure out based on nodes names which were deleted/added
			// see existing_state.go:transferStaticNodes
			for _, dnode := range desired.Nodes {
				found := slices.ContainsFunc(current, func(s string) bool { return s == dnode.Name })
				if !found {
					result.adding = true
				}
			}

			// Node names are transferred over from current state based on the public IP.
			// Thus, at this point we can figure out based on nodes names which were deleted/added
			// see existing_state.go:transferStaticNodes
			for _, cnode := range current {
				found := slices.ContainsFunc(desired.Nodes, func(dn *spec.Node) bool { return cnode == dn.Name })
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

	clusterID := desired.ClusterInfo.Id()

	k := slices.Collect(maps.Keys(k8sDiffResult.partialDeletedDynamic))
	slices.Sort(k)

	for _, nodepool := range k {
		inp := nodepools.FindByName(nodepool, ir.ClusterInfo.NodePools)
		cnp := nodepools.FindByName(nodepool, current.ClusterInfo.NodePools)

		log.Debug().Str("cluster", clusterID).Msgf("nodes from dynamic nodepool %q were partially deleted, crafting ir to include them", nodepool)
		inp.GetDynamicNodePool().Count = cnp.GetDynamicNodePool().Count
		fillDynamicNodes(clusterID, cnp, inp)
	}

	k = slices.Collect(maps.Keys(k8sDiffResult.partialDeletedStatic))
	slices.Sort(k)

	for _, nodepool := range k {
		log.Debug().Str("cluster", clusterID).Msgf("nodes from static nodepool %q were partially deleted, crafting ir to include them", nodepool)
		inp := nodepools.FindByName(nodepool, ir.ClusterInfo.NodePools)
		cnp := nodepools.FindByName(nodepool, current.ClusterInfo.NodePools)

		is := inp.GetStaticNodePool()
		cs := cnp.GetStaticNodePool()

		maps.Insert(is.NodeKeys, maps.All(cs.NodeKeys))
		transferStaticNodes(clusterID, cnp, inp)

		for _, cn := range cnp.Nodes {
			if slices.Contains(k8sDiffResult.partialDeletedStatic[nodepool], cn.Name) {
				inp.Nodes = append(inp.Nodes, cn)
			}
		}
	}

	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	k = slices.Collect(maps.Keys(deletedNodePools))
	slices.Sort(k)

	for _, nodepool := range k {
		log.Debug().Str("cluster", clusterID).Msgf("nodepool %q  deleted, crafting ir to include it", nodepool)
		np := nodepools.FindByName(nodepool, current.ClusterInfo.NodePools)
		ir.ClusterInfo.NodePools = append(ir.ClusterInfo.NodePools, np)
	}

	return ir
}

func craftLbsIR(current, desired, added []*spec.LBcluster) []*spec.LBcluster {
	// 1. Create an IR from the current state.
	// as for why current state, the desired state could have different
	// node counts for the lbs, and we don't want to perform any changes
	// with  the infrastructure. The LB IR will be the current state +
	// the roles in desired state + any added loadbalancers so that we
	// can perform any migrations safely.
	ir := lbClone(current)

	for _, lb := range desired {
		i := clusters.IndexLoadbalancerById(lb.ClusterInfo.Id(), ir)
		if i < 0 {
			// desired cluster is not in the current state.
			continue
		}

		// 2. for each role defined in the desired state
		// we try to match it against the current state
		// if missing we add it, if modified we merge it.
		for _, desired := range lb.Roles {
			var found *spec.Role

			for _, current := range ir[i].Roles {
				if desired.Name == current.Name {
					found = current
					break
				}
			}

			if found == nil {
				ir[i].Roles = append(ir[i].Roles, proto.Clone(desired).(*spec.Role))
			} else {
				found.MergeTargetPools(desired)
			}
		}
	}
	// 3. add new lbs.
	return append(ir, lbClone(added)...)
}

func endpointNodeDeleted(k8sDiffResult nodePoolDiffResult, current *spec.K8Scluster) bool {
	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	for nodepool := range deletedNodePools {
		np := nodepools.FindByName(nodepool, current.ClusterInfo.NodePools)
		if np.EndpointNode() != nil {
			return true
		}
	}

	clear(deletedNodePools)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.partialDeletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.partialDeletedStatic))

	for nodepool, nodes := range deletedNodePools {
		np := nodepools.FindByName(nodepool, current.ClusterInfo.NodePools)
		for _, deleted := range nodes {
			i := slices.IndexFunc(np.Nodes, func(node *spec.Node) bool { return node.Name == deleted })
			if i < 0 {
				continue
			}
			if np.Nodes[i].NodeType == spec.NodeType_apiEndpoint {
				return true
			}
		}
	}

	return false
}

func deletedTargetApiNodePools(k8sDiffResult nodePoolDiffResult, current *spec.K8Scluster, currentLbs []*spec.LBcluster) bool {
	deletedNodePools := make(map[string][]string)
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedDynamic))
	maps.Insert(deletedNodePools, maps.All(k8sDiffResult.deletedStatic))

	var deleted []*spec.NodePool
	for np := range deletedNodePools {
		deleted = append(deleted, nodepools.FindByName(np, current.ClusterInfo.NodePools))
	}

	_, d := targetPoolsDeleted(currentLbs, deleted)
	return d
}

func newAPIEndpointNodeCandidate(desired []*spec.NodePool) (string, string) {
	for _, np := range desired {
		if np.IsControl {
			return np.Name, np.Nodes[0].Name
		}
	}
	panic("no suitable api endpoint replacement candidate found, malformed state.")
}

// targetPoolsDeleted check whether the LB API cluster target pools are among those that get deleted, if yes returns the names.
func targetPoolsDeleted(current []*spec.LBcluster, nps []*spec.NodePool) ([]string, bool) {
	for _, role := range clusters.FindAssignedLbApiEndpoint(current).GetRoles() {
		if role.RoleType != spec.RoleType_ApiServer {
			continue
		}

		var matches []string
		for _, targetPool := range role.TargetPools {
			idx := slices.IndexFunc(nps, func(pool *spec.NodePool) bool {
				name, _ := nodepools.MatchNameAndHashWithTemplate(targetPool, pool.Name)
				return name == targetPool
			})
			if idx >= 0 {
				matches = append(matches, nps[idx].Name)
			}
		}

		if len(matches) == len(role.TargetPools) {
			return matches, true
		}
	}
	return nil, false
}

func k8sAutoscalerDiff(current, desired *spec.K8Scluster) bool {
	cnp := make(map[string]*spec.DynamicNodePool)
	for _, np := range current.GetClusterInfo().GetNodePools() {
		if dyn := np.GetDynamicNodePool(); dyn != nil {
			cnp[np.Name] = dyn
		}
	}

	for _, np := range desired.GetClusterInfo().GetNodePools() {
		if dyn := np.GetDynamicNodePool(); dyn != nil {
			if prev, ok := cnp[np.Name]; ok {
				equal := proto.Equal(prev.AutoscalerConfig, dyn.AutoscalerConfig)
				if !equal {
					return true
				}
			}
		}
	}

	return false
}

func labelsTaintsAnnotationsDiff(current, desired *spec.K8Scluster) bool {
	cnp := make(map[string]*spec.NodePool)
	for _, np := range current.GetClusterInfo().GetNodePools() {
		cnp[np.Name] = np
	}

	for _, np := range desired.GetClusterInfo().GetNodePools() {
		if prev, ok := cnp[np.Name]; ok {
			if !reflect.DeepEqual(prev.Annotations, np.Annotations) {
				return true
			}
			if !reflect.DeepEqual(prev.Labels, np.Labels) {
				return true
			}
			if !reflect.DeepEqual(prev.Taints, np.Taints) {
				return true
			}
		}
	}

	return false
}

// tryReachK8sNodes determines if the InputManifest should be rescheduled or not based on the desired state and the reachability of the
// kubernetes nodes of the cluster. If the InputManifest is not ready to be scheduled yet apply will be false. Only if apply is true
// will the function also returns events that need to be handled before any other.
func tryReachK8sNodes(logger zerolog.Logger, ips []string, state *spec.ClusterState) (apply bool, events []*spec.TaskEvent) {
	logger.Info().Msgf("%v kubernetes cluster nodes are unreachable", ips)

	state.State = &spec.Workflow{
		Stage:       spec.Workflow_NONE,
		Status:      spec.Workflow_ERROR,
		Description: fmt.Sprintf("%v kubernetes nodes are unreachable", ips),
	}

	kubectl := kubectl.Kubectl{
		Kubeconfig:        state.Current.K8S.Kubeconfig,
		MaxKubectlRetries: 5,
	}

	n, err := kubectl.KubectlGetNodeNames()
	if err != nil {
		logger.Err(err).Msgf("failed to retrieve actuall nodes present in the cluster, retrying later")
		return
	}

	nodesInCluster := make(map[string]struct{})
	for _, n := range strings.Split(string(n), "\n") {
		nodesInCluster[n] = struct{}{}
	}

	// look which out of the current nodes that claudie tracks is not present
	// inside the actual kubernetes cluster.
	toDelete := make(map[string]*spec.DeletedNodes)
	for _, np := range state.Current.GetK8S().GetClusterInfo().GetNodePools() {
		for _, n := range np.Nodes {
			// node name inside k8s cluster have stripped cluster prefix.
			name := strings.TrimPrefix(n.Name, fmt.Sprintf("%s-", state.Current.GetK8S().GetClusterInfo().Id()))
			if _, ok := nodesInCluster[name]; !ok {
				if _, ok := toDelete[np.Name]; !ok {
					toDelete[np.Name] = new(spec.DeletedNodes)
				}
				toDelete[np.Name].Nodes = append(toDelete[np.Name].Nodes, n.Name)
			}
		}
	}

	if len(toDelete) < 1 {
		logger.Info().Msgf("Fix the unreachable nodes or deleted them manually from the cluster via 'kubectl'")
		return
	}

	// delete node from desired state so that
	// the manifest would not be rescheduled to
	// again include the same node.
	// For Dynamic nodes this will trigger a new
	// workflow that will trigger a build of a new node
	// since we dont change the count.
	// and for static nodes this will just remove the node
	// from the state and not reschedule again.
	// if it is a static node and is still present require it to be deleted from the desired state.
	// for dynamic nodepools this is not needed.
	var fixNeeded bool
deletion:
	for np, d := range toDelete {
		for _, nodeName := range d.Nodes {
			static, node := nodepools.FindNode(state.Desired.K8S.ClusterInfo.NodePools, nodeName)
			if node != nil && static { // we found the node that was delete from the kubernetes cluster in the desired state.
				logger.Info().
					Msgf("static node %s with endpoint %s from nodepool %s needs to be removed from the desired state. Remove the node from the InputManifest and apply the changes", nodeName, node.Public, np)
				fixNeeded = true
				break deletion
			}
			logger.Info().Msgf("node %s from nodepool %s no longer part of the kubernetes cluster, scheduling deletion", nodeName, np)
		}
	}

	if fixNeeded {
		return // Wait for the user to fix the InputManifest.
	}

	events = append(events, &spec.TaskEvent{
		Id:          uuid.New().String(),
		Timestamp:   timestamppb.New(time.Now().UTC()),
		Event:       spec.Event_DELETE,
		Description: "deleting nodes from k8s cluster",
		Task: &spec.Task{
			DeleteState: &spec.DeleteState{Nodepools: toDelete},
		},
	})

	apply = true
	return
}
