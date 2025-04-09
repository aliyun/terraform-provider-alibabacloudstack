---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpc_cluster"
sidebar_current: "docs-alibabacloudstack-resource-ecs-hpc-cluster"
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

* `name` - (Required) The name of the ECS Hpc Cluster. This name must be unique within the same region and can consist of uppercase/lowercase letters, numbers, hyphens (-), and underscores (_). It cannot exceed 128 characters in length.
* `description` - (Optional) The description of the ECS Hpc Cluster. This description can consist of uppercase/lowercase letters, numbers, periods (.), colons (:), underscores (_), hyphens (-), and at symbols (@). It cannot exceed 256 characters in length.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `description` - (Computed) The description of the ECS Hpc Cluster.

* `id` - The ID of the ECS Hpc Cluster. This attribute is the same as the `name` argument and can be used for referencing this resource in other parts of your Terraform configuration.

## Import

ECS Hpc Cluster can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_ecs_hpc_cluster.example <id>
```