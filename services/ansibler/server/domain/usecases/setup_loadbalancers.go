package usecases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/berops/claudie/internal/templateUtils"
	commonUtils "github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/berops/claudie/services/ansibler/server/utils"
	"github.com/berops/claudie/services/ansibler/templates"
	"github.com/rs/zerolog"
)

const (
	// nodeExporterPlaybookFileName defines name for node exporter playbook.
	nodeExporterPlaybookFileName = "node-exporter.yml"
	// nginxPlaybookName defines name for nginx playbook.
	nginxPlaybookName = "nginx.yml"
	// nginxConfigName defines name for nginx config.
	nginxConfigName = "lb.conf"
)

func (u *Usecases) SetUpLoadbalancers(request *pb.SetUpLBRequest) (*pb.SetUpLBResponse, error) {
	logger := commonUtils.CreateLoggerWithProjectAndClusterName(request.ProjectName, commonUtils.GetClusterID(request.Desired.ClusterInfo))
	logger.Info().Msgf("Setting up the loadbalancers")

	currentLBClusters := make(map[string]*spec.LBcluster)
	for _, lbCluster := range request.CurrentLbs {
		currentLBClusters[lbCluster.ClusterInfo.Name] = lbCluster
	}

	lbClustersInfo := &utils.LBClustersInfo{
		FirstRun:              request.FirstRun,
		TargetK8sNodepool:     request.Desired.ClusterInfo.NodePools,
		PreviousAPIEndpointLB: request.PreviousAPIEndpoint,
		ClusterID:             commonUtils.GetClusterID(request.Desired.ClusterInfo),
	}
	for _, lbCluster := range request.DesiredLbs {
		lbClustersInfo.LbClusters = append(lbClustersInfo.LbClusters, &utils.LBClusterData{
			DesiredLbCluster: lbCluster,
			// if there is a value in the map it will return it, otherwise nil is returned.
			CurrentLbCluster: currentLBClusters[lbCluster.ClusterInfo.Name],
		})
	}

	if err := setUpLoadbalancers(request.Desired.ClusterInfo, lbClustersInfo, logger, u.SpawnProcessLimit); err != nil {
		logger.Err(err).Msgf("Error encountered while setting up the loadbalancers")
		return nil, fmt.Errorf("error encountered while setting up the loadbalancers for cluster %s project %s : %w", request.Desired.ClusterInfo.Name, request.ProjectName, err)
	}

	logger.Info().Msgf("Loadbalancers were successfully set up")
	return &pb.SetUpLBResponse{Desired: request.Desired, CurrentLbs: request.CurrentLbs, DesiredLbs: request.DesiredLbs}, nil
}

