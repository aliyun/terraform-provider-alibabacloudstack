---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpc_cluster"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-hpc-cluster"
description: |- 
  Provides a Alibabacloudstack ECS Hpc Cluster resource.
---

# alibabacloudstack_ecs_hpc_cluster
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_hpccluster`

Provides a ECS Hpc Cluster resource.

For information about ECS Hpc Cluster and how to use it, see [What is Hpc Cluster](https://www.alibabacloud.com/help/en/doc-detail/109138.htm).

## Example Usage

Basic Usage

```terraform
variable "name" {
    default = "tf-testaccecshpc_cluster26697"
}

resource "alibabacloudstack_ecs_hpc_cluster" "example" {
  name        = var.name
  description = "For Terraform Test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HPC cluster. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter or Chinese character but cannot start with `http://` or `https://`. This name must be unique within the same region.
* `description` - (Optional) The description of the HPC cluster. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`. It can contain uppercase/lowercase letters, numbers, periods (.), colons (:), underscores (_), hyphens (-), and at symbols (@). Default value: empty.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the HPC cluster.
* `hpc_cluster_id` - (Computed) The ID of the HPC cluster.
* `name` - (Computed) The name of the HPC cluster.
* `description` - (Computed) The description of the HPC cluster.

## Import

ECS Hpc Cluster can be imported using the HpcClusterId, e.g.

```
$ terraform import alibabacloudstack_ecs_hpc_cluster.example hpc-bp1a5zr3u7nq9cx****
```