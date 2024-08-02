package nodes

import (
	"encoding/json"
	"fmt"
	"strings"

	comm "github.com/berops/claudie/internal/command"
	"github.com/berops/claudie/internal/kubectl"
	"github.com/berops/claudie/internal/nodes"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb/spec"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	ProviderIdFormat          = "claudie://%s"
	patchProviderIDPathFormat = "{\"spec\":{\"providerID\":\"%s\"}}"
)

type patchJson struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

type PatchAnnotations struct {
	MetadataAnnotations `json:"metadata"`
}

type MetadataAnnotations struct {
	Annotations map[string]string `json:"annotations"`
}

type Patcher struct {
	clusterID        string
	desiredNodepools []*spec.NodePool
	kc               kubectl.Kubectl
	logger           zerolog.Logger
}

func NewPatcher(cluster *spec.K8Scluster, logger zerolog.Logger) *Patcher {
	kc := kubectl.Kubectl{Kubeconfig: cluster.Kubeconfig, MaxKubectlRetries: 3}
	clusterID := utils.GetClusterID(cluster.ClusterInfo)
	if log.Logger.GetLevel() == zerolog.DebugLevel {
		kc.Stdout = comm.GetStdOut(clusterID)
		kc.Stderr = comm.GetStdErr(clusterID)
	}

	return &Patcher{
		kc:               kc,
		desiredNodepools: cluster.ClusterInfo.NodePools,
		clusterID:        clusterID,
		logger:           logger,
	}
}

func (p *Patcher) PatchProviderID() error {
	var err error
	for _, np := range p.desiredNodepools {
		for _, node := range np.GetNodes() {
			nodeName := strings.TrimPrefix(node.Name, fmt.Sprintf("%s-", p.clusterID))
			patchPath := fmt.Sprintf(patchProviderIDPathFormat, fmt.Sprintf(ProviderIdFormat, nodeName))
			if err1 := p.kc.KubectlPatch("node", nodeName, patchPath); err1 != nil {
				p.logger.Err(err1).Str("node", nodeName).Msgf("Error while patching node with patch %s", patchPath)
				err = fmt.Errorf("error while patching one or more nodes with providerID")
			}
		}
	}
	return err
}

func (p *Patcher) PatchLabels() error {
	var err error
	for _, np := range p.desiredNodepools {
		nodeLabels, err1 := nodes.GetAllLabels(np, nil)
		if err1 != nil {
			return fmt.Errorf("failed to create labels for %s : %w, %w", np.Name, err, err1)
		}

		for _, node := range np.Nodes {
			nodeName := strings.TrimPrefix(node.Name, fmt.Sprintf("%s-", p.clusterID))
			for key, value := range nodeLabels {
				patchPath, err1 := buildJSONPatchString("replace", "/metadata/labels/"+key, value)
				if err1 != nil {
					return fmt.Errorf("failed to create label %s patch path for %s : %w, %w", key, np.Name, err, err1)
				}
				if err1 := p.kc.KubectlPatch("node", nodeName, patchPath, "--type", "json"); err1 != nil {
					p.logger.Err(err1).Str("node", nodeName).Msgf("Failed to patch labels on node with path %s", patchPath)
					err = fmt.Errorf("error while patching one or more nodes with labels")
				}
			}
		}
	}
	return err
}

func (p *Patcher) PatchAnnotations() error {
	var err error
	for _, np := range p.desiredNodepools {
		nodeAnnotations := np.Annotations
		for _, node := range np.Nodes {
			nodeName := strings.TrimPrefix(node.Name, fmt.Sprintf("%s-", p.clusterID))
			patchPath, err1 := buildJSONAnnotationPatch(nodeAnnotations)
			if err1 != nil {
				return fmt.Errorf("failed to create annotations patch for %s : %w, %w", np.Name, err, err1)
			}
			if err1 := p.kc.KubectlPatch("node", nodeName, patchPath, "--type", "merge"); err1 != nil {
				p.logger.Err(err1).Str("node", nodeName).Msgf("Failed to patch annotations on node with path %s", patchPath)
				err = fmt.Errorf("error while patching one or more nodes with annotations")
			}
		}
	}
	return err
}

func (p *Patcher) PatchTaints() error {
	var err error
	for _, np := range p.desiredNodepools {
		patchPath, err1 := buildJSONPatchString("replace", "/spec/taints", nodes.GetAllTaints(np))
		if err1 != nil {
			return fmt.Errorf("failed to create taints patch path for %s : %w", np.Name, err)
		}
		for _, node := range np.Nodes {
			nodeName := strings.TrimPrefix(node.Name, fmt.Sprintf("%s-", p.clusterID))
			if err1 := p.kc.KubectlPatch("node", nodeName, patchPath, "--type", "json"); err1 != nil {
				p.logger.Err(err1).Str("node", nodeName).Msgf("Failed to patch taints on node with path %s", patchPath)
				err = fmt.Errorf("error while patching one or more nodes with taints")
			}
		}
	}
	return err
}

func buildJSONAnnotationPatch(data map[string]string) (string, error) {
	metadata := PatchAnnotations{
		MetadataAnnotations{
			Annotations: data,
		},
	}
	jsonPatch, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(jsonPatch), nil
}

func buildJSONPatchString(op, path string, value any) (string, error) {
	patchJson := patchJson{Op: op, Path: path, Value: value}
	b, err := json.Marshal(patchJson)
	if err != nil {
		return "", err
	}
	return "[" + string(b) + "]", nil
}