// setUpLoadbalancers sets up the loadbalancers along with DNS and verifies their configuration
func setUpLoadbalancers(desiredK8sClusterInfo *spec.ClusterInfo, lbClustersInfo *utils.LBClustersInfo, logger zerolog.Logger, spawnProcessLimit chan struct{}) error {
	clusterName := desiredK8sClusterInfo.Name
	clusterBaseDirectory := filepath.Join(baseDirectory, outputDirectory, fmt.Sprintf("%s-%s-lbs", clusterName, commonUtils.CreateHash(commonUtils.HashLength)))

	if err := utils.GenerateLBBaseFiles(clusterBaseDirectory, lbClustersInfo); err != nil {
		return fmt.Errorf("error encountered while generating base files for %s : %w", clusterName, err)
	}

	err := commonUtils.ConcurrentExec(lbClustersInfo.LbClusters,
		func(_ int, lbCluster *utils.LBClusterData) error {
			var (
				loggerPrefix = "LB-cluster"
				lbClusterId  = commonUtils.GetClusterID(lbCluster.DesiredLbCluster.ClusterInfo)
			)

			logger.Info().Str(loggerPrefix, lbClusterId).Msg("Setting up the loadbalancer cluster")

			// Create the directory where files will be generated
			clusterDirectory := filepath.Join(clusterBaseDirectory, lbClusterId)
			if err := commonUtils.CreateDirectory(clusterDirectory); err != nil {
				return fmt.Errorf("failed to create directory %s : %w", clusterDirectory, err)
			}

			if err := commonUtils.CreateKeysForDynamicNodePools(commonUtils.GetCommonDynamicNodePools(lbCluster.DesiredLbCluster.ClusterInfo.NodePools), clusterDirectory); err != nil {
				return fmt.Errorf("failed to create key file(s) for dynamic nodepools : %w", err)
			}

			if err := commonUtils.CreateKeysForStaticNodepools(commonUtils.GetCommonStaticNodePools(lbCluster.DesiredLbCluster.ClusterInfo.NodePools), clusterDirectory); err != nil {
				return fmt.Errorf("failed to create key file(s) for static nodes : %w", err)
			}

			if err := setUpNodeExporter(lbCluster.DesiredLbCluster, clusterDirectory, spawnProcessLimit); err != nil {
				return err
			}

			if err := setUpNginx(lbCluster.DesiredLbCluster, lbClustersInfo.TargetK8sNodepool, clusterDirectory, spawnProcessLimit); err != nil {
				return err
			}

			logger.Info().Str(loggerPrefix, lbClusterId).Msg("Loadbalancer cluster successfully set up")
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("error while setting up the loadbalancers for cluster %s : %w", clusterName, err)
	}

	var desiredApiServerTypeLBCluster *utils.LBClusterData
	for _, lbClusterInfo := range lbClustersInfo.LbClusters {
		if commonUtils.HasAPIServerRole(lbClusterInfo.DesiredLbCluster.Roles) {
			desiredApiServerTypeLBCluster = lbClusterInfo
		}
	}

	// If we didn't find any Api server type LB cluster in the desired state,
	// it's possible that we've changed the role from an API server to some other role.
	// This won't be caught by the above check.
	// So we have to do an additional check for the desiredApiServerTypeLBCluster using the current state.
	if desiredApiServerTypeLBCluster == nil {
		desiredApiServerTypeLBCluster = utils.FindCurrentAPIServerTypeLBCluster(lbClustersInfo.LbClusters)
	}

	if err := utils.HandleAPIEndpointChange(desiredApiServerTypeLBCluster, desiredK8sClusterInfo, lbClustersInfo, clusterBaseDirectory, spawnProcessLimit); err != nil {
		return fmt.Errorf("failed to find a candidate for the Api Server: %w", err)
	}

	return os.RemoveAll(clusterBaseDirectory)
}

// setUpNodeExporter sets up node-exporter on each node of the LB cluster.
// Returns error if not successful, nil otherwise.
func setUpNodeExporter(lbCluster *spec.LBcluster, clusterDirectory string, spawnProcessLimit chan struct{}) error {
	var playbookParameters = utils.LBPlaybookParameters{Loadbalancer: lbCluster.ClusterInfo.Name}

	// Generate node-exporter Ansible playbook from template
	template, err := templateUtils.LoadTemplate(templates.NodeExporterPlaybookTemplate)
	if err != nil {
		return fmt.Errorf("error while loading %s template for node_exporter playbook : %w", lbCluster.ClusterInfo.Name, err)
	}

	tpl := templateUtils.Templates{Directory: clusterDirectory}
	if err := tpl.Generate(template, nodeExporterPlaybookFileName, playbookParameters); err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nodeExporterPlaybookFileName, lbCluster.ClusterInfo.Name, err)
	}

	// Run the Ansible playbook
	ansible := utils.Ansible{
		Directory:         clusterDirectory,
		Playbook:          nodeExporterPlaybookFileName,
		Inventory:         filepath.Join("..", utils.InventoryFileName),
		SpawnProcessLimit: spawnProcessLimit,
	}

	if err = ansible.RunAnsiblePlaybook(fmt.Sprintf("LB - %s-%s", lbCluster.ClusterInfo.Name, lbCluster.ClusterInfo.Hash)); err != nil {
		return fmt.Errorf("error while running ansible for %s : %w", lbCluster.ClusterInfo.Name, err)
	}

	return nil
}

