# InputManifest API reference
InputManifest is a definition of the user's infrastructure. It contains cloud provider specification, nodepool specification, Kubernetes and loadbalancer clusters.

## Status
Most recently observed status of the InputManifest

## Spec

Specification of the desired behavior of the InputManifest

- `providers` [Providers](#providers)

  Providers is a list of defined cloud provider configuration that will be used in infrastructure provisioning.

- `nodepools` [Nodepools](#nodepools)

  Describes nodepools used for either kubernetes clusters or loadbalancer cluster defined in this manifest.

- `kubernetes` [Kubernetes](#kubernetes)

  List of Kubernetes cluster this manifest will manage.

- `loadBalancers` [Loadbalancer](#loadbalancer)

  List of loadbalancer clusters the Kubernetes clusters may use.

## Providers

Contains configurations for supported cloud providers. At least one provider
needs to be defined.

- `name`

  The name of the provider specification. It has to be unique across all providers.

- `providerType`

  Type of a provider. The providerType defines mandatory fields that has to be included for a specific provider. A list of available providers can be found at [providers section](./providers). Allowed values are:

  | Value        | Description                              |
  | ------------ | ---------------------------------------- |
  | `aws`        | [AWS](#aws) provider type                |
  | `azure`      | [Azure](#azure) provider type            |
  | `cloudflare` | [Cloudflare](#cloudflare) provider type  |
  | `gcp`        | [GCP](#gcp) provider type                |
  | `hetzner`    | [Hetzner](#hetzner) provider type        |
  | `hetznerdns` | [Hetzner](#hetznerdns) DNS provider type |
  | `oci`        | [OCI](#oci) provider type                |
  
- `secretRef` [SecretRef](#secretref)

  Represents a Secret Reference. It has enough information to retrieve secret in any namespace.

Support for more cloud providers is in the [roadmap](https://github.com/berops/claudie/blob/master/docs/roadmap/roadmap.md).


## SecretRef
  
  SecretReference represents a Kubernetes Secret Reference. It has enough information to retrieve secret in any namespace.

- `name`

  Name of the secret, which holds data for the particular cloud provider instance.

- `namespace`

  Namespace of the secret which holds data for the particular cloud provider instance.

### Cloudflare

The fields that need to be included in a Kubernetes Secret resource to utilize the Cloudflare provider.
To find out how to configure Cloudflare follow the instructions [here](./providers/cloudflare.md)

- `apitoken`

  Credentials for the provider (API token).

## HetznerDNS

The fields that need to be included in a Kubernetes Secret resource to utilize the HetznerDNS provider.
To find out how to configure HetznerDNS follow the instructions [here](./providers/hetzner.md)

- `apitoken`

  Credentials for the provider (API token).

## GCP

The fields that need to be included in a Kubernetes Secret resource to utilize the GCP provider.
To find out how to configure GCP provider and service account, follow the instructions [here](./providers/gcp.md).

- `credentials`

  Credentials for the provider. Stringified JSON service account key.

- `gcpproject`

  Project id of an already existing GCP project where the infrastructure is to be created.

## Hetzner

The fields that need to be included in a Kubernetes Secret resource to utilize the Hetzner provider.
To find out how to configure Hetzner provider and service account, follow the instructions [here](./providers/hetzner.md).

- `credentials`

  Credentials for the provider (API token).

## OCI

The fields that need to be included in a Kubernetes Secret resource to utilize the OCI provider.
To find out how to configure OCI provider and service account, follow the instructions [here](./providers/oci.md).

- `privatekey`

  [Private key](https://docs.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#two) used to authenticate to the OCI.

- `keyfingerprint`

  Fingerprint of the user-supplied private key.

- `tenancyocid`
  
  OCID of the tenancy where `privateKey` is added as an API key

- `userocid`
  
  OCID of the user in the supplied tenancy

- `compartmentocid`

  OCID of the [compartment](https://docs.oracle.com/en/cloud/paas/integration-cloud/oracle-integration-oci/creating-oci-compartment.html) where VMs/VCNs/... will be created

## AWS

The fields that need to be included in a Kubernetes Secret resource to utilize the AWS provider.
To find out how to configure AWS provider and service account, follow the instructions [here](./providers/aws.md).

- `accesskey`

  Access key ID for your AWS account.

- `secretkey`

  Secret key for the Access key specified above.

## Azure

The fields that need to be included in a Kubernetes Secret resource to utilize the Azure provider.
To find out how to configure Azure provider and service account, follow the instructions [here](./providers/azure.md).

- `subscriptionid`

  Subscription ID of your subscription in Azure.

- `tenantid`
  
  Tenant ID of your tenancy in Azure.

- `clientid`

  Client ID of your client. The Claudie is design to use a service principal with appropriate permissions.

- `clientsecret`
  
  Client secret generated for your client.

## Nodepools

Collection of static and dynamic nodepool specification, to be referenced in the `kubernetes` or `loadBalancer` clusters.

- `dynamic` [Dynamic](#dynamic)

  List of dynamically to-be-created nodepools of not yet existing machines, used for Kubernetes or loadbalancer clusters.
  
  These are only blueprints, and will only be created per reference in `kubernetes` or `loadBalancer` clusters. E.g. if the nodepool isn't used, it won't even be created. Or if the same nodepool is used in two different clusters, it will be created twice.
In OOP analogy, a dynamic nodepool would be a class that would get instantiated `N >= 0` times depending on which clusters reference it.

- `static` [Static](#static)

  List of static nodepools of already existing machines, not created by of Claudie, used for Kubernetes or loadbalancer clusters. Typically, these would be on-premises machines.

## Dynamic

Dynamic nodepools are defined for cloud provider machines that Claudie is expected to create.

- `name`

  Name of the nodepool. Each nodepool will have a random hash appended to the name, so the whole name will be of format `<name>-<hash>`.

- `provideSpec` [Provider spec](#provider-spec)

  Collection of provider data to be used while creating the nodepool.  

- `count`

  Number of the nodes in the nodepool. Mutually exclusive with `autoscaler`.

- `serverType`
  
  Type of the machines in the nodepool.
  
  Currently, only AMD64 machines are supported.

- `image`

  OS image of the machine.
  
  Currently, only Ubuntu 22.04 AMD64 images are supported.

- `storageDiskSize`

  Size of the storage disk on the nodes in the nodepool in `GB`. The OS disk is created automatically with predefined size of `100GB` for kubernetes nodes and `50GB` for Loadbalancer nodes.
  
  Default value is `50`, minimum value is `50`. Value is used only for compute nodes.

  This field is optional, however, if compute nodepool does not define it, default value is used for creation of storage disk. Control nodepools and Loadbalancer nodepools ignore this field.

- `autoscaler` [Autoscaler Configuration](#autoscaler-configuration)
  
  Autoscaler configuration for this nodepool. Mutually exclusive with `count`.

## Provider Spec

Provider spec is an additional specification built on top of the data from any of the provider instance. Here are provider configuration examples for each individual provider: [aws](providers/aws.md), [azure](providers/azure.md), [gcp](providers/gcp.md), [cloudflare](providers/cloudflare.md), [hetzner](providers/hetzner.md) and [oci](providers/oci.md).

- `name`

  Name of the provider instance specified in [providers](#providers)

- `region`

  Region of the nodepool.

- `zone`

  Zone of the nodepool.

## Autoscaler Configuration

Autoscaler configuration on per nodepool basis. Defines the number of nodes, autoscaler will scale up or down specific nodepool.

- `min`
  
  Minimum number of nodes in nodepool.

- `max`

  Maximum number of nodes in nodepool.

## Static

Static nodepools are defined for static machines which Claudie will not manage. Used for on premise nodes.

- `name`

  Name of the static nodepool.

- `nodes` [Static Node](#static-node)

  List of static nodes for a particular static nodepool.

## Static node

Static node defines single static node from a static nodepool.

- `endpoint`

  Endpoint under which Claudie will access this node.

- `secretRef` [SecretRef](#secretref)

  Secret from which private key will be taken used to SSH into the machine (as root).

  The field in the secret must be `privatekey`, i.e.
  
  ```yaml
  apiVersion: v1
  type: Opaque
  kind: Secret
    name: private-key-node-1
    namespace: claudie-secrets
  data:
    privatekey: <base64 encoded private key>
  ```

## Kubernetes

Defines Kubernetes clusters.

- `clusters` [Cluster-k8s](#cluster-k8s)

  List of Kubernetes clusters Claudie will create.

## Cluster-k8s

Collection of data used to define a Kubernetes cluster.

- `name`

  Name of the Kubernetes cluster. Each cluster will have a random hash appended to the name, so the whole name will be of format `<name>-<hash>`.

- `version`

  Kubernetes version of the cluster.

  Version should be defined in format `vX.Y`. In terms of supported versions of Kubernetes, Claudie follows `kubeone` releases and their supported versions. The current `kubeone` version used in Claudie is `1.5`. To see the list of supported versions, please refer to `kubeone` [documentation](https://docs.kubermatic.com/kubeone/v1.5/architecture/compatibility/supported-versions/#supported-kubernetes-versions).

- `network`

  Network range for the VPN of the cluster. The value should be defined in format `A.B.C.D/mask`.

- `pools`

  List of nodepool names this cluster will use. Remember that nodepools defined in [nodepools](#nodepools) are only "blueprints". The actual nodepool will be created once referenced here.

## LoadBalancer

Defines loadbalancer clusters.

- `roles` [Role](#role)
  
  List of roles loadbalancers use to forward the traffic. Single role can be used in multiple loadbalancer clusters.

- `clusters` [Cluster-lb](#cluster-lb)

  List of loadbalancer clusters used in the Kubernetes clusters defined under [clusters](#cluster-k8s).

## Role

Role defines a concrete loadbalancer configuration. Single loadbalancer can have multiple roles.

- `name`

  Name of the role. Used as a reference in [clusters](#cluster-lb).

- `protocol`

  Protocol of the rule. Allowed values are:

  | Value | Description                |
  | ----- | -------------------------- |
  | `tcp` | Role will use TCP protocol |
  | `udp` | Role will use UDP protocol |

- `port`

  Port of the incoming traffic on the loadbalancer.

- `targetPort`

  Port where loadbalancer forwards the traffic.

- `target`

  Defines a target group of nodes. Allowed values are:

  | Value             | Description                          |
  | ----------------- | ------------------------------------ |
  | `k8sAllNodes`     | All nodes in the cluster             |
  | `k8sControlPlane` | Only control/master nodes in cluster |
  | `k8sComputePlane` | Only compute/worker nodes in cluster |

## Cluster-lb

Collection of data used to define a loadbalancer cluster.

- `name`

  Name of the loadbalancer.

- `roles`
  
  List of roles the loadbalancer uses.

- `dns` [DNS](#dns)
  
  Specification of the loadbalancer's DNS record.
  
- `targetedK8s`

  Name of the Kubernetes cluster targetted by this loadbalancer.

- `pools`

  List of nodepool names this loadbalancer will use. Remember, that nodepools defined in [nodepools](#nodepools) are only "blueprints". The actual nodepool will be created once referenced here.

## DNS

Collection of data Claudie uses to create a DNS record for the loadbalancer.

- `dnsZone`

  DNS zone inside which the records will be created. GCP/AWS/OCI/Azure/Cloudflare/Hetzner DNS zone is accepted.

  The record created in this zone must be accessible to the public. Therefore, a public DNS zone is required.

- `provider`

  Name of [provider](#providers) to be used for creating an A record entry in defined DNS zone.

- `hostname`
  
  Custom hostname for your A record. If left empty, the hostname will be a random hash.
