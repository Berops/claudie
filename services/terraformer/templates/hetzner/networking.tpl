{{- $clusterName := .ClusterData.ClusterName }}
{{- $clusterHash := .ClusterData.ClusterHash }}

{{- range $_, $region := .Regions }}
{{- $specName := $.Provider.SpecName }}

resource "hcloud_ssh_key" "claudie_{{ $region }}" {
  provider   = hcloud.nodepool_{{ $region }}_{{ $specName }}
  name       = "key-{{ $clusterName }}-{{ $clusterHash }}"
  public_key = file("./public.pem")

  labels = {
    "managed-by"      : "Claudie"
    "claudie-cluster" : "{{ $clusterName }}-{{ $clusterHash }}"
  }
}

resource "hcloud_firewall" "firewall_{{ $region }}" {
  provider = hcloud.nodepool_{{ $region }}_{{ $specName }}
  name     = "{{ $clusterName }}-{{ $clusterHash }}-firewall"
  rule {
    direction  = "in"
    protocol   = "icmp"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "22"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction  = "in"
    protocol   = "udp"
    port       = "51820"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

{{- if eq $.ClusterData.ClusterType "LB" }}
  {{- range $role := index $.Metadata "roles" }}
  rule {
    direction  = "in"
    protocol   = "{{ $role.Protocol }}"
    port       = "{{ $role.Port }}"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
  {{- end }}
{{- end }}

{{- if eq $.ClusterData.ClusterType "K8s" }}
  {{- if index $.Metadata "loadBalancers" | targetPorts | isMissing 6443 }}
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "6443"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
  {{- end }}
{{- end }}

  labels = {
    "managed-by"      : "Claudie"
    "claudie-cluster" : "{{ $clusterName }}-{{ $clusterHash }}"
  }
}
{{- end }}
