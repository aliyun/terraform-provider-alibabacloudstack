---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessrule"
sidebar_current: "docs-Alibabacloudstack-nas-accessrule"
description: |- 
  Provides a nas Accessrule resource.
---

# alibabacloudstack_nas_accessrule
-> **NOTE:** Alias name has: `alibabacloudstack_nas_access_rule`

Provides a nas Accessrule resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
    default = "tf-testaccnasaccess_rule99652"
}

resource "alibabacloudstack_nas_access_group" "example" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
}

resource "alibabacloudstack_nas_accessrule" "default" {
  access_group_name = alibabacloudstack_nas_access_group.example.access_group_name
  source_cidr_ip    = "1.1.1.1/0"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 1
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) The name of the permission group. This parameter is required and cannot be modified once set.
* `source_cidr_ip` - (Required) The address or address segment that you want to allow access to the NAS file system. For example, `1.1.1.1/0`.
* `rw_access_type` - (Optional) The read-write permission type for the rule. Valid values are:
  * `RDWR`: Read and write access (default).
  * `RDONLY`: Read-only access.
* `user_access_type` - (Optional) The user permission type for the rule. Valid values are:
  * `no_squash`: No restrictions on root users (default).
  * `root_squash`: Restricts root users from having full privileges.
  * `all_squash`: Restricts all users from having full privileges.
* `priority` - (Optional) The priority level of the rule. Valid range is 1-100. Lower numbers indicate higher priority. Default value is `1`.
* `access_rule_id` - (Optional) - The unique identifier for the NAS access rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of this resource. It is formatted as `<access_group_name>:<access_rule_id>`.
* `access_rule_id` - The unique identifier for the NAS access rule.
