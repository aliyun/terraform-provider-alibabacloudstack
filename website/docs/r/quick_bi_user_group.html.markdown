---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user_group"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-user-group"
description: |-
  Provides a Alibabacloudstack Quick BI UserGroup resource.
---

# alibabacloudstack_quick_bi_user_group

Provides a Quick BI UserGroup resource.


## Example Usage

Basic Usage

```terraform

resource "alibabacloudstack_quick_bi_user_group" "example" {
  user_group_name = "example_value"
  user_group_description = "example_value"
  parent_user_group_id = "-1"
}

```

## Argument Reference

The following arguments are supported:

* `user_group_name` - (Required) User group name.
* `user_group_description` - (Required) User group description.
* `parent_user_group_id` - (Required) Parent user group ID. You can add a new user group to this grouping.When you enter -1, the newly created user group will be added to the root directory.
* `user_group_id` - (Optional)  User group ID.

## Attributes Reference

The following attributes are exported:

* `user_group_id` -  The resource ID in terraform of UserGroup.

## Import

Quick BI UserGroup can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_quick_bi_user_group.example <id>
```