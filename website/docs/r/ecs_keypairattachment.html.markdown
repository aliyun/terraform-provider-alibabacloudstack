---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-keypairattachment"
description: |-
  Provides a ECS Key Pair Attachment resource.
---

# alibabacloudstack_ecs_keypairattachment

-> **NOTE:** This resource can also be referred to by the following alias: `alibabacloudstack_key_pair_attachment`.

Provides a ECS Key Pair Attachment resource to bind an SSH key pair to one or more Linux instances.

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

* `key_name` - (Required, ForceNew) The name of the SSH key pair to bind. The name must be 2 to 128 characters in length. Changing this parameter forces a new resource to be created.
* `instance_ids` - (Required, ForceNew) A list of ECS instance IDs to which the SSH key pair will be bound. Up to 50 instance IDs can be specified in a JSON array format. Changing this parameter forces a new resource to be created.
* `force` - (Optional, ForceNew) If set to `true`, the instances will be rebooted automatically after binding the key pair to ensure it takes effect immediately. Default to `false`. Changing this parameter forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, formatted as `<key_name>:<instance_ids>` where `instance_ids` is a JSON array string.
* `key_name` - The name of the SSH key pair that has been bound.
* `instance_ids` - A list of ECS instance IDs to which the SSH key pair is bound.

## Import

ECS Key Pair Attachment can be imported using the key pair name and instance IDs in the format `<key_name>:<instance_ids>`, e.g.

```
$ terraform import alibabacloudstack_ecs_keypairattachment.example test-key-pair:["i-bp1d6tsvznfghy7y****","i-bp1ippxbaql9zet7****"]
```