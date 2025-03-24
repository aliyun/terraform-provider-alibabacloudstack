---
subcategory: "Container Service (CS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_kubernetes_clusters"
sidebar_current: "docs-alibabacloudstack-datasource-cs-kubernetes-clusters"
description: |-
  Provides a list of Container Service Kubernetes Clusters to be used by the alibabacloudstack_cs_kubernetes_cluster resource.
---

# alibabacloudstack_cs_kubernetes_clusters
-> **NOTE:** Alias name has: `alibabacloudstack_ack_clusters`

This data source provides a list Container Service Kubernetes Clusters on AlibabacloudStack.


## Example Usage

```
# Declare the data source
data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
  cluster_id = "xx"
}
output "kubeconfig" {
  value = data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig.kubeconfig
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) Cluster ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `kubeconfig` - The kubeconfig content for the Kubernetes cluster.

Clusters
name: The name of the cluster, referenced in contexts.
cluster:
server: The URL of the Kubernetes API server.
certificate-authority: The path to the CA certificate, used to verify the identity of the API server.
insecure-skip-tls-verify: If set to true, it skips TLS verification (not recommended).
Users
name: The username, referenced in contexts.
user:
client-certificate: The path to the client certificate, used to authenticate with the API server.
client-key: The path to the private key corresponding to the client certificate.
token: If using token-based authentication, this is where you provide the token string.
username/password: For some simple authentication mechanisms, a username and password may be required.
Contexts
name: The context name, which can be named arbitrarily.
context:
cluster: The associated cluster name.
user: The user to use for accessing the cluster.
namespace: The default namespace.