---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_network_interface"
sidebar_current: "docs-apsarastack-resource-network-interface"
description: |-
  Provides an ECS Elastic Network Interface resource.
---

# apsarastack\_network\_interface

Provides an ECS Elastic Network Interface resource.

## Example Usage

```
resource "apsarastack_security_group" "secgroup" {
  name        = "SecurityGroup-Name"
  description = "Security Group"
  vpc_id      = apsarastack_vpc.vpc.id
}
resource "apsarastack_vpc" "vpc" {
  name       = "VPC-Name"
  cidr_block = "10.0.0.0/16"
}

resource "apsarastack_vswitch" "vsw" {
  name       = "VSW-Name"
  vpc_id            = apsarastack_vpc.vpc.id
  cidr_block        = apsarastack_vpc.vpc.cidr_block
  availability_zone = "cn-beijing-b"
}
resource "apsarastack_instance" "apsarainstance" {
  image_id              = "gj2j1g3-45h3nnc-454hj5g"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  security_groups      = [apsarastack_security_group.secgroup.id]
  instance_name        = "apsarainstance"
  vswitch_id           = apsarastack_vswitch.vsw.id
}

resource "apsarastack_network_interface" "NetInterface" {
  name              = "ENI"
  vswitch_id        = apsarastack_vswitch.vsw.id
  security_groups   = apsarastack_security_group.secgroup.id
  private_ips_count = 1
  description = "Network Interface"
}
```

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The VSwitch to create the ENI in.
* `security_groups` - (Required) A list of security group ids to associate with.
* `private_ip` - (Optional, ForceNew) The primary private IP of the ENI.
* `name` - (Optional) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `private_ips`  - (Optional) List of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `private_ips_count` - (Optional) Number of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ENI ID.
* `mac` -The MAC address of an ENI.


