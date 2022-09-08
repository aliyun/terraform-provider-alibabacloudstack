---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_association"
sidebar_current: "docs-alibabacloudstack-resource-eip-association"
description: |-
  Provides a ECS EIP Association resource.
---

# alibabacloudstack\_eip\_association

Provides an Alibabacloudstack EIP Association resource for associating Elastic IP to ECS Instance, SLB Instance or Nat Gateway.

-> **NOTE:** `alibabacloudstack_eip_association` is useful in scenarios where EIPs are either
 pre-existing or distributed to customers or users and therefore cannot be changed.

-> **NOTE:** The resource support to associate EIP to SLB Instance or Nat Gateway.

-> **NOTE:** One EIP can only be associated with ECS or SLB instance which in the VPC.

## Example Usage

```
# Create a new EIP association and use it to associate a EIP form a instance.

data "alibabacloudstack_zones" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "vsw" {
  vpc_id            = "${alibabacloudstack_vpc.vpc.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"

  depends_on = [
    "alibabacloudstack_vpc.vpc",
  ]
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_instance" "ecs_instance" {
  image_id          = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type     = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  security_groups   = ["${alibabacloudstack_security_group.group.id}"]
  vswitch_id        = "${alibabacloudstack_vswitch.vsw.id}"
  instance_name     = "hello"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_eip" "eip" {
}

resource "alibabacloudstack_eip_association" "eip_asso" {
  allocation_id = "${alibabacloudstack_eip.eip.id}"
  instance_id   = "${alibabacloudstack_instance.ecs_instance.id}"
}

resource "alibabacloudstack_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
  vpc_id      = "${alibabacloudstack_vpc.vpc.id}"
}
```


## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required, ForcesNew) The allocation EIP ID.
* `instance_id` - (Required, ForcesNew) The ID of the ECS or SLB instance or Nat Gateway.
* `instance_type` - (Optional, ForceNew) The type of cloud product that the eip instance to bind.


## Attributes Reference

The following attributes are exported:

* `allocation_id` - As above.
* `instance_id` - As above.
