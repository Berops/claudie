package cluster_builder

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"

	comm "github.com/berops/claudie/internal/command"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/templates"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/terraform"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	TemplatesRootDir     = "services/terraformer/templates"
	Output               = "services/terraformer/server/clusters"
	baseSubnetCIDR       = "10.0.0.0/24"
	defaultOctetToChange = 2
)

type K8sInfo struct{ LoadBalancers []*spec.LBcluster }
type LBInfo struct{ Roles []*spec.Role }

// ClusterBuilder wraps data needed for building a cluster.
type ClusterBuilder struct {
	// DesiredClusterInfo contains the information about the
	// desired state of the cluster.
	DesiredClusterInfo *spec.ClusterInfo
	// CurrentClusterInfo contains the information about the
	// current state of the cluster.
	CurrentClusterInfo *spec.ClusterInfo
	// ProjectName is the name of the manifest.
	ProjectName string
	// ClusterType is the type of the cluster being build
	// LoadBalancer or K8s.
	ClusterType spec.ClusterType
	// K8sInfo contains additional data for when building kubernetes clusters.
	K8sInfo K8sInfo
	// LBInfo contains additional data for when building loadbalancer clusters.
	LBInfo LBInfo
	// SpawnProcessLimit represents a synchronization channel which limits the number of spawned terraform
	// processes. This values should always be non-nil and be buffered, where the capacity indicates
	// the limit.
	SpawnProcessLimit chan struct{}
}

// CreateNodepools creates node pools for the cluster.
func (c ClusterBuilder) CreateNodepools() error {
	clusterID := utils.GetClusterID(c.DesiredClusterInfo)
	clusterDir := filepath.Join(Output, clusterID)

	// Calculate CIDR, so they do not change if nodepool order changes
	// https://github.com/berops/claudie/issues/647
	// Order them by provider and region
	for _, nps := range utils.GroupNodepoolsByProviderRegion(c.DesiredClusterInfo) {
		if err := calculateCIDR(baseSubnetCIDR, utils.GetDynamicNodePools(nps)); err != nil {
			return fmt.Errorf("error while generating CIDR for nodepools : %w", err)
		}
	}

	if err := c.generateFiles(clusterID, clusterDir); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	terraform := terraform.Terraform{
		Directory:         clusterDir,
		SpawnProcessLimit: c.SpawnProcessLimit,
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		terraform.Stdout = comm.GetStdOut(clusterID)
		terraform.Stderr = comm.GetStdErr(clusterID)
	}

	if err := terraform.Init(); err != nil {
		return fmt.Errorf("error while running terraform init in %s : %w", clusterID, err)
	}

	if err := terraform.Apply(); err != nil {
		return fmt.Errorf("error while running terraform apply in %s : %w", clusterID, err)
	}
	oldNodes := c.getCurrentNodes()

	// fill new nodes with output
	for _, nodepool := range c.DesiredClusterInfo.NodePools {
		if np := nodepool.GetDynamicNodePool(); np != nil {
			k := fmt.Sprintf(
				"%s_%s_%s",
				nodepool.Name,
				np.GetProvider().GetSpecName(),
				templates.Fingerprint(templates.ExtractTargetPath(np.GetProvider().GetTemplates())),
			)
			output, err := terraform.Output(k)
			if err != nil {
				return fmt.Errorf("error while getting output from terraform for %s : %w", nodepool.Name, err)
			}
			out, err := readIPs(output)
			if err != nil {
				return fmt.Errorf("error while reading the terraform output for %s : %w", nodepool.Name, err)
			}
			fillNodes(&out, nodepool, oldNodes)
		}
	}

	// Clean after terraform
	if err := os.RemoveAll(clusterDir); err != nil {
		return fmt.Errorf("error while deleting files in %s : %w", clusterDir, err)
	}

	return nil
}

