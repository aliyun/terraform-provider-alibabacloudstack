---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterfaceattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-networkinterfaceattachment"
description: |- 
  Provides a ecs Networkinterfaceattachment resource.
---

# alibabacloudstack_ecs_networkinterfaceattachment
-> **NOTE:** Alias name has: `alibabacloudstack_network_interface_attachment`

Provides a ecs Networkinterfaceattachment resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

data "alibabacloudstack_instance_types" "eni2" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount       = 2
  sorted_by        = "Memory"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "default" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  availability_zone = reverse(data.alibabacloudstack_zones.default.zones)[0].id
  vpc_id            = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  availability_zone         = reverse(data.alibabacloudstack_zones.default.zones)[0].id
  security_groups          = [alibabacloudstack_security_group.default.id]
  instance_type            = data.alibabacloudstack_instance_types.eni2.instance_types[0].id
  system_disk_category     = "cloud_efficiency"
  image_id                 = data.alibabacloudstack_images.default.images[0].id
  instance_name            = var.name
  vswitch_id               = alibabacloudstack_vswitch.default.id
  internet_max_bandwidth_out = 10
}

resource "alibabacloudstack_network_interface" "default" {
  name           = var.name
  vswitch_id     = alibabacloudstack_vswitch.default.id
  security_groups = [alibabacloudstack_security_group.default.id]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
  instance_id           = alibabacloudstack_instance.default.id
  network_interface_id  = alibabacloudstack_network_interface.default.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the ECS instance to which the elastic network interface (ENI) will be attached. Changing this value will force the creation of a new attachment.
* `network_interface_id` - (Required, ForceNew) The ID of the elastic network interface (ENI) that will be attached to the specified instance. Changing this value will force the creation of a new attachment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the ENI attachment resource. It is formatted as `<network_interface_id>:<instance_id>`.
* `instance_id` - (Computed) The ID of the ECS instance to which the elastic network interface (ENI) is attached. 
* `network_interface_id` - (Computed) The ID of the elastic network interface (ENI) attached to the specified instance. 