---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vswitch"
sidebar_current: "docs-Alibabacloudstack-vpc-vswitch"
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

* `zone_id` - (Required) The ID of the zone to which the vSwitches belong. You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/36064.html) operation to query the most recent zone list.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) to which the vSwitches belong.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable the switch IPv6 CIDR block. Valid values:
  * `false` (Default): disables IPv6 CIDR blocks.
  * `true`: enables IPv6 CIDR blocks. If the `enable_ipv6` is `true`, IPv6 must also be enabled for the VPC directed by the `vpc_id`. The system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.
* `vswitch_name` - (Optional) The name of the vSwitch. Defaults to null.
* `description` - (Optional) The description of the vSwitch. The description must be 1 to 256 characters in length and cannot start with `http://` or `https://`.
* `tags` - (Optional, Map) The tags of the VSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vSwitch.
* `zone_id` - The ID of the zone to which the vSwitches belong.
* `availability_zone` - The AZ for the vSwitch.
* `cidr_block` - The CIDR block for the vSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the vSwitch.
* `vpc_id` - The VPC ID.
* `vswitch_name` - The name of the vSwitch.
* `description` - The description of the vSwitch.
