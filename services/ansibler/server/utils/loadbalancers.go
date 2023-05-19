package utils

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/berops/claudie/internal/templateUtils"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
)

/*
Loadbalancers are set up via ansible playbooks

The layout of the files/directories for a single k8s cluster loadbalancers is:

clusters/
└── k8s-cluster-1/
	│
	├── lb-cluster-1/
	│	├── key.pem
	│	├── lb.conf
	│	└── nginx.yml
	│
	├── lb-cluster-2/
	│	├── key.pem
	│	├── lb.conf
	│	└── nginx.yml
	├── k8s.pem
	└── inventory.ini
*/

const (
	lbInventoryFile         = "lb-inventory.goini"
	confFile                = "conf.gotpl"
	nginxPlaybookTpl        = "nginx.goyml"
	nginxPlaybook           = "nginx.yml"
	nodeExporterPlaybookTpl = "node-exporter.goyml"
	nodeExporterPlaybook    = "node-exporter.yml"
	nodeExporterService     = "node-exporter.service.j2"
	apiChangePlaybook       = "../../ansible-playbooks/apiEndpointChange.yml"
	loggerPrefix            = "LB-cluster"
)

type APIEndpointChangeState string

const (
	// NoChange represents the 1st. case, no change is needed as LB is currently
	// attached and the desired spec contains no changes.
	NoChange APIEndpointChangeState = "no-change"

	// AttachingLoadBalancer represents 2nd. case, the cluster previously
	// didn't have a LB and the ports needed to communicate with the API server
	// were exposed. After attaching a LB to the existing cluster the ports
	// were closed and are no longer accessible, and thus we need to change the API endpoint.
	AttachingLoadBalancer APIEndpointChangeState = "attaching-load-balancer"

	// DetachingLoadBalancer represents 3rd. case, the cluster had an existing
	// LB attached but the new state removed the LB and thus the API endpoint
	// needs to be changed back to one of the control nodes of the cluster.
	DetachingLoadBalancer APIEndpointChangeState = "detaching-load-balancer"

	// EndpointRenamed represents the 4th. case, the cluster had an existing
	// LB attached and also keeps it but the endpoint has changed.
	EndpointRenamed APIEndpointChangeState = "endpoint-renamed"

	// RoleChangedToAPIServer represents the 5th case, the cluster had an existing
	// LB attached that didn't have a ApiServer role attach but the desired state does.
	RoleChangedToAPIServer APIEndpointChangeState = "role-changed-to-api-server"

	// RoleChangedFromAPIServer represents the 6th case, the cluster had an existing
	// LB attached that had an ApiServer role attached but the desired state doesn't.
	RoleChangedFromAPIServer APIEndpointChangeState = "role-changed-from-api-server"
)

type (
	// LBClustersInfo wraps all Load-balancers and Nodepools used for a single k8s cluster.
	LBClustersInfo struct {
		// LbClusters are Load-Balancers that share the targeted k8s cluster.
		LbClusters []*LBClusterData
		// TargetK8sNodepool are all nodepools used by the targeted k8s cluster.
		TargetK8sNodepool []*pb.NodePool
		// TargetK8sNodepoolKey is the key used for the nodepools.
		TargetK8sNodepoolKey string
		// PreviousAPIEndpointLB holds the endpoint of the previous Load-Balancer endpoint
		// if there was any to be able to handle the endpoint change.
		PreviousAPIEndpointLB string
		// ClusterID contains the ClusterName-Hash- prefix of the kubernetes cluster
		ClusterID string
		// Indicates whether the manifest has no current state i.e. it's the first time it's being build.
		FirstRun bool
	}

	LBClusterData struct {
		// CurrentLbCluster is the current spec.
		// A value of nil means that the current k8s cluster
		// didn't have a LB attached to it.
		CurrentLbCluster *pb.LBcluster

		// DesiredLbCluster is the desired spec.
		// A value of nil means that the targeted k8s cluster
		// will no longer use a LoadBalancer.
		DesiredLbCluster *pb.LBcluster
	}

	LbPlaybookData struct {
		Loadbalancer string
	}

	ConfData struct {
		Roles []LBConfiguration
	}

	LBConfiguration struct {
		Role        *pb.Role
		TargetNodes []*pb.Node
	}
)

