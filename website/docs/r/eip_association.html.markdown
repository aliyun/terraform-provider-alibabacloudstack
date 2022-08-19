---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_eip_association"
sidebar_current: "docs-apsarastack-resource-eip-association"
description: |-
  Provides a ECS EIP Association resource.
---

# apsarastack\_eip\_association

Provides an Apsarastack EIP Association resource for associating Elastic IP to ECS Instance, SLB Instance or Nat Gateway.

-> **NOTE:** `apsarastack_eip_association` is useful in scenarios where EIPs are either
 pre-existing or distributed to customers or users and therefore cannot be changed.

-> **NOTE:** The resource support to associate EIP to SLB Instance or Nat Gateway.

-> **NOTE:** One EIP can only be associated with ECS or SLB instance which in the VPC.

## Example Usage

```
# Create a new EIP association and use it to associate a EIP form a instance.

data "apsarastack_zones" "default" {
}

resource "apsarastack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "vsw" {
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"

  depends_on = [
    "apsarastack_vpc.vpc",
  ]
}

data "apsarastack_instance_types" "default" {
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "apsarastack_instance" "ecs_instance" {
  image_id          = "${data.apsarastack_images.default.images.0.id}"
  instance_type     = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  security_groups   = ["${apsarastack_security_group.group.id}"]
  vswitch_id        = "${apsarastack_vswitch.vsw.id}"
  instance_name     = "hello"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "apsarastack_eip" "eip" {
}

resource "apsarastack_eip_association" "eip_asso" {
  allocation_id = "${apsarastack_eip.eip.id}"
  instance_id   = "${apsarastack_instance.ecs_instance.id}"
}

resource "apsarastack_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
  vpc_id      = "${apsarastack_vpc.vpc.id}"
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
