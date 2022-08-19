---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ecs_hpc_cluster"
sidebar_current: "docs-apsarastack-resource-ecs-hpc-cluster"
description: |-
  Provides a Apsarastack ECS Hpc Cluster resource.
---

# apsarastack\_ecs\_hpc\_cluster

Provides a ECS Hpc Cluster resource.

For information about ECS Hpc Cluster and how to use it, see [What is Hpc Cluster](https://www.alibabacloud.com/help/en/doc-detail/109138.htm).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_ecs_hpc_cluster" "example" {
  name        = "tf-testAcc"
  description = "For Terraform Test"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of ECS Hpc Cluster.
* `name` - (Required) The name of ECS Hpc Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hpc Cluster.

## Import

ECS Hpc Cluster can be imported using the id, e.g.

```
$ terraform import apsarastack_ecs_hpc_cluster.example <id>
```
