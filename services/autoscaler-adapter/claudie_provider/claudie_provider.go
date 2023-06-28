package claudie_provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/externalgrpc/protos"

	"github.com/berops/claudie/internal/envs"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	"github.com/berops/claudie/services/autoscaler-adapter/node_manager"
	frontend "github.com/berops/claudie/services/frontend/client"
)

const (
	// Default GPU label.
	GpuLabel = "claudie.io/gpu-node"
)

var (
	//Error for functions which are not implemented.
	ErrNotImplemented = errors.New("not implemented")
)

type nodeCache struct {
	// Nodegroup as per Cluster Autoscaler definition.
	nodeGroup *protos.NodeGroup
	// Nodepool as per Claudie definition.
	nodepool *pb.NodePool
	// Target size of node group.
	targetSize int32
}

type ClaudieCloudProvider struct {
	protos.UnimplementedCloudProviderServer

	// Name of the Claudie config.
	projectName string
	// Kubernetes InputManifest resource name
	resourceName string
	// Kubernetes InputManifest resource namespace
	resourceNamespace string
	// Cluster as described in Claudie config.
	configCluster *pb.K8Scluster
	// Map of cached info regarding nodes.
	nodesCache map[string]*nodeCache
	// Node manager.
	nodeManager *node_manager.NodeManager
	// Server mutex
	lock sync.Mutex
}

// NewClaudieCloudProvider returns a ClaudieCloudProvider with initialised caches.
func NewClaudieCloudProvider(projectName, clusterName string) *ClaudieCloudProvider {
	// Connect to Claudie and retrieve *pb.K8Scluster
	var (
		cluster    *pb.K8Scluster
		err        error
		rName      string
		rNamespace string
		nm         *node_manager.NodeManager
	)
	if cluster, rName, rNamespace, err = getClaudieState(projectName, clusterName); err != nil {
		panic(fmt.Sprintf("Error while getting cluster %s : %v", clusterName, err))
	}
	if nm, err = node_manager.NewNodeManager(cluster.ClusterInfo.NodePools); err != nil {
		panic(fmt.Sprintf("Error while creating node manager : %v", err))
	}
	// Initialise all other variables.
	log.Logger = log.Logger.With().Str("cluster", utils.GetClusterID(cluster.ClusterInfo)).Logger()
	return &ClaudieCloudProvider{
		projectName:       projectName,
		configCluster:     cluster,
		resourceName:      rName,
		resourceNamespace: rNamespace,
		nodesCache:        getNodesCache(cluster.ClusterInfo.NodePools),
		nodeManager:       nm,
	}
}

// getClaudieState returns a *pb.K8Scluster, resourceName and resourceNamespace from Claudie, for this particular ClaudieCloudProvider instance.
func getClaudieState(projectName, clusterName string) (*pb.K8Scluster, string, string, error) {
	var cc *grpc.ClientConn
	var err error
	var res *pb.GetConfigFromDBResponse
	cboxURL := strings.ReplaceAll(envs.ContextBoxURL, ":tcp://", "")

	if cc, err = utils.GrpcDialWithRetryAndBackoff("context-box", cboxURL); err != nil {
		return nil, "", "", fmt.Errorf("failed to dial context-box at %s : %w", cboxURL, err)
	}
	defer func() {
		if err := cc.Close(); err != nil {
			log.Err(err).Msgf("Failed to close context-box connection")
		}
	}()

	c := pb.NewContextBoxServiceClient(cc)
	if res, err = c.GetConfigFromDB(context.Background(), &pb.GetConfigFromDBRequest{Id: projectName, Type: pb.IdType_NAME}); err != nil {
		return nil, "", "", fmt.Errorf("failed to get config for project %s : %w", projectName, err)
	}

	for _, cluster := range res.Config.DesiredState.Clusters {
		if cluster.ClusterInfo.Name == clusterName {
			return cluster, res.Config.ResourceName, res.Config.ResourceNamespace, nil
		}
	}
	return nil, "", "", fmt.Errorf("failed to find cluster %s in config for a project %s", clusterName, projectName)
}

// getNodesCache returns a map of nodeCache, regarding all information needed based on the nodepools with autoscaling enabled.
func getNodesCache(nodepools []*pb.NodePool) map[string]*nodeCache {
	var nc = make(map[string]*nodeCache, len(nodepools))
	for _, np := range nodepools {
		if np.GetDynamicNodePool() != nil {
			// Cache nodepools, which are autoscaled.
			if np.GetDynamicNodePool().AutoscalerConfig != nil {
				// Create nodeGroup struct.
				ng := &protos.NodeGroup{
					Id:      np.Name,
					MinSize: np.GetDynamicNodePool().AutoscalerConfig.Min,
					MaxSize: np.GetDynamicNodePool().AutoscalerConfig.Max,
					Debug:   fmt.Sprintf("Nodepool %s [min %d, max %d]", np.Name, np.GetDynamicNodePool().AutoscalerConfig.Min, np.GetDynamicNodePool().AutoscalerConfig.Max),
				}
				// Append ng to the final slice.
				nc[np.Name] = &nodeCache{nodeGroup: ng, nodepool: np, targetSize: np.GetDynamicNodePool().Count}
			}
		}
	}
	return nc
}

