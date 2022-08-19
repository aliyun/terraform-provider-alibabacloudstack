---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_key_pair_attachment"
sidebar_current: "docs-apsarastack-resource-key-pair-attachment"
description: |-
  Provides a ApsaraStack key pair attachment resource to bind key pair for several ECS instances.
---

# apsarastack\_key\_pair\_attachment

Provides a key pair attachment resource to bind key pair for several ECS instances.

-> **NOTE:** After the key pair is attached with some instances, there instances must be rebooted to make the key pair affect.

## Example Usage

Basic Usage

```
data "apsarastack_zones" "default" {
  available_disk_category     = "cloud_ssd"
  available_resource_creation = "VSwitch"
}
data "apsarastack_instance_types" "type" {
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "apsarastack_images" "images" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
variable "name" {
  default = "keyPairAttachmentName"
}

resource "apsarastack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "vswitch" {
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "group" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = "${apsarastack_vpc.vpc.id}"
}

resource "apsarastack_instance" "instance" {
  instance_name   = "${var.name}-${count.index + 1}"
  image_id        = "${data.apsarastack_images.images.images.0.id}"
  instance_type   = "${data.apsarastack_instance_types.type.instance_types.0.id}"
  count           = 2
  security_groups = ["${apsarastack_security_group.group.id}"]
  vswitch_id      = "${apsarastack_vswitch.vswitch.id}"
  internet_max_bandwidth_out = 5
  password                   = "Test12345"
  system_disk_category = "cloud_ssd"
}

resource "apsarastack_key_pair" "pair" {
  key_name = "${var.name}"
}

resource "apsarastack_key_pair_attachment" "attachment" {
  key_name     = "${apsarastack_key_pair.pair.id}"
  instance_ids = ["${apsarastack_instance.instance.*.id}"]
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
