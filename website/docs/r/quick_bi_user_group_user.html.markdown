---
subcategory: "QuickBI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user_group_user"
sidebar_current: "docs-Alibabacloudstack-resource-quick-bi-user-group-user"
description: |-
  Provides a Alibabacloudstack Quick BI UserGroupUser resource.
---

# alibabacloudstack_quick_bi_user_group_user

Provides a Quick BI UserGroupUser resource.


## Example Usage

Basic Usage

```terraform

resource "alibabacloudstack_quick_bi_user_group_user" "example" {
  user_group_id = alibabacloudstack_quick_bi_user_group.default.id
  account_id    = alibabacloudstack_quick_bi_user.default.id
}

resource "alibabacloudstack_quick_bi_user" "default" {
  nick_name       = "example_value"
  account_name    = "example_value"
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}

resource "alibabacloudstack_quick_bi_user_group" "default" {
  user_group_name        = "example_value"
  user_group_description = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `user_group_id` - (Required) User group ID. Modifying this parameter will force a new resource to be created.
* `account_id` - (Required) Quick BI user ID. Modifying this parameter will force a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format of `<user_group_id>:<account_id>`.

## Import

Quick BI UserGroupUser can be imported using the user_group_id and account_id, e.g.

```
$ terraform import alibabacloudstack_quick_bi_user_group_user.example <user_group_id>:<account_id>
```
