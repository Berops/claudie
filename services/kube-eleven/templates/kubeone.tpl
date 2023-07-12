apiVersion: kubeone.k8c.io/v1beta2
kind: KubeOneCluster
name: '{{ .ClusterName }}'

versions:
  kubernetes: '{{ .KubernetesVersion }}'

clusterNetwork:
  cni:
    external: {}

cloudProvider:
  none: {}
  external: false

addons:
  enable: true
  # In case when the relative path is provided, the path is relative
  # to the KubeOne configuration file.
  path: "../../addons"

apiEndpoint:
  host: '{{ .APIEndpoint }}'
  port: 6443

{{- $privateKey := "./private.pem" }}
controlPlane:
  hosts:
{{- range $nodepool := .Nodepools }}
  {{- range $nodeInfo := $nodepool.Nodes }}
    {{- if ge $nodeInfo.Node.NodeType 1}}
  - publicAddress: '{{ $nodeInfo.Node.Public }}'
    privateAddress: '{{ $nodeInfo.Node.Private }}'
    sshUsername: root
    {{- if $nodepool.IsDynamic }}
    sshPrivateKeyFile: '{{ $privateKey }}'
    {{- else }}
    sshPrivateKeyFile: './{{ $nodeInfo.Name }}.pem'
    {{- end }}
    hostname: '{{ $nodeInfo.Name }}'
    {{- if eq $nodeInfo.Node.Public $.APIEndpoint }}
    isLeader: true
    {{- end }}
    taints:
    - key: "node-role.kubernetes.io/control-plane"
      effect: "NoSchedule"
    {{- end}}
  {{- end}}
{{- end}}

staticWorkers:
  hosts:
{{- range $nodepool := .Nodepools }}
  {{- range $nodeInfo := $nodepool.Nodes }}
    {{- if eq $nodeInfo.Node.NodeType 0}}
  - publicAddress: '{{ $nodeInfo.Node.Public }}'
    privateAddress: '{{ $nodeInfo.Node.Private }}'
    sshUsername: root
    {{- if $nodepool.IsDynamic }}
    sshPrivateKeyFile: '{{ $privateKey }}'
    {{- else }}
    sshPrivateKeyFile: './{{ $nodeInfo.Name }}.pem'
    {{- end }}
    hostname: '{{ $nodeInfo.Name }}'
    {{- end}}
  {{- end}}
{{- end}}

machineController:
  deploy: false
