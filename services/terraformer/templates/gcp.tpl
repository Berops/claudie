{{- $clusterName := .ClusterName}}
{{- $clusterHash := .ClusterHash}}
{{- $index :=  0}}
{{- range $i, $region := .Regions}}
provider "google" {
  credentials = "${file("{{ (index $.NodePools $index).Provider.SpecName }}")}"
  project     = "{{ (index $.NodePools 0).Provider.GcpProject }}"
  region      = "{{ $region }}"
  alias       = "k8s_nodepool_{{ $region }}"
}

resource "google_compute_network" "network_{{ $clusterName}}_{{ $clusterHash}}_{{ $region }}" {
  provider                = google.k8s_nodepool_{{ $region }}
  name                    = "{{ $clusterName}}-{{ $clusterHash}}-{{ $region }}-network"
  auto_create_subnetworks = false
  description             = "Managed by Claudie for cluster {{ $clusterName }}-{{ $clusterHash }}"
}

resource "google_compute_firewall" "firewall_{{ $clusterName}}_{{ $clusterHash}}_{{ $region }}" {
  provider    = google.k8s_nodepool_{{ $region }}
  name        = "{{ $clusterName}}-{{ $clusterHash}}-{{ $region }}-firewall"
  network     = google_compute_network.network_{{ $clusterName}}_{{ $clusterHash}}_{{ $region }}.self_link
  description = "Managed by Claudie for cluster {{ $clusterName }}-{{ $clusterHash }}"

  allow {
    protocol = "UDP"
    ports    = ["51820"]
  }

  {{- if index $.Metadata "loadBalancers" | targetPorts | isMissing 6443 }}
  allow {
      protocol = "TCP"
      ports    = ["6443"]
  }
  {{- end }}

  allow {
      protocol = "TCP"
      ports    = ["22"]
  }

  allow {
      protocol = "icmp"
   }

  source_ranges = [
      "0.0.0.0/0",
   ]
}

{{- end }}

{{- range $i, $nodepool := .NodePools }}
resource "google_compute_subnetwork" "{{ $nodepool.Name }}_subnet" {
  provider      = google.k8s_nodepool_{{ $nodepool.Region }}
  name          = "{{ $nodepool.Name }}-{{ $clusterHash }}-subnet"
  network       = google_compute_network.network_{{ $clusterName}}_{{ $clusterHash}}_{{ $nodepool.Region }}.self_link
  ip_cidr_range = "{{getCIDR "10.0.0.0/24" 2 $i }}"
  description   = "Managed by Claudie for cluster {{ $clusterName }}-{{ $clusterHash }}"
}

{{- range $node := $nodepool.Nodes }}
resource "google_compute_instance" "{{ $node.Name }}" {
  provider                  = google.k8s_nodepool_{{ $nodepool.Region }}
  zone                      = "{{ $nodepool.Zone }}"
  name                      = "{{ $node.Name }}"
  machine_type              = "{{ $nodepool.ServerType }}"
  description   = "Managed by Claudie for cluster {{ $clusterName }}-{{ $clusterHash }}"
  allow_stopping_for_update = true
  boot_disk {
    initialize_params {
      size = "{{ $nodepool.DiskSize }}"
      image = "{{ $nodepool.Image }}"
    }
  }
  network_interface {
    subnetwork = google_compute_subnetwork.{{ $nodepool.Name }}_subnet.self_link
    access_config {}
  }
  metadata = {
    ssh-keys = "root:${file("./public.pem")}"
  }
  metadata_startup_script = "echo 'PermitRootLogin without-password' >> /etc/ssh/sshd_config && echo 'PubkeyAuthentication yes' >> /etc/ssh/sshd_config && service sshd restart"
  
  labels = {
    managed-by = "claudie"
    claudie-cluster = "{{ $clusterName }}-{{ $clusterHash }}"
  }
}
{{- end }}

output "{{ $nodepool.Name }}" {
  value = {
  {{- range $node := $nodepool.Nodes }}
    "${google_compute_instance.{{ $node.Name }}.name}" = google_compute_instance.{{ $node.Name }}.network_interface.0.access_config.0.nat_ip
  {{- end }}
  }
}
{{- end }}


