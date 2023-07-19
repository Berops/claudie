package claudie_provider

import (
	"fmt"
	"time"

	"github.com/berops/claudie/internal/nodes"
	kuberNodes "github.com/berops/claudie/services/kuber/server/domain/utils/nodes"
	k8sV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Default node template with prefilled conditions.
	defaultNodeTemplate = k8sV1.Node{
		Status: k8sV1.NodeStatus{
			Conditions: buildReadyConditions(),
		},
	}
)

// getNodeGroupTemplateNodeInfo returns a template for the new node from specified node group.
func (c *ClaudieCloudProvider) getNodeGroupTemplateNodeInfo(nodeGroupId string) (*k8sV1.Node, error) {
	if ngc, ok := c.nodesCache[nodeGroupId]; ok {
		// Create a new node struct.s
		node := defaultNodeTemplate
		// Fill dynamic fields.
		l, err := nodes.GetAllLabels(ngc.nodepool, c.nodeManager)
		if err != nil {
			return nil, err
		}
		node.Labels = l
		node.Spec.Taints = nodes.GetAllTaints(ngc.nodepool)
		node.Status.Capacity = c.nodeManager.GetCapacity(ngc.nodepool)
		node.Status.Allocatable = node.Status.Capacity
		node.Spec.ProviderID = fmt.Sprintf(kuberNodes.ProviderIdFormat, fmt.Sprintf("%s-N", ngc.nodepool.Name))
		return &node, nil
	}
	return nil, nil
}

// buildReadyConditions returns default ready condition for the node.
func buildReadyConditions() []k8sV1.NodeCondition {
	lastTransition := time.Now().Add(-time.Minute)
	return []k8sV1.NodeCondition{
		{
			Type:               k8sV1.NodeReady,
			Status:             k8sV1.ConditionTrue,
			LastTransitionTime: metaV1.Time{Time: lastTransition},
		},
		{
			Type:               k8sV1.NodeNetworkUnavailable,
			Status:             k8sV1.ConditionFalse,
			LastTransitionTime: metaV1.Time{Time: lastTransition},
		},
		{
			Type:               k8sV1.NodeDiskPressure,
			Status:             k8sV1.ConditionFalse,
			LastTransitionTime: metaV1.Time{Time: lastTransition},
		},
		{
			Type:               k8sV1.NodeMemoryPressure,
			Status:             k8sV1.ConditionFalse,
			LastTransitionTime: metaV1.Time{Time: lastTransition},
		},
	}
}
