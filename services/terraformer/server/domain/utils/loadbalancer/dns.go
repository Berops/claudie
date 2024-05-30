package loadbalancer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	comm "github.com/berops/claudie/internal/command"
	"github.com/berops/claudie/internal/utils"
	"github.com/berops/claudie/proto/pb"
	cluster_builder "github.com/berops/claudie/services/terraformer/server/domain/utils/cluster-builder"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/templates/backend"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/templates/provider"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/templates/templates"
	"github.com/berops/claudie/services/terraformer/server/domain/utils/terraform"
)

const (
	TemplatesRootDir = "services/terraformer/templates"
)

type DNS struct {
	ProjectName string
	ClusterName string
	ClusterHash string

	DesiredNodeIPs []string
	CurrentNodeIPs []string

	CurrentDNS *pb.DNS
	DesiredDNS *pb.DNS

	// SpawnProcessLimit represents a synchronization channel which limits the number of spawned terraform
	// processes. This values should always be non-nil and be buffered, where the capacity indicates
	// the limit.
	SpawnProcessLimit chan struct{}
}

// CreateDNSRecords creates DNS records for the Loadbalancer cluster.
func (d DNS) CreateDNSRecords(logger zerolog.Logger) (string, error) {
	sublogger := logger.With().Str("endpoint", d.DesiredDNS.Endpoint).Logger()

	clusterID := fmt.Sprintf("%s-%s", d.ClusterName, d.ClusterHash)
	dnsID := fmt.Sprintf("%s-dns", clusterID)
	dnsDir := filepath.Join(cluster_builder.Output, dnsID)

	terraform := terraform.Terraform{
		Directory:         dnsDir,
		SpawnProcessLimit: d.SpawnProcessLimit,
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		terraform.Stdout = comm.GetStdOut(clusterID)
		terraform.Stderr = comm.GetStdErr(clusterID)
	}

	if utils.ChangedDNSProvider(d.CurrentDNS, d.DesiredDNS) {
		sublogger.Info().Msg("Destroying old DNS records")
		if err := d.generateFiles(dnsID, dnsDir, d.CurrentDNS, d.CurrentNodeIPs); err != nil {
			return "", fmt.Errorf("error while creating dns .tf files for %s : %w", dnsID, err)
		}
		if err := terraform.Init(); err != nil {
			return "", err
		}
		if err := terraform.Destroy(); err != nil {
			return "", err
		}

		if err := os.RemoveAll(dnsDir); err != nil {
			return "", fmt.Errorf("error while removing files in dir %q: %w", dnsDir, err)
		}
		sublogger.Info().Msg("Old DNS records were successfully destroyed")
	}

	sublogger.Info().Msg("Creating new DNS records")
	if err := d.generateFiles(dnsID, dnsDir, d.DesiredDNS, d.DesiredNodeIPs); err != nil {
		return "", fmt.Errorf("error while creating dns .tf files for %s : %w", dnsID, err)
	}
	if err := terraform.Init(); err != nil {
		return "", err
	}
	if err := terraform.Apply(); err != nil {
		return "", err
	}

	outputID := fmt.Sprintf("%s-%s", clusterID, "endpoint")
	output, err := terraform.Output(clusterID)
	if err != nil {
		return "", fmt.Errorf("error while getting output from terraform for %s : %w", clusterID, err)
	}

	out, err := readDomain(output)
	if err != nil {
		return "", fmt.Errorf("error while reading output from terraform for %s : %w", clusterID, err)
	}

	sublogger.Info().Msg("DNS records were successfully set up")
	if err := os.RemoveAll(dnsDir); err != nil {
		return validateDomain(out.Domain[outputID]), fmt.Errorf("error while deleting files in %s: %w", dnsDir, err)
	}

	return validateDomain(out.Domain[outputID]), nil
}

// DestroyDNSRecords destroys DNS records for the Loadbalancer cluster.
func (d DNS) DestroyDNSRecords(logger zerolog.Logger) error {
	sublogger := logger.With().Str("endpoint", d.CurrentDNS.Endpoint).Logger()

	sublogger.Info().Msg("Destroying DNS records")
	dnsID := fmt.Sprintf("%s-%s-dns", d.ClusterName, d.ClusterHash)
	dnsDir := filepath.Join(cluster_builder.Output, dnsID)

	if err := d.generateFiles(dnsID, dnsDir, d.CurrentDNS, d.CurrentNodeIPs); err != nil {
		return fmt.Errorf("error while creating dns records for %s : %w", dnsID, err)
	}

	terraform := terraform.Terraform{
		Directory:         dnsDir,
		SpawnProcessLimit: d.SpawnProcessLimit,
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		terraform.Stdout = comm.GetStdOut(dnsID)
		terraform.Stderr = comm.GetStdErr(dnsID)
	}

	if err := terraform.Init(); err != nil {
		return err
	}
	if err := terraform.Destroy(); err != nil {
		return err
	}
	sublogger.Info().Msg("DNS records were successfully destroyed")

	if err := os.RemoveAll(dnsDir); err != nil {
		return fmt.Errorf("error while deleting files in %s : %w", dnsDir, err)
	}

	return nil
}

// generateFiles creates all the necessary terraform files used to create/destroy DNS.
func (d DNS) generateFiles(dnsID, dnsDir string, dns *pb.DNS, nodeIPs []string) error {
	backend := backend.Backend{
		ProjectName: d.ProjectName,
		ClusterName: dnsID,
		Directory:   dnsDir,
	}

	if err := backend.CreateTFFile(); err != nil {
		return err
	}

	providers := provider.Provider{
		ProjectName: d.ProjectName,
		ClusterName: dnsID,
		Directory:   dnsDir,
	}

	if err := providers.CreateProviderDNS(dns); err != nil {
		return err
	}

	repo := templates.Repository{
		TemplatesRootDirectory: TemplatesRootDir,
	}

	if err := repo.Download(dns.GetTemplates()); err != nil {
		return fmt.Errorf("failed to download template repository for %q: %w", dnsID, err)
	}

	g := templates.DNSGenerator{
		DnsID:             dnsID,
		TargetDirectory:   dnsDir,
		ReadFromDirectory: TemplatesRootDir,
		DNS:               dns,
	}

	data := templates.DNSData{
		DNSZone:      dns.DnsZone,
		HostnameHash: dns.Hostname,
		ClusterName:  d.ClusterName,
		ClusterHash:  d.ClusterHash,
		NodeIPs:      nodeIPs,
		Provider:     dns.Provider,
	}

	if err := g.GenerateDNS(&data); err != nil {
		return fmt.Errorf("failed to generate dns templates for %q: %w", dnsID, err)
	}

	return nil
}

// validateDomain validates the domain does not start with ".".
func validateDomain(s string) string {
	if s[len(s)-1] == '.' {
		return s[:len(s)-1]
	}
	return s
}

// readDomain reads full domain from terraform output.
func readDomain(data string) (templates.DNSDomain, error) {
	var result templates.DNSDomain
	err := json.Unmarshal([]byte(data), &result.Domain)
	return result, err
}
