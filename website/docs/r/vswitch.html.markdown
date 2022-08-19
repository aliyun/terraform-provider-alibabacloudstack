---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_vswitch"
sidebar_current: "docs-apsarastack-resource-vswitch"
description: |-
  Provides a Apsarastack VPC switch resource.
---

# apsarastack\_vswitch

Provides a VPC switch resource.

## Example Usage

Basic Usage

```
resource "apsarastack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "${var.cidr_block}"
}

resource "apsarastack_vswitch" "vsw" {
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  cidr_block        = "${var.cidr_block}"
  availability_zone = "${var.availability_zone}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The AZ for the switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `cidr_block` - (Required, ForceNew) The CIDR block for the switch.
* `name` - (Optional) The name of the switch. Defaults to null.
* `description` - (Optional) The switch description. Defaults to null.

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


