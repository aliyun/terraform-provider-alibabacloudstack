---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessgroup"
sidebar_current: "docs-Alibabacloudstack-nas-accessgroup"
description: |- 
  Provides a nas Accessgroup resource.
---

# alibabacloudstack_nas_accessgroup
-> **NOTE:** Alias name has: `alibabacloudstack_nas_access_group`

Provides a nas Accessgroup resource.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_nas_accessgroup" "foo" {
  access_group_name = "CreateAccessGroup"
  access_group_type = "Vpc"
  description       = "test_AccessG"
  file_system_type = "extreme"
}
```

Example with all parameters

```terraform
variable "name" {
  default = "tf-testaccnasaccess_group31001"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
  access_group_name = "accssGroupExtremeVpcTest"
  file_system_type  = "extreme"
  access_group_type = "Vpc"
  description       = "test"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) The name of the permission group. Once set, it cannot be modified.
* `file_system_type` - (Optional, ForceNew) The type of file system. Valid values: `standard` and `extreme`. Default to `standard`. Note that the `extreme` type only supports the `Vpc` network.
* `access_group_type` - (Required, ForceNew) Permission group types. Valid values: `Vpc` and `Classic`.
* `description` - (Optional) Permission group description information. This provides additional details about the permission group for better identification.

### Deprecated Arguments

The following arguments are deprecated and replaced in version 1.92.0:

* `name` - (Deprecated) Replaced by `access_group_name` after version 1.92.0.
* `type` - (Deprecated) Replaced by `access_group_type` after version 1.92.0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Access Group. The value is formatted as `<access_group_name>:<file_system_type>` after version 1.92.0.
* `access_group_name` - The name of the permission group.
* `access_group_type` - Permission group types, including `Vpc` and `Classic`.
* `file_system_type` - File system type. Valid values: `standard` and `extreme`.

## Import

NAS Access Group can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_nas_accessgroup.foo tf_testAccNasConfig:standard
```