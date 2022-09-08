---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_route_table_attachment"
sidebar_current: "docs-alibabacloudstack-resource-route-table-attachment"
description: |-
  Provides an alibabacloudstack Route Table Attachment resource.
---

# alibabacloudstack\_route\_table\_attachment

Provides an alibabacloudstack Route Table Attachment resource for associating Route Table to VSwitch Instance.

-> **NOTE:** Terraform will auto build route table attachment while it uses `alibabacloudstack_route_table_attachment` to build a route table attachment resource.


## Example Usage

Basic Usage

```
variable "name" {
  default = "route-table-attachment-example-name"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_route_table" "foo" {
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
  name        = "${var.name}"
  description = "route_table_attachment"
}

resource "alibabacloudstack_route_table_attachment" "foo" {
  vswitch_id     = "${alibabacloudstack_vswitch.foo.id}"
  route_table_id = "${alibabacloudstack_route_table.foo.id}"
}
```
## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch_id of the route table attachment, the field can't be changed.
* `route_table_id` - (Required, ForceNew) The route_table_id of the route table attachment, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the route table attachment id and formates as `<route_table_id>:<vswitch_id>`.