// DestroyNodepools destroys nodepools for the cluster.
func (c ClusterBuilder) DestroyNodepools() error {
	clusterID := utils.GetClusterID(c.CurrentClusterInfo)
	clusterDir := filepath.Join(Output, clusterID)

	// Calculate CIDR, in case some nodepools do not have it, due to error.
	// https://github.com/berops/claudie/issues/647
	// Order them by provider and region
	for _, nps := range utils.GroupNodepoolsByProviderRegion(c.CurrentClusterInfo) {
		if err := calculateCIDR(baseSubnetCIDR, utils.GetDynamicNodePools(nps)); err != nil {
			return fmt.Errorf("error while generating CIDR for nodepools : %w", err)
		}
	}

	if err := c.generateFiles(clusterID, clusterDir); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	terraform := terraform.Terraform{
		Directory:         clusterDir,
		SpawnProcessLimit: c.SpawnProcessLimit,
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		terraform.Stdout = comm.GetStdOut(clusterID)
		terraform.Stderr = comm.GetStdErr(clusterID)
	}

	if err := terraform.Init(); err != nil {
		return fmt.Errorf("error while running terraform init in %s : %w", clusterID, err)
	}

	if err := terraform.Destroy(); err != nil {
		return fmt.Errorf("error while running terraform apply in %s : %w", clusterID, err)
	}

	if err := os.RemoveAll(clusterDir); err != nil {
		return fmt.Errorf("error while deleting files in %s : %w", clusterDir, err)
	}

	return nil
}

// generateFiles creates all the necessary terraform files used to create/destroy node pools.
func (c *ClusterBuilder) generateFiles(clusterID, clusterDir string) error {
	backend := templates.Backend{
		ProjectName: c.ProjectName,
		ClusterName: clusterID,
		Directory:   clusterDir,
	}

	if err := backend.CreateTFFile(); err != nil {
		return err
	}

	// generate Providers terraform configuration
	usedProviders := templates.UsedProviders{
		ProjectName: c.ProjectName,
		ClusterName: clusterID,
		Directory:   clusterDir,
	}

	if err := usedProviders.CreateUsedProvider(c.CurrentClusterInfo, c.DesiredClusterInfo); err != nil {
		return err
	}

	var clusterInfo *spec.ClusterInfo
	if c.DesiredClusterInfo != nil {
		clusterInfo = c.DesiredClusterInfo
	} else if c.CurrentClusterInfo != nil {
		clusterInfo = c.CurrentClusterInfo
	}

	// Init node slices if needed
	for _, np := range clusterInfo.NodePools {
		if n := np.GetDynamicNodePool(); n != nil {
			nodes := make([]*spec.Node, 0, n.Count)
			nodeNames := make(map[string]struct{}, n.Count)
			// Copy existing nodes into new slice
			for i, node := range np.Nodes {
				if i == int(n.Count) {
					break
				}
				log.Debug().Str("cluster", clusterID).Msgf("Nodepool is reusing node %s", node.Name)
				nodes = append(nodes, node)
				nodeNames[node.Name] = struct{}{}
			}
			// Fill the rest of the nodes with assigned names
			nodepoolID := fmt.Sprintf("%s-%s", clusterID, np.Name)
			for len(nodes) < int(n.Count) {
				// Get a unique name for the new node
				nodeName := getUniqueNodeName(nodepoolID, nodeNames)
				nodeNames[nodeName] = struct{}{}
				nodes = append(nodes, &spec.Node{Name: nodeName})
			}
			np.Nodes = nodes
		}
	}

	clusterData := templates.ClusterData{
		ClusterName: clusterInfo.Name,
		ClusterHash: clusterInfo.Hash,
		ClusterType: c.ClusterType.String(),
	}

	if err := c.generateProviderTemplates(c.CurrentClusterInfo, c.DesiredClusterInfo, clusterID, clusterDir, clusterData); err != nil {
		return fmt.Errorf("error while generating provider templates: %w", err)
	}

	groupedNodepools := utils.GroupNodepoolsByProviderNames(clusterInfo)
	for providerNames, nodepools := range groupedNodepools {
		if providerNames.CloudProviderName == spec.StaticNodepoolInfo_STATIC_PROVIDER.String() {
			continue
		}

		p := nodepools[0].GetDynamicNodePool().GetProvider()
		templatesDownloadDir := filepath.Join(TemplatesRootDir, clusterID, p.SpecName)

		if err := templates.DownloadProvider(templatesDownloadDir, p); err != nil {
			msg := fmt.Sprintf("cluster %q failed to download template repository", clusterID)
			log.Error().Msgf(msg)
			return fmt.Errorf("%s: %w", msg, err)
		}

		nps := make([]templates.NodePoolInfo, 0, len(nodepools))

		for _, np := range nodepools {
			if dnp := np.GetDynamicNodePool(); dnp != nil {
				nps = append(nps, templates.NodePoolInfo{
					Name:      np.Name,
					Nodes:     np.Nodes,
					Details:   np.GetDynamicNodePool(),
					IsControl: np.IsControl,
				})

				if err := utils.CreateKeyFile(dnp.GetPublicKey(), clusterDir, np.GetName()); err != nil {
					return fmt.Errorf("error public key file for %s : %w", clusterDir, err)
				}
			}
		}

		// based on the cluster type fill out the nodepools data to be used
		nodepoolData := templates.Nodepools{
			ClusterData: clusterData,
			NodePools:   nps,
		}

		g := templates.Generator{
			ID:                clusterID,
			TargetDirectory:   clusterDir,
			ReadFromDirectory: templatesDownloadDir,
			TemplatePath:      templates.ExtractTargetPath(p.GetTemplates()),
		}

		if err := g.GenerateNetworking(&templates.Networking{
			ClusterData: clusterData,
			Provider:    nodepools[0].GetDynamicNodePool().GetProvider(),
			Regions:     utils.GetRegions(utils.GetDynamicNodePools(nodepools)),
			K8sData: templates.K8sData{
				HasAPIServer: !slices.Contains(
					utils.ExtractTargetPorts(c.K8sInfo.LoadBalancers),
					6443,
				),
			},
			LBData: templates.LBData{
				Roles: c.LBInfo.Roles,
			},
		}); err != nil {
			return fmt.Errorf("failed to generate networking_common template files: %w", err)
		}

		if err := g.GenerateNodes(&nodepoolData); err != nil {
			return fmt.Errorf("failed to generate nodepool specific templates files: %w", err)
		}

		if err := utils.CreateKeyFile(
			utils.GetAuthCredentials(nps[0].Details.Provider),
			clusterDir,
			providerNames.SpecName,
		); err != nil {
			return fmt.Errorf("error creating provider credential key file for provider %s in %s : %w", providerNames.SpecName, clusterDir, err)
		}
	}

	return nil
}

