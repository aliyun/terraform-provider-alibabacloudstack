---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_user"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-user"
description: |-
  Provides a Alibabacloudstack maxcompute user resource.
---

# alibabacloudstack_maxcompute_user

The user is the basic unit of operation in maxcompute. It is similar to the concept of Database or Schema in traditional databases, and sets the boundary for maxcompute multi-user isolation and access control.
->**NOTE:** Available in 1.0.18+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_maxcompute_user" "example" {
  user_name             = "%s"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
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
* `organization_name` - (Optional) The name of the organization. 

## Attributes Reference

* `id` - The ID of the user.
* `user_id` - Alias of the key `id`.
* `user_pk` - The PK of the user.
* `user_type` - The type of the user.
* `user_name` - (Computed) The name of the user.
* `description` - (Computed) The description of the user.
* `organization_name` - The name of the organization. 

## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import alibabacloudstack_maxcompute_cu.example tf_maxcompute_cu
```