// setUpNginx sets up the nginx loadbalancer based on the input manifest specification.
// Return error if not successful, nil otherwise
func setUpNginx(lbCluster *spec.LBcluster, targetK8sNodepool []*spec.NodePool, clusterDirectory string, spawnProcessLimit chan struct{}) error {
	lbClusterRolesInfo := targetPools(lbCluster, targetK8sNodepool)
	// Generate the nginx config file
	nginxConfTemplate, err := templateUtils.LoadTemplate(templates.NginxConfigTemplate)
	tpl := templateUtils.Templates{Directory: clusterDirectory}
	if err != nil {
		return fmt.Errorf("error while loading nginx config template : %w", err)
	}
	nginxPlaybookTemplate, err := templateUtils.LoadTemplate(templates.NginxPlaybookTemplate)
	if err != nil {
		return fmt.Errorf("error while loading nginx playbook template : %w", err)
	}

	if err := tpl.Generate(nginxConfTemplate, nginxConfigName, utils.NginxConfigTemplateParameters{Roles: lbClusterRolesInfo}); err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nginxConfigName, lbCluster.ClusterInfo.Name, err)
	}

	if err := tpl.Generate(nginxPlaybookTemplate, nginxPlaybookName, utils.LBPlaybookParameters{Loadbalancer: lbCluster.ClusterInfo.Name}); err != nil {
		return fmt.Errorf("error while generating %s for %s : %w", nginxPlaybookName, lbCluster.ClusterInfo.Name, err)
	}

	ansible := utils.Ansible{
		Playbook:          nginxPlaybookName,
		Inventory:         filepath.Join("..", utils.InventoryFileName),
		Directory:         clusterDirectory,
		SpawnProcessLimit: spawnProcessLimit,
	}

	err = ansible.RunAnsiblePlaybook(fmt.Sprintf("LB - %s-%s", lbCluster.ClusterInfo.Name, lbCluster.ClusterInfo.Hash))
	if err != nil {
		return fmt.Errorf("error while running ansible for %s : %w", lbCluster.ClusterInfo.Name, err)
	}

	return nil
}

func targetPools(lbCluster *spec.LBcluster, targetK8sNodepool []*spec.NodePool) []utils.LBClusterRolesInfo {
	//TODO: remove in favor of targetNodepools
	targetControlNodes, targetComputeNodes := splitNodesByType(targetK8sNodepool)

	var lbClusterRolesInfo []utils.LBClusterRolesInfo
	for _, role := range lbCluster.Roles {
		if len(role.TargetPools) == 0 {
			// TODO: remove in favor of targetNodepools
			target := assignTarget(targetControlNodes, targetComputeNodes, role.Target)
			lbClusterRolesInfo = append(lbClusterRolesInfo, utils.LBClusterRolesInfo{Role: role, TargetNodes: target})
		} else {
			lbClusterRolesInfo = append(lbClusterRolesInfo, utils.LBClusterRolesInfo{
				Role:        role,
				TargetNodes: targetNodes(role.TargetPools, targetK8sNodepool),
			})
		}
	}

	return lbClusterRolesInfo
}

func targetNodes(targetPools []string, targetk8sPools []*spec.NodePool) (nodes []*spec.Node) {
	var pools []*spec.NodePool

	for _, target := range targetPools {
		for _, np := range targetk8sPools {
			if np.GetDynamicNodePool() != nil {
				if name, _ := commonUtils.MatchNameAndHashWithTemplate(target, np.Name); name != "" {
					pools = append(pools, np)
				}
			} else if np.GetStaticNodePool() != nil {
				if target == np.Name {
					pools = append(pools, np)
				}
			}
		}
	}

	for _, np := range pools {
		nodes = append(nodes, np.Nodes...)
	}

	return
}

// splitNodesByType returns two slices of *pb.Node, one for control nodes and one for compute nodes.
func splitNodesByType(nodepools []*spec.NodePool) (controlNodes, computeNodes []*spec.Node) {
	for _, nodepool := range nodepools {
		for _, node := range nodepool.Nodes {
			if node.NodeType == spec.NodeType_apiEndpoint || node.NodeType == spec.NodeType_master {
				controlNodes = append(controlNodes, node)
			} else {
				computeNodes = append(computeNodes, node)
			}
		}
	}

	return controlNodes, computeNodes
}

// assignTarget returns a target nodes for pb.Target.
// If no target matches the pb.Target enum, returns nil
func assignTarget(targetControlNodes, targetComputeNodes []*spec.Node, target spec.Target) (targetNodes []*spec.Node) {
	if target == spec.Target_k8sAllNodes {
		targetNodes = append(targetControlNodes, targetComputeNodes...)
	} else if target == spec.Target_k8sControlPlane {
		targetNodes = targetControlNodes
	} else if target == spec.Target_k8sComputePlane {
		targetNodes = targetComputeNodes
	}

	return targetNodes
}
