package provider

import (
	_ "embed"
	"fmt"

	"github.com/berops/claudie/internal/templateUtils"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
)

//go:embed providers.tpl
var providersTemplate string

// Provider package struct
type Provider struct {
	ProjectName string
	ClusterName string
	Directory   string
}

// templateData is data structure passed to providers.tpl
type templateData struct {
	Gcp          bool
	Hetzner      bool
	Aws          bool
	Oci          bool
	Azure        bool
	Cloudflare   bool
	HetznerDNS   bool
	GenesisCloud bool
}

// CreateProviderDNS creates provider file used for DNS management.
func (p Provider) CreateProviderDNS(dns *pb.DNS) error {
	template := templateUtils.Templates{Directory: p.Directory}

	tpl, err := templateUtils.LoadTemplate(providersTemplate)
	if err != nil {
		return fmt.Errorf("error while parsing template file providers.tpl for cluster %s: %w", p.ClusterName, err)
	}

	var data templateData
	getDNSProvider(dns, &data)
	return template.Generate(tpl, "providers.tf", data)
}

// CreateProvider creates provider file used for infrastructure management.
func (p Provider) CreateProvider(currentCluster, desiredCluster *pb.ClusterInfo) error {
	template := templateUtils.Templates{Directory: p.Directory}

	var data templateData

	getProvidersUsed(utils.GetDynamicNodePoolsFromCI(currentCluster), &data)
	getProvidersUsed(utils.GetDynamicNodePoolsFromCI(desiredCluster), &data)

	tpl, err := templateUtils.LoadTemplate(providersTemplate)
	if err != nil {
		return fmt.Errorf("error while parsing template file providers.tpl for cluster %s : %w", p.ClusterName, err)
	}

	if err := template.Generate(tpl, "providers.tf", data); err != nil {
		return fmt.Errorf("error while creating provider.tf for %s : %w", p.ClusterName, err)
	}

	return nil
}

// getProvidersUsed modifies templateData to reflect current providers used.
func getProvidersUsed(nodepools []*pb.DynamicNodePool, data *templateData) {
	if len(nodepools) == 0 {
		return
	}

	for _, nodepool := range nodepools {
		if nodepool.Provider.CloudProviderName == "gcp" {
			data.Gcp = true
		}
		if nodepool.Provider.CloudProviderName == "hetzner" {
			data.Hetzner = true
		}
		if nodepool.Provider.CloudProviderName == "aws" {
			data.Aws = true
		}
		if nodepool.Provider.CloudProviderName == "oci" {
			data.Oci = true
		}
		if nodepool.Provider.CloudProviderName == "azure" {
			data.Azure = true
		}
		if nodepool.Provider.CloudProviderName == "genesiscloud" {
			data.GenesisCloud = true
		}
	}
}

// getProvidersUsed modifies templateData to reflect current providers used in DNS.
func getDNSProvider(dns *pb.DNS, data *templateData) {
	if dns == nil {
		return
	}

	switch dns.Provider.CloudProviderName {
	case "gcp":
		data.Gcp = true
	case "hetzner":
		data.Hetzner = true
	case "aws":
		data.Aws = true
	case "oci":
		data.Oci = true
	case "azure":
		data.Azure = true
	case "cloudflare":
		data.Cloudflare = true
	case "hetznerdns":
		data.HetznerDNS = true
	}
}
