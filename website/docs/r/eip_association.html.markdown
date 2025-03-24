---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_association"
sidebar_current: "docs-Alibabacloudstack-eip-association"
description: |- 
  Provides a eip Association resource.
---

# alibabacloudstack_eip_association

Provides a EIP Association resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccEipAssociation10702"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id           = alibabacloudstack_vpc.default.id
  cidr_block       = "10.1.1.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name             = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name        = var.name
  description = "New security group"
  vpc_id      = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  vswitch_id         = alibabacloudstack_vswitch.default.id
  image_id          = data.alibabacloudstack_images.default.images.0.id
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_ssd"
  instance_type     = local.default_instance_type_id

  security_groups = [alibabacloudstack_security_group.default.id]
  instance_name   = var.name
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_eip" "default" {
  name = var.name
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = alibabacloudstack_eip.default.id
  instance_id   = alibabacloudstack_instance.default.id
  force         = false
  instance_type = "EcsInstance"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForceNew) The ID of the EIP instance.
* `instance_id` - (Required, ForceNew) The ID of the instance with which you want to associate the EIP. You can enter the ID of a NAT gateway, CLB instance, ECS instance, secondary ENI, HAVIP, or IP address.
* `force` - (Optional, ForceNew) Specifies whether to disassociate the EIP from a NAT gateway if a DNAT or SNAT entry is added to the NAT gateway. Valid values:
  * **false** (default): Does not force the unbinding of the EIP.
  * **true**: Forces the unbinding of the EIP.
* `instance_type` - (Optional, ForceNew) The type of the instance with which you want to associate the EIP. Valid values:
  * **Nat**: NAT gateway
  * **SlbInstance**: CLB instance
  * **EcsInstance** (default): ECS instance
  * **NetworkInterface**: Secondary ENI
  * **HaVip**: HAVIP
  * **IpAddress**: IP address
  > The default value is **EcsInstance**. If the instance with which you want to associate the EIP is not an ECS instance, this parameter is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - The ID of the EIP instance.
* `instance_id` - The ID of the instance with which the EIP is associated.
* `instance_type` - The type of the instance with which the EIP is associated. Valid values:
  * **Nat**: NAT gateway
  * **SlbInstance**: CLB instance
  * **EcsInstance** (default): ECS instance
  * **NetworkInterface**: Secondary ENI
  * **HaVip**: HAVIP
  * **IpAddress**: IP address
  > The default value is **EcsInstance**. If the instance with which you want to associate the EIP is not an ECS instance, this parameter is required.