// APIEndpointState determines if the API endpoint should be updated with a new
// address, as otherwise communication with the cluster wouldn't be possible.
func (lb *LBClusterData) APIEndpointState() APIEndpointChangeState {
	if lb.CurrentLbCluster == nil && lb.DesiredLbCluster == nil {
		return NoChange
	}

	if lb.CurrentLbCluster == nil && lb.DesiredLbCluster != nil {
		return AttachingLoadBalancer
	}

	if lb.CurrentLbCluster != nil && lb.DesiredLbCluster == nil {
		return DetachingLoadBalancer
	}

	if lb.CurrentLbCluster.Dns.Endpoint != lb.DesiredLbCluster.Dns.Endpoint {
		return EndpointRenamed
	}

	// check if role changed.
	isAPIServer := utils.HasAPIServerRole(lb.CurrentLbCluster.Roles)
	if utils.HasAPIServerRole(lb.DesiredLbCluster.Roles) && !isAPIServer {
		return RoleChangedToAPIServer
	}

	if isAPIServer && !utils.HasAPIServerRole(lb.DesiredLbCluster.Roles) {
		return RoleChangedFromAPIServer
	}

	return NoChange
}

func handleAPIEndpointChange(apiServer *LBClusterData, k8sCluster *LBClustersInfo, k8sDirectory string) error {
	if apiServer == nil {
		// if there is no ApiSever LB that means that the ports 6443 are exposed
		// on the nodes, and thus we don't need to anything.
		return nil
	}

	var oldEndpoint string
	var newEndpoint string

	switch apiServer.APIEndpointState() {
	case NoChange:
		return nil
	case EndpointRenamed:
		oldEndpoint = apiServer.CurrentLbCluster.Dns.Endpoint
		newEndpoint = apiServer.DesiredLbCluster.Dns.Endpoint
	case RoleChangedFromAPIServer:
		oldEndpoint = apiServer.CurrentLbCluster.Dns.Endpoint

		// 1st find if any control node was an API server.
		if node, err := utils.FindAPIEndpointNode(k8sCluster.TargetK8sNodepool); err == nil {
			newEndpoint = node.Public
			break
		}

		// 2nd choose one of the control nodes as the api endpoint.
		node, err := utils.FindControlNode(k8sCluster.TargetK8sNodepool)
		if err != nil {
			return err
		}

		node.NodeType = pb.NodeType_apiEndpoint

		newEndpoint = node.Public
	case RoleChangedToAPIServer:
		newEndpoint = apiServer.DesiredLbCluster.Dns.Endpoint

		// 1st check if there was any APISERVER-LB previously attached to the k8scluster.
		if k8sCluster.PreviousAPIEndpointLB != "" {
			oldEndpoint = k8sCluster.PreviousAPIEndpointLB
			break
		}

		// 2nd check if any other LB was previously an ApiServer.
		if oldAPIServer := findCurrentAPILoadBalancer(k8sCluster.LbClusters); oldAPIServer != nil {
			oldEndpoint = oldAPIServer.CurrentLbCluster.Dns.Endpoint
			break
		}

		// 3rd pick the control node as the previous ApiServer.
		node, err := utils.FindAPIEndpointNode(k8sCluster.TargetK8sNodepool)
		if err != nil {
			return fmt.Errorf("failed to find ApiEndpoint k8s node, couldn't update Api server endpoint")
		}

		oldEndpoint = node.Public
	case AttachingLoadBalancer:
		newEndpoint = apiServer.DesiredLbCluster.Dns.Endpoint

		if k8sCluster.FirstRun {
			// it's the first time the manifest goes through the workflow,
			// thus we don't need to change the api endpoint.
			return nil
		}

		// We know that it's not a first run, so before we use the node as the old APIServer
		// endpoint we check a few other possibilities.

		// 1st. check if there was any APIServer-LB previously attached to the k8scluster
		if k8sCluster.PreviousAPIEndpointLB != "" {
			oldEndpoint = k8sCluster.PreviousAPIEndpointLB
			break
		}

		// 2nd check if any other LB was previously an APIServer.
		if oldAPIServer := findCurrentAPILoadBalancer(k8sCluster.LbClusters); oldAPIServer != nil {
			oldEndpoint = oldAPIServer.CurrentLbCluster.Dns.Endpoint
			break
		}

		// 3rd pick the control node as the previous ApiServer.
		node, err := utils.FindAPIEndpointNode(k8sCluster.TargetK8sNodepool)
		if err != nil {
			return fmt.Errorf("failed to find APIEndpoint k8s node, couldn't update Api server endpoint")
		}

		node.NodeType = pb.NodeType_master // remove the Endpoint type from the node.
		oldEndpoint = node.Public
	case DetachingLoadBalancer:
		oldEndpoint = apiServer.CurrentLbCluster.Dns.Endpoint

		// 1st find if any control node was an API server.
		if node, err := utils.FindAPIEndpointNode(k8sCluster.TargetK8sNodepool); err == nil {
			newEndpoint = node.Public
			break
		}

		// 2nd choose one of the control nodes as the api endpoint.
		node, err := utils.FindControlNode(k8sCluster.TargetK8sNodepool)
		if err != nil {
			return err
		}

		node.NodeType = pb.NodeType_apiEndpoint

		newEndpoint = node.Public
	}

	lbCluster := apiServer.DesiredLbCluster
	if lbCluster == nil {
		lbCluster = apiServer.CurrentLbCluster
	}
	log.Debug().Str(loggerPrefix, utils.GetClusterID(lbCluster.ClusterInfo)).Msgf("Changing the API endpoint from %s to %s", oldEndpoint, newEndpoint)

	if err := changeAPIEndpoint(lbCluster.ClusterInfo.Name, oldEndpoint, newEndpoint, k8sDirectory); err != nil {
		return fmt.Errorf("error while changing the endpoint for %s : %w", lbCluster.ClusterInfo.Name, err)
	}

	return nil
}

