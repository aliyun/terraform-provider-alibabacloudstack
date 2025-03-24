---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-keypairattachment"
description: |- 
  Provides a ecs Keypairattachment resource.
---

# alibabacloudstack_ecs_keypairattachment
-> **NOTE:** Alias name has: `alibabacloudstack_key_pair_attachment`

Provides a ecs Keypairattachment resource.

## Example Usage

Basic Usage

```hcl
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

variable "password" {
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
  password                   = var.password
  system_disk_category = "cloud_ssd"
}

resource "alibabacloudstack_key_pair" "pair" {
  key_name = "${var.name}"
}

resource "alibabacloudstack_ecs_keypairattachment" "attachment" {
  key_name     = "${alibabacloudstack_key_pair.pair.key_name}"
  instance_ids = ["${alibabacloudstack_instance.instance.*.id}"]
  force        = true
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required, ForceNew) The name of the key pair used to bind.
* `instance_ids` - (Required, ForceNew) The list of ECS instance IDs to which the key pair will be attached.
* `force` - (Optional, ForceNew) If set to `true`, the instances will be rebooted immediately after attaching the key pair to ensure that the key pair takes effect without requiring manual intervention.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_name` - The name of the key pair that has been attached.
* `instance_ids` - The list of ECS instance IDs to which the key pair is attached.