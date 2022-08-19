---
subcategory: "MaxCompute"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_maxcompute_user"
sidebar_current: "docs-apsarastack-resource-maxcompute-user"
description: |-
  Provides a Apsarastack maxcompute user resource.
---

# apsarastack\_maxcompute\_user

The user is the basic unit of operation in maxcompute. It is similar to the concept of Database or Schema in traditional databases, and sets the boundary for maxcompute multi-user isolation and access control.
->**NOTE:** Available in 1.0.18+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_maxcompute_user" "example" {
  user_name             = "%s"
  description           = "TestAccApsaraStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
```
## Argument Reference

The following arguments are supported:
* `user_name` - (Required, ForceNew) The name of the user that you want to create.
* `description` - (Required, ForceNew) The description of the user that you want to create.
* `organization_id` - (Optional) The id of the organization. 

## Attributes Reference



## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import apsarastack_maxcompute_cu.example tf_maxcompute_cu
```
