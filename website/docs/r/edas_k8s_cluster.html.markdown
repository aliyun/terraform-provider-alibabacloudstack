---
subcategory: "Enterprise Distributed Application Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_cluster"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-cluster"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alibabacloudstack_edas_k8s_cluster

Provides an EDAS K8s cluster resource. For information about EDAS K8s Cluster and how to use it, see [What is EDAS K8s Cluster](https://www.alibabacloud.com/help/en/doc-detail/85108.htm).



## Example Usage

Basic Usage

```
resource "alibabacloudstack_edas_k8s_cluster" "default" {
  cs_cluster_id = "xxxx-xxx-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `cs_cluster_id` - (Required, ForceNew) The ID of the Container Service Kubernetes cluster that you want to import. You can call the [GetK8sCluster](https://www.alibabacloud.com/help/en/doc-detail/85108.htm) operation to query the cluster ID.
* `namespace_id` - (Optional, ForceNew) The ID of the namespace where you want to import. You can call the [ListUserDefineRegion](https://www.alibabacloud.com/help/en/doc-detail/149377.htm) operation to query the namespace ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EDAS K8s cluster.
* `cluster_name` - The name of the cluster.
* `cluster_type` - The type of the cluster. Valid values: `5`: Container Service K8s cluster or Serverless K8s cluster.
* `network_mode` - The network type of the cluster. Valid values: `1`: Classic network. `2`: VPC.
* `vpc_id` - The ID of the Virtual Private Cloud (VPC) for the cluster.
* `cluster_import_status` - The import status of the cluster. Valid values:
    * `1`: Success.
    * `2`: Failed.
    * `3`: Importing.
    * `4`: Deleted.
* `cs_cluster_id` - The ID of the Container Service Kubernetes cluster that you want to import.
* `namespace_id` - The ID of the namespace where you want to import.

## Import

EDAS K8s Cluster can be imported using the cluster ID (the EDAS internal cluster ID returned after import), e.g.

```
$ terraform import alibabacloudstack_edas_k8s_cluster.example 81453e4b-4df0-4592-xxxx-b835a2eexxxx
```