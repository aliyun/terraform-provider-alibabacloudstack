---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-user"
description: |-
  Provides a Alibabacloudstack Quick BI User resource.
---

# alibabacloudstack\_quick\_bi\_user

Provides a Quick BI User resource.

For information about Quick BI User and how to use it, see [What is User](https://www.alibabacloud.com/help/doc-detail/33813.htm).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_quick_bi_user" "example" {
  account_name    = "example_value"
  admin_user      = false
  auth_admin_user = false
  nick_name       = "example_value"
  user_type       = "Analyst"
}

```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional, ForceNew) Alibaba Cloud account ID.
* `account_name` - (Required) An Alibaba Cloud account, Alibaba Cloud name.
* `admin_user` - (Required) Whether it is the administrator. Valid values: `true` and `false`.
* `auth_admin_user` - (Required) Whether this is a permissions administrator. Valid values: `false`, `true`.
* `nick_name` - (Required, ForceNew) The nickname of the user.
* `user_type` - (Required) The members of the organization of the type of role separately. Valid values: `Analyst`, `Developer` and `Visitor`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User.

## Import

Quick BI User can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_quick_bi_user.example <id>
```
