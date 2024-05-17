# Claudie storage solution

## Concept

Running stateful workloads is a complex task, even more so when considering the multi-cloud environment. Claudie therefore needs to be able to accommodate stateful workloads, regardless of the underlying infrastructure providers.

Claudie orchestrates storage on the kubernetes cluster nodes by creating one "storage cluster" across multiple providers. This "storage cluster" has a series of `zones`, one for each cloud provider instance. Each `zone` then stores its own persistent volume data.

This concept is translated into longhorn implementation, where each `zone` is represented by a Storage Class which is backed up by the nodes defined under the same cloud provider instance. Furthermore, each node uses separate disk to the one, where OS is installed, to assure clear data separation. The size of the storage disk can be configured in `storageDiskSize` field of the nodepool specification.

## Longhorn

A Claudie-created cluster comes with the `longhorn` deployment preinstalled and ready to be used. By default, only **worker** nodes are used to store data.

Longhorn installed in the cluster is set up in a way that it provides one default `StorageClass` called `longhorn`, which, if used, creates a volume that is then replicated across random nodes in the cluster.

Besides the default storage class, Claudie can also create custom storage classes, which force persistent volumes to be created on specific nodes based on the provider instance they have. In other words, you can use a specific provider instance to provision nodes for your storage needs, while using another provider instance for computing tasks.

## Example

To follow along, have a look at the example of `InputManifest` below.

``` yaml title="storage-classes-example.yaml"
apiVersion: claudie.io/v1beta1
kind: InputManifest
metadata:
  name: ExampleManifestForStorageClasses
  labels:
    app.kubernetes.io/part-of: claudie
spec:

  providers:
    - name: storage-provider
      providerType: hetzner
      secretRef:
        name: storage-provider-secrets
        namespace: claudie-secrets

    - name: compute-provider
      providerType: hetzner
      secretRef:
        name: storage-provider-secrets
        namespace: claudie-secrets

    - name: dns-provider
      providerType: cloudflare
      secretRef:
        name: dns-provider-secret
        namespace: claudie-secrets

  nodePools:
    dynamic:
        - name: control
          providerSpec:
            name: compute-provider
            region: hel1
            zone: hel1-dc2
          count: 3
          serverType: cpx21
          image: ubuntu-22.04

        - name: datastore
          providerSpec:
            name: storage-provider
            region: hel1
            zone: hel1-dc2
          count: 5
          serverType: cpx21
          image: ubuntu-22.04
          storageDiskSize: 800
          taints:
            - key: node-type
              value: datastore
              effect: NoSchedule
  
        - name: compute
          providerSpec:
            name: compute-provider
            region: hel1
            zone: hel1-dc2
          count: 10
          serverType: cpx41
          image: ubuntu-22.04
          taints:
            - key: node-type
              value: compute
              effect: NoSchedule

        - name: loadbalancer
          providerSpec:
            name: compute-provider
            region: hel1
            zone: hel1-dc2
          count: 1
          serverType: cpx21
          image: ubuntu-22.04

  kubernetes:
    clusters:
      - name: my-awesome-claudie-cluster
        version: 1.27.0
        network: 192.168.2.0/24
        pools:
          control:
            - control
          compute:
            - datastore
            - compute

  loadBalancers:
    roles:
      - name: apiserver
        protocol: tcp
        port: 6443
        targetPort: 6443
        targetPools: 
          - control

    clusters:
      - name: apiserver-lb
        roles:
          - apiserver
        dns:
          dnsZone: dns-zone
          provider: dns-provider
        targetedK8s: my-awesome-claudie-cluster
        pools:
          - loadbalancer
```

When Claudie applies this input manifest, the following storage classes are installed:

- `longhorn` - the default storage class, which stores data on random nodes
- `longhorn-storage-provider-zone` - storage class, which stores data only on nodes of the `storage-provider` provider instance.
- `longhorn-compute-provider-zone` - storage class, which stores data only on nodes of the `compute-provider` provider instance.

Now all you have to do is specify correct storage class when defining your PVCs.

In case you are interested in using different cloud provider for `datastore-nodepool` or `compute-nodepool` of this `InputManifest` example, see the [list of supported providers instance](../getting-started/detailed-guide.md#supported-providers)

For more information on how Longhorn works you can check out [Longhorn's official documentation](https://longhorn.io/docs/1.4.0/what-is-longhorn/).
