---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetableattachment"
sidebar_current: "docs-Alibabacloudstack-vpc-routetableattachment"
description: |- 
  Provides a vpc Routetableattachment resource.
---

# alibabacloudstack_vpc_routetableattachment
-> **NOTE:** Alias name has: `alibabacloudstack_route_table_attachment`

Provides a vpc Routetableattachment resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccvpcroute_table_attachment24025"
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
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_route_table" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_description"
}

resource "alibabacloudstack_vpc_routetableattachment" "default" {
  route_table_id = "${alibabacloudstack_route_table.default.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, ForceNew) The ID of the route table to be bound to the switch. This field cannot be modified after creation.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch to which the route table will be attached. This field cannot be modified after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the route table attachment. It is formatted as `<route_table_id>:<vswitch_id>`.