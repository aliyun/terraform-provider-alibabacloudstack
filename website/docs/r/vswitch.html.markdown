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
  ipv6_cidr_block   = "${alibabacloudstack_vpc.vpc.ipv6_cidr_block}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The AZ for the switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.
* `ipv6_cidr_block` - (Optional) The ipv6 cidr block of VPC.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vswitch (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vswitch. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.


