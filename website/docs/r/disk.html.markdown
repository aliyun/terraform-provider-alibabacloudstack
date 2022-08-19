---
subcategory: "ECS"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_disk"
sidebar_current: "docs-apsarastack-resource-disk"
description: |-
  Provides a ECS Disk resource.
---

# apsarastack\_disk

Provides a ECS disk resource.

-> **NOTE:** One of `size` or `snapshot_id` is required when specifying an ECS disk. If all of them be specified, `size` must more than the size of snapshot which `snapshot_id` represents. Currently, `apsarastack_disk` doesn't resize disk.

## Example Usage

```
# Create a new ECS disk.
resource "apsarastack_disk" "ecs_disk" { 
  availability_zone = "${var.availability_zone}"
  name              = "New-disk"
  description       = "ECS-Disk"
  category          = "cloud_efficiency"
  size              = "30"

  tags = {
    Name = "TerraformTest"
  }
}
```
## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The Zone to create the disk in.
* `name` - (Optional) Name of the ECS disk. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the disk. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `category` - (Optional, ForceNew) Category of the disk. Valid values are `cloud`, `cloud_efficiency`, `cloud_ssd`. Default is `cloud`.
* `size` - (Optional) The size of the disk in GiBs. When resize the disk, the new size must be greater than the former value, or you would get an error `InvalidDiskSize.TooSmall`.
* `snapshot_id` - (Optional) A snapshot to base the disk off of. If the disk size required by a snapshot is greater than `size`, the `size` will be ignored, conflict with `encrypted`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `encrypted` - (Optional) If true, the disk will be encrypted, conflict with `snapshot_id`.
* `delete_auto_snapshot` - (Optional) Indicates whether the automatic snapshot is deleted when the disk is released. Default value: false.
* `delete_with_instance` - (Optional) Indicates whether the disk is released together with the instance: Default value: false.
* `enable_auto_snapshot` - (Optional) Indicates whether to apply a created automatic snapshot policy to the disk. Default value: false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk.
* `status` - The disk status.
