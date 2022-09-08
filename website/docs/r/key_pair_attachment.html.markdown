---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_key_pair_attachment"
sidebar_current: "docs-alibabacloudstack-resource-key-pair-attachment"
description: |-
  Provides a AlibabacloudStack key pair attachment resource to bind key pair for several ECS instances.
---

# alibabacloudstack\_key\_pair\_attachment

Provides a key pair attachment resource to bind key pair for several ECS instances.

-> **NOTE:** After the key pair is attached with some instances, there instances must be rebooted to make the key pair affect.

## Example Usage

Basic Usage

```
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_ssd"
  available_resource_creation = "VSwitch"
}
data "alibabacloudstack_instance_types" "type" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "alibabacloudstack_images" "images" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
variable "name" {
  default = "keyPairAttachmentName"
}

resource "alibabacloudstack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "vswitch" {
  vpc_id            = "${alibabacloudstack_vpc.vpc.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "group" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_instance" "instance" {
  instance_name   = "${var.name}-${count.index + 1}"
  image_id        = "${data.alibabacloudstack_images.images.images.0.id}"
  instance_type   = "${data.alibabacloudstack_instance_types.type.instance_types.0.id}"
  count           = 2
  security_groups = ["${alibabacloudstack_security_group.group.id}"]
  vswitch_id      = "${alibabacloudstack_vswitch.vswitch.id}"
  internet_max_bandwidth_out = 5
  password                   = "Test12345"
  system_disk_category = "cloud_ssd"
}

resource "alibabacloudstack_key_pair" "pair" {
  key_name = "${var.name}"
}

resource "alibabacloudstack_key_pair_attachment" "attachment" {
  key_name     = "${alibabacloudstack_key_pair.pair.id}"
  instance_ids = ["${alibabacloudstack_instance.instance.*.id}"]
}
```
## Argument Reference

The following arguments are supported:

* `key_name` - (Required, ForceNew) The name of key pair used to bind.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.
* `force` - (ForceNew) Set it to true and it will reboot instances which attached with the key pair to make key pair affect immediately.

## Attributes Reference

* `key_name` - The name of the key pair.
* `instance_ids` The list of ECS instance's IDs.