// NodeGroups returns all node groups configured for this cloud provider.
func (c *ClaudieCloudProvider) NodeGroups(_ context.Context, req *protos.NodeGroupsRequest) (*protos.NodeGroupsResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got NodeGroups request")
	ngs := make([]*protos.NodeGroup, 0, len(c.nodesCache))
	for _, ngc := range c.nodesCache {
		ngs = append(ngs, ngc.nodeGroup)
	}
	return &protos.NodeGroupsResponse{NodeGroups: ngs}, nil
}

// NodeGroupForNode returns the node group for the given node.
// The node group id is an empty string if the node should not
// be processed by cluster autoscaler.
func (c *ClaudieCloudProvider) NodeGroupForNode(_ context.Context, req *protos.NodeGroupForNodeRequest) (*protos.NodeGroupForNodeResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got NodeGroupForNode request")
	nodeName := req.Node.Name
	// Initialise as empty response.
	nodeGroup := &protos.NodeGroup{}
	// Try to find if node is from any NodeGroup
	for id, ngc := range c.nodesCache {
		// If node name contains ng.Id (nodepool name), return this NodeGroup.
		if strings.Contains(nodeName, id) {
			nodeGroup = ngc.nodeGroup
			break
		}
	}
	return &protos.NodeGroupForNodeResponse{NodeGroup: nodeGroup}, nil
}

// PricingNodePrice returns a theoretical minimum price of running a node for
// a given period of time on a perfectly matching machine.
// Implementation optional.
func (c *ClaudieCloudProvider) PricingNodePrice(_ context.Context, req *protos.PricingNodePriceRequest) (*protos.PricingNodePriceResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got PricingNodePrice request; Not implemented")
	return nil, ErrNotImplemented
}

// PricingPodPrice returns a theoretical minimum price of running a pod for a given
// period of time on a perfectly matching machine.
// Implementation optional.
func (c *ClaudieCloudProvider) PricingPodPrice(_ context.Context, req *protos.PricingPodPriceRequest) (*protos.PricingPodPriceResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got PricingPodPrice request; Not implemented")
	return nil, ErrNotImplemented
}

// GPULabel returns the label added to nodes with GPU resource.
func (c *ClaudieCloudProvider) GPULabel(_ context.Context, req *protos.GPULabelRequest) (*protos.GPULabelResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got GPULabel request")
	return &protos.GPULabelResponse{Label: GpuLabel}, nil
}

// GetAvailableGPUTypes return all available GPU types cloud provider supports.
func (c *ClaudieCloudProvider) GetAvailableGPUTypes(_ context.Context, req *protos.GetAvailableGPUTypesRequest) (*protos.GetAvailableGPUTypesResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got GetAvailableGPUTypes request")
	return &protos.GetAvailableGPUTypesResponse{}, nil
}

// Cleanup cleans up open resources before the cloud provider is destroyed, i.e. go routines etc.
func (c *ClaudieCloudProvider) Cleanup(_ context.Context, req *protos.CleanupRequest) (*protos.CleanupResponse, error) {
	log.Info().Msgf("Got Cleanup request")
	return &protos.CleanupResponse{}, nil
}

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
func (c *ClaudieCloudProvider) Refresh(_ context.Context, req *protos.RefreshRequest) (*protos.RefreshResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Info().Msgf("Got Refresh request")
	return &protos.RefreshResponse{}, c.refresh()
}

// refresh refreshes the state of the claudie provider based of the state from Claudie.
func (c *ClaudieCloudProvider) refresh() error {
	log.Info().Msgf("Refreshing the state")
	if cluster, rName, rNamespace, err := getClaudieState(c.projectName, c.configCluster.ClusterInfo.Name); err != nil {
		log.Err(err).Msgf("Error while refreshing a state of the cluster")
		return fmt.Errorf("error while refreshing a state for the cluster %s : %w", c.configCluster.ClusterInfo.Name, err)
	} else {
		c.configCluster = cluster
		c.resourceName = rName
		c.resourceNamespace = rNamespace
		c.nodesCache = getNodesCache(cluster.ClusterInfo.NodePools)
		if err := c.nodeManager.Refresh(cluster.ClusterInfo.NodePools); err != nil {
			return fmt.Errorf("failed to refresh node manager : %w", err)
		}
	}
	return nil
}

// SendAutoscalerEvent will sent the resourceName and resourceNamespace to the InputManifest controller,
// when a scaleup or scaledown occurs
func (c *ClaudieCloudProvider) sendAutoscalerEvent() error {
	var cc *grpc.ClientConn
	var err error
	frontendURL := strings.ReplaceAll(envs.FrontendURL, ":tcp://", "")
	log.Info().Msgf("Sending autoscale event to %s: %s, %s, ", frontendURL, c.resourceName, c.resourceNamespace)
	if cc, err = utils.GrpcDialWithRetryAndBackoff("frontend", frontendURL); err != nil {
		return fmt.Errorf("failed to dial frontend at %s : %w", envs.FrontendURL, err)
	}
	client := pb.NewFrontendServiceClient(cc)
	if err := frontend.SendAutoscalerEvent(client, &pb.SendAutoscalerEventRequest{
		InputManifestName:      c.resourceName,
		InputManifestNamespace: c.resourceNamespace,
	}); err != nil {
		return fmt.Errorf("error while sending autoscaling event to Frontend : %w", err)
	}
	return nil
}
