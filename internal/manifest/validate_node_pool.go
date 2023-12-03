package manifest

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	k8sV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Validate validates the parsed data inside the NodePool section of the manifest.
// It checks for missing/invalid filled out values defined in the NodePool section of
// the manifest.
func (p *NodePool) Validate(m *Manifest) error {
	// check for name uniqueness across node pools.
	names := make(map[string]bool)

	for _, n := range p.Dynamic {
		if err := n.Validate(); err != nil {
			return fmt.Errorf("failed to validate DynamicNodePool %q: %w", n.Name, err)
		}
		// check if the provider is defined in the manifest
		if _, err := m.GetProvider(n.ProviderSpec.Name); err != nil {
			return fmt.Errorf("provider %q specified for DynamicNodePool %q doesn't exists", n.ProviderSpec.Name, n.Name)
		}

		// check if the name is already used by a different node pool
		if _, ok := names[n.Name]; ok {
			return fmt.Errorf("name %q is used across multiple node pools, must be unique", n.Name)
		}
		names[n.Name] = true

		// check if the count and autoscaler are mutually exclusive
		if n.Count != 0 && n.AutoscalerConfig.isDefined() {
			return fmt.Errorf("nodepool %s cannot have both, autoscaler enabled and \"count\" defined", n.Name)
		}
		if err := checkTaints(n.Taints); err != nil {
			return fmt.Errorf("nodepool %s has incorrectly defined taints : %w", n.Name, err)
		}
		if err := checkLabels(n.Labels); err != nil {
			return fmt.Errorf("nodepool %s has incorrectly defined labels : %w", n.Name, err)
		}
	}

	for _, n := range p.Static {
		if err := n.Validate(); err != nil {
			return fmt.Errorf("failed to validate StaticNodePool %q: %w", n.Name, err)
		}

		// check if the name is already used by a different node pool
		if _, ok := names[n.Name]; ok {
			return fmt.Errorf("name %q is used across multiple node pools, must be unique", n.Name)
		}
		names[n.Name] = true
		if err := checkTaints(n.Taints); err != nil {
			return fmt.Errorf("nodepool %s has incorrectly defined taints : %w", n.Name, err)
		}
		if err := checkLabels(n.Labels); err != nil {
			return fmt.Errorf("nodepool %s has incorrectly defined labels : %w", n.Name, err)
		}
	}

	return nil
}

func (d *DynamicNodePool) Validate() error {
	if d.StorageDiskSize < 0 || d.StorageDiskSize >= 50 {
		return fmt.Errorf("storageDiskSize size must be either -1 or >= 50")
	}

	return validator.New().Struct(d)
}
func (s *StaticNodePool) Validate() error { return validator.New().Struct(s) }

func (a *AutoscalerConfig) isDefined() bool { return a.Min >= 0 && a.Max > 0 }

func checkTaints(taints []k8sV1.Taint) error {
	for _, t := range taints {
		// Check if effect is supported
		if !(t.Effect == k8sV1.TaintEffectNoSchedule || t.Effect == k8sV1.TaintEffectNoExecute || t.Effect == k8sV1.TaintEffectPreferNoSchedule) {
			return fmt.Errorf("taint effect \"%s\" is not supported", t.Effect)
		}
	}
	return nil
}

func checkLabels(labels map[string]string) error {
	for k, v := range labels {
		if errs := validation.IsQualifiedName(k); len(errs) > 0 {
			return fmt.Errorf("key %v is not valid  : %v", k, errs)
		}
		if errs := validation.IsValidLabelValue(v); len(errs) > 0 {
			return fmt.Errorf("value %v is not valid  : %v", v, errs)
		}
	}
	return nil
}
