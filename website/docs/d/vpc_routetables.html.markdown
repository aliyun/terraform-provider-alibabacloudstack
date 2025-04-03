---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetables"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-routetables"
description: |- 
  Provides a list of vpc routetables owned by an alibabacloudstack account.
---

# alibabacloudstack_vpc_routetables
-> **NOTE:** Alias name has: `alibabacloudstack_route_tables`

This data source provides a list of vpc routetables in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "vpc-routetables-datasource-example-name"
}

resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

resource "alibabacloudstack_route_table" "foo" {
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
  name        = "${var.name}"
  description = "${var.name}"
}

data "alibabacloudstack_vpc_routetables" "foo" {
  ids         = ["${alibabacloudstack_route_table.foo.id}"]
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
  name_regex  = "${var.name}"
}

output "route_table_ids" {
  value = "${data.alibabacloudstack_vpc_routetables.foo.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Route Table IDs. 
* `name_regex` - (Optional) A regex string to filter Route Tables by name.
* `vpc_id` - (Optional) The ID of the VPC to which the Route Table belongs.
* `tags` - (Optional) A mapping of tags, each tag is represented as key-value pairs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of Route Table IDs.
* `names` - A list of Route Table names.
* `tables` - A list of Route Tables. Each element contains the following attributes:
  * `id` - The ID of the Route Table.
  * `router_id` - The router ID to which the routing table belongs.
  * `route_table_type` - The type of routing table. Values:
    - `custom`: Custom routing table.
    - `system`: System routing table.
  * `name` - The name of the Route Table.
  * `description` - The description of the Route Table.
  * `creation_time` - The time when the Route Table was created.