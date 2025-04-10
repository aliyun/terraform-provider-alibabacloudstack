subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_mounttarget"
sidebar_current: "docs-Alibabacloudstack-nas-mounttarget"
description: |- 
  Provides a nas Mounttarget resource.

---

# alibabacloudstack_nas_mounttarget
-> **NOTE:** Alias name has: `alibabacloudstack_nas_mount_target`

Provides a NAS Mount Target resource. For information about NAS Mount Target and how to use it, see [Manage NAS Mount Targets](https://www.alibabacloud.com/help/en/doc-detail/27531.htm).

-> **NOTE**: Available in v1.34.0+.

-> **NOTE**: Currently this resource supports creating a mount point in a classic network only when the current region is China mainland regions.

-> **NOTE**: You must grant NAS with specific RAM permissions when creating a classic mount target, and it only can be achieved by creating a classic mount target manually. See [Add a mount point](https://www.alibabacloud.com/help/doc-detail/60431.htm) and [Why do I need RAM permissions to create a mount point in a classic network](https://www.alibabacloud.com/help/faq-detail/42176.htm).

## Example Usage

```terraform
variable "name" {
    default = "tf-testaccnasmount_target81413"
}

resource "alibabacloudstack_vpc" "default" {
    cidr_block = "172.16.0.0/16"
    name       = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
    available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "172.16.0.0/21"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name             = "${var.name}"
}

variable "storage_type" {
    default = "Capacity"
}

data "alibabacloudstack_nas_protocols" "default" {
    type = "${var.storage_type}"
}

resource "alibabacloudstack_nas_file_system" "default" {
    description  = "${var.name}"
    storage_type = "${var.storage_type}"
    protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}

resource "alibabacloudstack_nas_access_group" "default" {
    access_group_name = "tf-testAccNasConfig-resource-test86663"
    access_group_type = "Vpc"
    description       = "tf-testAccNasConfig"
}

resource "alibabacloudstack_nas_mounttarget" "default" {
    file_system_id    = "${alibabacloudstack_nas_file_system.default.id}"
    vswitch_id        = "${alibabacloudstack_vswitch.default.id}"
    access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
    status            = "Active"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Optional) <!--  AI CREATE  --> The name of the permission group that applies to the mount target.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch in the VPC where the mount target resides.
* `status` - (Optional) The status of the mount target. Valid values: `Active` and `Inactive`. Default value is `Active`. Before you mount a file system, make sure that the mount target is in the Active state.
* `security_group_id` - (Optional, ForceNew, Available in v1.95.0+) The ID of the security group.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - This ID of this resource. It is formatted to `<file_system_id>:<mount_target_domain>`. Before version 1.95.0, the value is `<mount_target_domain>`.
* `status` - (Computed) <!--  AI CREATE  --> The current status of the Mount point, including `Active` and `Inactive`. You can use the file system only when the status is `Active`.