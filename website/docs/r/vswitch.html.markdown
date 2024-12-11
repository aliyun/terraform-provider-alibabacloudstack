---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vswitch"
sidebar_current: "docs-alibabacloudstack-resource-vswitch"
description: |-
  Provides a Alibabacloudstack VPC switch resource.
---

# alibabacloudstack\_vswitch

Provides a VPC switch resource.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_vpc" "vpc" {
  vpc_name       = "${var.vpc_name}"
  cidr_block = "${var.cidr_block}"
  enable_ipv6    = true
}

resource "alibabacloudstack_vswitch" "vsw" {
  vpc_id            = "${alibabacloudstack_vpc.vpc.id}"
  cidr_block        = "${var.cidr_block}"
  availability_zone = "${var.availability_zone}"
  enable_ipv6       = true
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The AZ for the switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable the switch IPv6 CIDR block. Valid values: `false` (Default): disables IPv6 CIDR blocks. `true`: enables IPv6 CIDR blocks. If the `enable_ipv6` is `true`, ipv6 must also be enabled for the vpc directed by the `vpc_id`, the system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.
* `tags` - (Optional, Map) The tags of VSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `ipv6_cidr_block` - (Optional) The ipv6 cidr block of switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.