// findCurrentAPILoadBalancers finds the current Load-Balancer for the API server
func findCurrentAPILoadBalancer(lbs []*LBClusterData) *LBClusterData {
	for _, lb := range lbs {
		if lb.CurrentLbCluster != nil {
			if utils.HasAPIServerRole(lb.CurrentLbCluster.Roles) {
				return lb
			}
		}
	}

	return nil
}

// setUpNginx sets up the nginx loadbalancer based on the input manifest specification
// return error if not successful, nil otherwise
func setUpNginx(lb *pb.LBcluster, targetedNodepool []*pb.NodePool, directory string) error {
	//prepare data for .conf
	templateLoader := templateUtils.TemplateLoader{Directory: templateUtils.AnsiblerTemplates}
	template := templateUtils.Templates{Directory: directory}
	tpl, err := templateLoader.LoadTemplate(confFile)
	if err != nil {
		return fmt.Errorf("error while loading %s template for %w", confFile, err)
	}
	//get control and compute nodes
	controlTarget, computeTarget := splitNodesByType(targetedNodepool)
	var lbRoles []LBConfiguration
	for _, role := range lb.Roles {
		target := assignTarget(controlTarget, computeTarget, role.Target)
		if target == nil {
			return fmt.Errorf("target %v did not specify any nodes", role.Target)
		}
		lbRoles = append(lbRoles, LBConfiguration{Role: role, TargetNodes: target})
	}
	//create .conf file
	err = template.Generate(tpl, "lb.conf", ConfData{Roles: lbRoles})
	if err != nil {
		return fmt.Errorf("error while generating lb.conf for %s : %w", lb.ClusterInfo.Name, err)
	}
	tpl, err = templateLoader.LoadTemplate(nginxPlaybookTpl)
	if err != nil {
		return fmt.Errorf("error while loading %s for %s : %w", nginxPlaybookTpl, lb.ClusterInfo.Name, err)
	}
	err = template.Generate(tpl, nginxPlaybook, LbPlaybookData{Loadbalancer: lb.ClusterInfo.Name})
	if err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nginxPlaybook, lb.ClusterInfo.Name, err)
	}
	//run the playbook
	ansible := Ansible{Playbook: nginxPlaybook, Inventory: filepath.Join("..", inventoryFile), Directory: directory}
	err = ansible.RunAnsiblePlaybook(fmt.Sprintf("LB - %s-%s", lb.ClusterInfo.Name, lb.ClusterInfo.Hash))
	if err != nil {
		return fmt.Errorf("error while running ansible for %s : %w", lb.ClusterInfo.Name, err)
	}
	return nil
}

