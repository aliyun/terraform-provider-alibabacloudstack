---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vswitch"
sidebar_current: "docs-Alibabacloudstack-resource-vpc-vswitch"
description: |- 
  Provides a vpc VSwitch resource.
---

# alibabacloudstack_vpc_vswitch
-> **NOTE:** Alias name has: `alibabacloudstack_vswitch`

Provides a vpc VSwitch resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
    default = "tf-testaccvpcvswitch97984"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  description     = "modify_description"
  vswitch_name   = "tf-testaccvpcvswitch97984"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}
```

IPv6 Enabled Usage

```hcl
resource "alibabacloudstack_vpc_vpc" "ipv6_example" {
  vpc_name       = "ipv6_vpc"
  cidr_block     = "192.168.0.0/16"
  enable_ipv6    = true
}

resource "alibabacloudstack_vpc_vswitch" "ipv6_vswitch" {
  vswitch_name   = "ipv6_vswitch"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.ipv6_example.id}"
  cidr_block     = "192.168.0.0/24"
  enable_ipv6    = true
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) to which the vSwitches belong.
* `cidr_block` - (Optional, ForceNew) The CIDR block for the switch. The CIDR block requirements are as follows:
  * The mask length of the CIDR block ranges from 16 to 29 bits.
  * The CIDR block must be a subset of the VPC CIDR block.
  * The CIDR block cannot be the same as the destination CIDR block of any route entry in the VPC, but can be a subset of it.
  * The CIDR block cannot be 100.64.0.0/10 or its subnet.
* `zone_id` - (Optional) The ID of the zone to which the vSwitches belong. You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/36064.html) operation to query the most recent zone list.
* `vswitch_name` - (Optional) The name of the vSwitch. Defaults to null.
* `description` - (Optional) The description of the vSwitch. The description must be 1 to 256 characters in length and cannot start with `http://` or `https://`.
* `enable_ipv6` - (Optional) Specifies whether to enable the switch IPv6 CIDR block. Valid values:
  * `false` (Default): disables IPv6 CIDR blocks.
  * `true`: enables IPv6 CIDR blocks. If the `enable_ipv6` is `true`, IPv6 must also be enabled for the VPC directed by the `vpc_id`. The system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.
* `is_cgw` - (Optional) Specifies whether the vSwitch is a CGW (Cloud Gateway) vSwitch.
* `tags` - (Optional, Map) The tags of the VSwitch.

-> **NOTE:** The parameter `availability_zone` is deprecated. Please use `zone_id` instead.
-> **NOTE:** The parameter `name` is deprecated. Please use `vswitch_name` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vSwitch.
* `zone_id` - The ID of the zone to which the vSwitches belong.
* `availability_zone` - (Deprecated) The AZ for the vSwitch. Please use `zone_id` instead.
* `cidr_block` - The CIDR block for the vSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the vSwitch.
* `vpc_id` - The VPC ID.
* `vswitch_name` - The name of the vSwitch.
* `name` - (Deprecated) The name of the vSwitch. Please use `vswitch_name` instead.
* `description` - The description of the vSwitch.
* `is_cgw` - Indicates whether the vSwitch is a CGW vSwitch.

## Import

VSwitch can be imported using the VSwitchId, e.g.

```
$ terraform import alibabacloudstack_vswitch.example vsw-12345678
```
