---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ess-attachment"
description: |- 
  Provides an ESS Attachment resource to attach or remove ECS instances.
---

# alibabacloudstack_ess_attachment

Attaches several ECS instances to a specified scaling group or remove them from it.

-> **NOTE:** ECS instances can be attached or remove only when the scaling group is active and it has no scaling
activity in progress.

-> **NOTE:** There are two types ECS instances in a scaling group: "AutoCreated" and "Attached". The total number of
them can not larger than the scaling group "MaxSize".

## Example Usage


```
variable "name" {
  default = "essattachmentconfig"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vswitch.default.id]
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id  = alibabacloudstack_ess_scaling_group.default.id
  image_id          = data.alibabacloudstack_images.default.images[0].id
  instance_type     = data.alibabacloudstack_instance_types.default.instance_types[0].id
  security_group_id = alibabacloudstack_security_group.default.id
  force_delete      = true
  active            = true
  enable            = true
}

resource "alibabacloudstack_instance" "default" {
  image_id                   = data.alibabacloudstack_images.default.images[0].id
  instance_type              = data.alibabacloudstack_instance_types.default.instance_types[0].id
  count                      = 2
  security_groups            = [alibabacloudstack_security_group.default.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alibabacloudstack_vswitch.default.id
  instance_name              = var.name
}

resource "alibabacloudstack_ess_attachment" "default" {
  scaling_group_id = alibabacloudstack_ess_scaling_group.default.id
  instance_ids     = [alibabacloudstack_instance.default[0].id, alibabacloudstack_instance.default[1].id]
  force            = true
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required) ID of the scaling group of a scaling configuration. 
* `instance_ids` - (Required) ID of the ECS instance to be attached to the scaling group. You can input up to 20 IDs.
* `force` - (Optional) Whether to remove forcibly "AutoCreated" ECS instances in order to release scaling group
  capacity "MaxSize" for attaching ECS instances. Default to false.

-> **NOTE:** "AutoCreated" ECS instance will be deleted after it is removed from scaling group, but "Attached" will be
not.

-> **NOTE:** Restrictions on attaching ECS instances:

- The attached ECS instances and the scaling group must have the same region and network type(`Classic` or `VPC`).
- The attached ECS instances and the instance with active scaling configurations must have the same instance type.
- The attached ECS instances must in the running state.
- The attached ECS instances has not been attached to other scaling groups.
- The attached ECS instances supports Subscription and Pay-As-You-Go payment methods.

## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS attachment resource ID.
* `instance_ids` - (Required)ID of list "Attached" ECS instance.
* `force` - Whether to delete "AutoCreated" ECS instances.

## Import

ESS attachment can be imported using the id or scaling group id, e.g.

```
$ terraform import alibabacloudstack_ess_attachment.example asg-abc123456
```