// getCurrentNodes returns all nodes which are in a current state
func (c *ClusterBuilder) getCurrentNodes() []*spec.Node {
	// group all the nodes together to make searching with respect to IP easy
	var oldNodes []*spec.Node
	if c.CurrentClusterInfo != nil {
		for _, oldNodepool := range c.CurrentClusterInfo.NodePools {
			if oldNodepool.GetDynamicNodePool() != nil {
				oldNodes = append(oldNodes, oldNodepool.Nodes...)
			}
		}
	}
	return oldNodes
}

// calculateCIDR will make sure all nodepools have subnet CIDR calculated.
func calculateCIDR(baseCIDR string, nodepools []*spec.DynamicNodePool) error {
	exists := make(map[string]struct{})
	// Save CIDRs which already exist.
	for _, np := range nodepools {
		exists[np.Cidr] = struct{}{}
	}

	// Calculate new ones if needed.
	for _, np := range nodepools {
		if np.Cidr != "" {
			continue
		}

		cidr, err := getCIDR(baseCIDR, defaultOctetToChange, exists)
		if err != nil {
			return fmt.Errorf("failed to parse CIDR for nodepool : %w", err)
		}

		log.Debug().Msgf("Calculating new VPC subnet CIDR for nodepool. New CIDR [%s]", cidr)
		np.Cidr = cidr
		// Cache calculated CIDR.
		exists[cidr] = struct{}{}
	}

	return nil
}

// fillNodes creates pb.Node slices in desired state, with the new nodes and old nodes
func fillNodes(terraformOutput *templates.NodepoolIPs, newNodePool *spec.NodePool, oldNodes []*spec.Node) {
	// fill slices from terraformOutput maps with names of nodes to ensure an order
	var tempNodes []*spec.Node
	// get sorted list of keys
	utils.IterateInOrder(terraformOutput.IPs, func(nodeName string, IP any) error {
		var nodeType spec.NodeType
		var private string

		if newNodePool.IsControl {
			nodeType = spec.NodeType_master
		} else {
			nodeType = spec.NodeType_worker
		}

		if len(oldNodes) > 0 {
			for _, node := range oldNodes {
				//check if node was defined before
				if fmt.Sprint(IP) == node.Public && nodeName == node.Name {
					// carry privateIP to desired state, so it will not get overwritten in Ansibler
					private = node.Private
					// carry node type since it might be API endpoint, which should not change once set
					nodeType = node.NodeType
					break
				}
			}
		}

		tempNodes = append(tempNodes, &spec.Node{
			Name:     nodeName,
			Public:   fmt.Sprint(terraformOutput.IPs[nodeName]),
			Private:  private,
			NodeType: nodeType,
		})
		return nil
	})

	newNodePool.Nodes = tempNodes
}

