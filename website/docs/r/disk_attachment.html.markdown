---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_disk_attachment"
sidebar_current: "docs-alibabacloudstack-resource-disk-attachment"
description: |-
  Provides a ECS Disk Attachment resource.
---

# apsaratack\_disk\_attachment

Provides an alibabacloudstack ECS Disk Attachment as a resource, to attach and detach disks from ECS Instances.

## Example Usage

Basic usage

```
# Create a new ECS disk-attachment and use it attach one disk to a new instance.

resource "alibabacloudstack_disk" "ecs_disk" {
  availability_zone = "${var.availability_zone}"
  size              = "50"

  tags = {
    Name = "TerraformTest-disk"
  }
}

resource "alibabacloudstack_instance" "instance" {
  image_id              = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  system_disk_name     = "test_sys_disk"
  security_groups      = [var.security_group_id]
  instance_name        = "test_apsara_instance"
  vswitch_id           = var.vswitch_id
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_disk_attachment" "ecs_disk_att" {
  disk_id     = "${alibabacloudstack_disk.ecs_disk.id}"
  instance_id = "${alibabacloudstack_instance.instance.id}"
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, Forces new resource) ID of the Instance to attach to.
* `disk_id` - (Required, Forces new resource) ID of the Disk to be attached.
* `device_name` - (Optional, Forces new resource) The device name which is used when attaching disk, it will be allocated automatically by system according to default order from /dev/xvdb to /dev/xvdz.
                                                          
## Attributes Reference

The following attributes are exported:

* `instance_id` - ID of the Instance.
* `disk_id` - ID of the Disk.
* `device_name` - The device name exposed to the instance.
