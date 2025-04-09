---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterface"
sidebar_current: "docs-Alibabacloudstack-ecs-networkinterface"
description: |- 
  Provides a ecs Networkinterface resource.
---

# alibabacloudstack_ecs_networkinterface
-> **NOTE:** Alias name has: `alibabacloudstack_network_interface`

Provides a ecs Networkinterface resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccecsnetwork_interface15831"
}

resource "alibabacloudstack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.0.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
  name       = "${var.name}"
  vpc_id     = alibabacloudstack_vpc.vpc.id
  cidr_block = "10.0.0.0/24"
  availability_zone = "cn-beijing-b"
}

resource "alibabacloudstack_security_group" "secgroup" {
  name        = "${var.name}"
  description = "Security Group"
  vpc_id      = alibabacloudstack_vpc.vpc.id
}

resource "alibabacloudstack_ecs_networkinterface" "default" {
  network_interface_name = "${var.name}-eni"
  vswitch_id             = alibabacloudstack_vswitch.vsw.id
  security_groups        = [alibabacloudstack_security_group.secgroup.id]
  primary_ip_address     = "10.0.0.10"
  private_ips_count      = 2
  description            = "Test ENI"
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_name` - (Optional) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `vswitch_id` - (Required, ForceNew) The VSwitch ID to create the ENI in.
* `security_groups` - (Required) A list of security group IDs to associate with the ENI.
* `primary_ip_address` - (Optional, ForceNew) The primary private IP address of the ENI. If not specified, Alibaba Cloud will automatically assign one within the CIDR block of the VSwitch.
* `private_ips` - (Optional) List of secondary private IPs to assign to the ENI. Do not use both `private_ips` and `private_ips_count` in the same ENI resource block.
* `private_ips_count` - (Optional) Number of secondary private IPs to assign to the ENI. Do not use both `private_ips` and `private_ips_count` in the same ENI resource block.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `mac_address` - (Optional) The MAC address of the ENI.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ENI ID.
* `mac_address` - The MAC address of the ENI.
* `network_interface_name` - Name of the ENI.
* `primary_ip_address` - The primary private IP address of the ENI.
* `private_ips` - List of all private IP addresses assigned to the ENI, including the primary IP and any secondary IPs.
* `private_ips_count` - Total number of private IPs assigned to the ENI, including the primary IP and any secondary IPs.