{{- $clusterName := .ClusterName }}
{{- $clusterHash := .ClusterHash }}
{{- $index :=  0 }}
provider "hcloud" {
  token = "{{ (index .NodePools $index).Provider.Credentials }}" 
  alias = "lb_nodepool"
}

resource "hcloud_ssh_key" "claudie" {
  provider   = hcloud.lb_nodepool
  name       = "key-{{ $clusterName }}-{{ $clusterHash }}"
  public_key = file("./public.pem")

  labels = {
    "managed-by"      : "Claudie"
    "claudie-cluster" : "{{ $clusterName }}-{{ $clusterHash }}"
  }
}

resource "hcloud_firewall" "firewall" {
  provider = hcloud.lb_nodepool
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

  rule {
    direction  = "in"
    protocol   = "udp"
    port       = "51820"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  labels = {
    "managed-by"      : "Claudie"
    "claudie-cluster" : "{{ $clusterName }}-{{ $clusterHash }}"
  }
}


{{- range $nodepool := .NodePools }}
{{- range $node := $nodepool.Nodes }}
resource "hcloud_server" "{{ $node.Name }}" {
  provider     = hcloud.lb_nodepool
  name         = "{{ $node.Name }}"
  server_type  = "{{ $nodepool.ServerType }}"
  image        = "{{ $nodepool.Image }}"
  firewall_ids = [hcloud_firewall.firewall.id]
  datacenter   = "{{ $nodepool.Zone }}"
  ssh_keys = [
    hcloud_ssh_key.claudie.id,
  ]

  labels = {
    "managed-by"      : "Claudie"
    "claudie-cluster" : "{{ $clusterName }}-{{ $clusterHash }}"
  }
}

resource "hcloud_volume" "{{ $node.Name }}_volume" {
  provider  = hcloud.lb_nodepool
  name      = "{{ $node.Name }}-volume"
  size      = {{ $node.DiskSize }}
  server_id = hcloud_server.{{ $node.Name }}.id
  automount = true
  format    = "ext4"
}

resource "hcloud_volume_attachment" "{{ $node.Name }}_volume_att" {
  provider  = hcloud.lb_nodepool
  volume_id = hcloud_volume.{{ $node.Name }}_volume.id
  server_id = hcloud_server.{{ $node.Name }}.id
  automount = true
}
{{- end }}

output "{{ $nodepool.Name }}" {
  value = {
    {{- range $node := $nodepool.Nodes }}
    "${hcloud_server.{{ $node.Name }}.name}" = hcloud_server.{{ $node.Name }}.ipv4_address
    {{- end }}
  }
}
{{- end }}