// setUpNodeExporter sets up node-exporter on the LB node
// return error if not successful, nil otherwise
func setUpNodeExporter(lb *pb.LBcluster, directory string) error {
	// generate node-exporter playbook template
	templateLoader := templateUtils.TemplateLoader{Directory: templateUtils.AnsiblerTemplates}
	template := templateUtils.Templates{Directory: directory}
	tpl, err := templateLoader.LoadTemplate(nodeExporterPlaybookTpl)
	if err != nil {
		return fmt.Errorf("error while loading %s template for %s : %w", nodeExporterPlaybookTpl, lb.ClusterInfo.Name, err)
	}
	if err = template.Generate(tpl, nodeExporterPlaybook, LbPlaybookData{Loadbalancer: lb.ClusterInfo.Name}); err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nodeExporterPlaybook, lb.ClusterInfo.Name, err)
	}

	// create node-exporter.service.j2 for the node-exporter playbook
	tpl, err = templateLoader.LoadTemplate(nodeExporterService)
	if err != nil {
		return fmt.Errorf("error while loading %s template for %s : %w", nodeExporterService, lb.ClusterInfo.Name, err)
	}
	if err = template.Generate(tpl, nodeExporterService, LbPlaybookData{Loadbalancer: lb.ClusterInfo.Name}); err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nodeExporterService, lb.ClusterInfo.Name, err)
	}

	//run the playbook
	ansible := Ansible{Playbook: nodeExporterPlaybook, Inventory: filepath.Join("..", inventoryFile), Directory: directory}
	if err = ansible.RunAnsiblePlaybook(fmt.Sprintf("LB - %s-%s", lb.ClusterInfo.Name, lb.ClusterInfo.Hash)); err != nil {
		return fmt.Errorf("error while running ansible for %s : %w", lb.ClusterInfo.Name, err)
	}

	return nil
}

// splitNodesByType returns two slices of *pb.Node, one for control nodes and one for compute
func splitNodesByType(nodepools []*pb.NodePool) (controlNodes, computeNodes []*pb.Node) {
	for _, nodepools := range nodepools {
		for _, node := range nodepools.Nodes {
			if node.NodeType == pb.NodeType_apiEndpoint || node.NodeType == pb.NodeType_master {
				controlNodes = append(controlNodes, node)
			} else {
				computeNodes = append(computeNodes, node)
			}
		}
	}
	return controlNodes, computeNodes
}

// generateK8sBaseFiles generates the base loadbalancer files, like inventory, keys, etc.
// return error if not successful, nil otherwise
func generateK8sBaseFiles(k8sDirectory string, lbInfo *LBClustersInfo) error {
	if err := utils.CreateDirectory(k8sDirectory); err != nil {
		return fmt.Errorf("failed to create directory %s : %w", k8sDirectory, err)
	}

	if err := utils.CreateKeyFile(lbInfo.TargetK8sNodepoolKey, k8sDirectory, "k8s.pem"); err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	var lbSlice []*pb.LBcluster
	for _, lb := range lbInfo.LbClusters {
		if lb.DesiredLbCluster != nil {
			lbSlice = append(lbSlice, lb.DesiredLbCluster)
		}
	}
	//generate inventory
	err := generateInventoryFile(lbInventoryFile, k8sDirectory, LbInventoryData{
		K8sNodepools: lbInfo.TargetK8sNodepool,
		LBClusters:   lbSlice,
		ClusterID:    lbInfo.ClusterID,
	})
	if err != nil {
		return fmt.Errorf("error while generating inventory file for %s : %w", k8sDirectory, err)
	}
	return nil
}

// assignTarget returns a target nodes for pb.Target
// if no target matches the pb.Target enum, returns nil
func assignTarget(controlTarget, computeTarget []*pb.Node, target pb.Target) (targetNodes []*pb.Node) {
	if target == pb.Target_k8sAllNodes {
		targetNodes = append(controlTarget, computeTarget...)
	} else if target == pb.Target_k8sControlPlane {
		targetNodes = controlTarget
	} else if target == pb.Target_k8sComputePlane {
		targetNodes = computeTarget
	}
	return targetNodes
}
