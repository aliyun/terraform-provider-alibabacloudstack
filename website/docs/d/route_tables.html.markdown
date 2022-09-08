---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_route_tables"
sidebar_current: "docs-alibabacloudstack-datasource-route-tables"
description: |-
    Provides a list of Route Tables owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack\_route\_tables

This data source provides a list of Route Tables owned by an Alibabacloudstack Cloud account.


## Example Usage

```
variable "name" {
  default = "route-tables-datasource-example-name"
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

data "alibabacloudstack_route_tables" "foo" {
  ids = ["${alibabacloudstack_route_table.foo.id}"]
}

output "route_table_ids" {
  value = "${data.alibabacloudstack_route_tables.foo.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Route Tables IDs.
* `name_regex` - (Optional) A regex string to filter route tables by name.
* `vpc_id` - (Optional) Vpc id of the route table.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Route Tables IDs.
* `names` - A list of Route Tables names.
* `tables` - A list of Route Tables. Each element contains the following attributes:
  * `id` - ID of the Route Table.
  * `router_id` - Router Id of the route table.
  * `route_table_type` - The type of route table.
  * `name` - Name of the route table.
  * `description` - The description of the route table instance.
  * `creation_time` - Time of creation.
  
