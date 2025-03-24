---
subcategory: "ACK"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ack-clusters"
description: |- 
  Provides a list of ack clusters owned by an AlibabacloudStack account.
---

# alibabacloudstack_ack_clusters
-> **NOTE:** Alias name has: `alibabacloudstack_cs_kubernetes_clusters`

This data source provides a list of ACK clusters in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_ack_clusters" "example" {
  name_regex = "my-first-ack"
}

output "ack_cluster_ids" {
  value = data.alibabacloudstack_ack_clusters.example.ids
}

output "ack_cluster_names" {
  value = data.alibabacloudstack_ack_clusters.example.names
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Cluster IDs to filter results. If not specified, all clusters will be considered.
* `name_regex` - (Optional) A regex string to filter results by cluster name.
* `state` - (Optional) Filter results by cluster state. Valid values include `Running`, `Creating`, `Updating`, etc.
* `enable_details` - (Optional) Boolean, defaults to `false`. Setting this parameter to `true` will return more details about each cluster, such as `master_nodes`, `worker_nodes`, and `connections`.
* `kube_config` - (Optional) Boolean, set to `true` if you want to obtain kubeconfig for the clusters specified in `ids`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of matched ACK clusters.
* `ids` - A list of IDs of matched ACK clusters.
* `clusters` - A list of matched ACK clusters. Each element contains the following attributes:
  * `id` - The ID of the ACK cluster.
  * `name` - The name of the ACK cluster.
  * `availability_zone` - The availability zone where the cluster is located.
  * `slb_internet_enabled` - Indicates whether an internet-facing load balancer for the API server is created.
  * `security_group_id` - The security group ID associated with the worker nodes of the cluster.
  * `nat_gateway_id` - The NAT gateway ID used for launching the Kubernetes cluster.
  * `vpc_id` - The VPC ID where the cluster is located.
  * `vswitch_ids` - A list of VSwitch IDs where the cluster is located.
  * `master_instance_types` - A list of instance types for the master nodes.
  * `worker_instance_types` - A list of instance types for the worker nodes.
  * `worker_numbers` - The number of worker nodes in the cluster.
  * `pod_cidr` - The CIDR block for the pod network when using Flannel.
  * `cluster_network_type` - The network type used by the cluster, such as `flannel` or `terway`.
  * `node_cidr_mask` - The network mask used on pods for each node.
  * `log_config` - A list of one element containing information about the associated log store. It includes:
    * `type` - Type of collecting logs.
    * `project` - Log Service project name.
  * `image_id` - The image ID used for the nodes.
  * `master_disk_size` - The system disk size of the master nodes.
  * `state` - The current state of the cluster.
  * `master_disk_category` - The system disk category of the master nodes.
  * `worker_disk_size` - The system disk size of the worker nodes.
  * `worker_disk_category` - The system disk category of the worker nodes.
  * `master_nodes` - A list of master nodes in the cluster. Each element includes:
    * `id` - The ID of the master node.
    * `name` - The name of the master node.
    * `private_ip` - The private IP address of the master node.
  * `worker_nodes` - A list of worker nodes in the cluster. Each element includes:
    * `id` - The ID of the worker node.
    * `name` - The name of the worker node.
    * `private_ip` - The private IP address of the worker node.
  * `connections` - A map of connection information for the Kubernetes cluster. It includes:
    * `api_server_internet` - The internet endpoint for the API server.
    * `api_server_intranet` - The intranet endpoint for the API server.
    * `master_public_ip` - The public IP address for SSH access to the master node.
    * `service_domain` - The domain for accessing services within the cluster.