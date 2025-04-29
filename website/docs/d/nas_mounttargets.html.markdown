---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_mounttargets"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-mounttargets"
description: |- 
  Provides a list of nas mounttargets owned by an alibabacloudstack account.
---

# alibabacloudstack_nas_mounttargets
-> **NOTE:** Alias name has: `alibabacloudstack_nas_mount_targets`

This data source provides a list of NAS MountTargets in an AlibabacloudStack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_nas_mounttargets" "example" {
  file_system_id    = "1a2sc4d"
  access_group_name = "tf-testAccNasConfig"
  network_type      = "VPC"
  vpc_id            = "vpc-abc1234567890"
  vswitch_id        = "vsw-abc1234567890"
  status            = "Active"
}

output "first_mount_target_domain" {
  value = data.alibabacloudstack_nas_mounttargets.example.targets.0.mount_target_domain
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Optional, ForceNew) The name of the permission group.
* `mount_target_domain` - (Optional, Deprecated in 1.53.+) The domain name of the Mount point. This field has been deprecated from provider version 1.53.0. Use `ids` instead.
* `type` - (Optional, Deprecated in 1.95.0+) Field `type` has been deprecated from provider version 1.95.0. Use `network_type` instead.
* `network_type` - (Optional, ForceNew, Available 1.95.0+) Network type. Valid values include `VPC`, etc.
* `vpc_id` - (Optional, ForceNew) VPC ID.
* `vswitch_id` - (Optional, ForceNew) VSwitch ID.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `ids` - (Optional, ForceNew, Available 1.53.0+) A list of MountTargetDomain IDs.
* `status` - (Optional, ForceNew, Available 1.95.0+) The current status of the Mount point. Valid values: `Active`, `Inactive`, and `Pending`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `targets` - A list of NAS MountTargets. Each element contains the following attributes:
  * `access_group_name` - The name of the permission group.
  * `id` - The unique identifier of the MountTarget.
  * `mount_target_domain` - The domain name of the Mount point.
  * `network_type` - Network type.
  * `type` - (Deprecated in 1.95.0+) This field has been deprecated from provider version 1.95.0. Use `network_type` instead.
  * `status` - The current status of the Mount point. Valid values: `Active`, `Inactive`, and `Pending`.
  * `vpc_id` - VPC ID.
  * `vswitch_id` - VSwitch ID.