// readIPs reads json output format from terraform and unmarshal it into map[string]map[string]string readable by Go.
func readIPs(data string) (templates.NodepoolIPs, error) {
	var result templates.NodepoolIPs
	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(data), &result.IPs)
	return result, err
}

// getUniqueNodeName returns new node name, which is guaranteed to be unique, based on the provided existing names.
func getUniqueNodeName(nodepoolID string, existingNames map[string]struct{}) string {
	index := uint8(1)
	for {
		candidate := fmt.Sprintf("%s-%02x", nodepoolID, index)
		if _, ok := existingNames[candidate]; !ok {
			return candidate
		}
		index++
	}
}

// getCIDR function returns CIDR in IPv4 format, with position replaced by value
// The function does not check if it is a valid CIDR/can be used in subnet spec
func getCIDR(baseCIDR string, position int, existing map[string]struct{}) (string, error) {
	_, ipNet, err := net.ParseCIDR(baseCIDR)
	if err != nil {
		return "", fmt.Errorf("cannot parse a CIDR with base %s, position %d", baseCIDR, position)
	}
	ip := ipNet.IP
	ones, _ := ipNet.Mask.Size()
	var i int
	for {
		if i > 255 {
			return "", fmt.Errorf("maximum number of IPs assigned")
		}
		ip[position] = byte(i)
		if _, ok := existing[fmt.Sprintf("%s/%d", ip.String(), ones)]; ok {
			// CIDR already assigned.
			i++
			continue
		}
		// CIDR does not exist yet, return.
		return fmt.Sprintf("%s/%d", ip.String(), ones), nil
	}
}

// generateProviderTemplates generates only the `provider.tpl` templates so terraform can destroy the infra if needed.
func (c *ClusterBuilder) generateProviderTemplates(current, desired *spec.ClusterInfo, clusterID, directory string, clusterData templates.ClusterData) error {
	currentNodepools := utils.GroupNodepoolsByProviderNames(current)
	desiredNodepools := utils.GroupNodepoolsByProviderNames(desired)

	// merge together into a single map instead of creating a new.
	for name, np := range desiredNodepools {
		// Continue if static node pool provider.
		if name.CloudProviderName == spec.StaticNodepoolInfo_STATIC_PROVIDER.String() {
			continue
		}

		if cnp, ok := currentNodepools[name]; !ok {
			currentNodepools[name] = np
		} else {
			// merge them together as different regions could be used.
			// (regions are used for generating the providers for various regions)
			for _, pool := range np {
				if found := utils.GetNodePoolByName(pool.GetName(), cnp); found == nil {
					currentNodepools[name] = append(currentNodepools[name], pool)
				}
			}
		}
	}

	info := desired
	if info == nil {
		info = current
	}

	for providerName, nodepools := range currentNodepools {
		// Continue if static node pool provider.
		if providerName.CloudProviderName == spec.StaticNodepoolInfo_STATIC_PROVIDER.String() {
			continue
		}

		p := nodepools[0].GetDynamicNodePool().GetProvider()
		templatesDownloadDir := filepath.Join(TemplatesRootDir, clusterID, p.SpecName)

		if err := templates.DownloadProvider(templatesDownloadDir, p); err != nil {
			msg := fmt.Sprintf("cluster %q failed to download template repository", clusterID)
			log.Error().Msgf(msg)
			return fmt.Errorf("%s: %w", msg, err)
		}

		g := templates.Generator{
			ID:                clusterID,
			TargetDirectory:   directory,
			ReadFromDirectory: templatesDownloadDir,
			TemplatePath:      templates.ExtractTargetPath(p.GetTemplates()),
		}

		err := g.GenerateProvider(&templates.Provider{
			ClusterData: clusterData,
			Provider:    p,
			Regions:     utils.GetRegions(utils.GetDynamicNodePools(nodepools)),
		})

		if err != nil {
			return fmt.Errorf("failed to generate provider templates: %w", err)
		}

		if err := utils.CreateKeyFile(utils.GetAuthCredentials(p), g.TargetDirectory, p.SpecName); err != nil {
			return fmt.Errorf("error creating provider credential key file for provider %s in %s : %w", p.SpecName, g.TargetDirectory, err)
		}
	}

	return nil
}
