---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_route_table"
sidebar_current: "docs-apsarastack-resource-route-table"
description: |-
  Provides a Apsarastack Route Table resource.
---

# apsarastack\_route_table

Provides a route table resource to add customized route tables.

-> **NOTE:** Terraform will auto build route table instance while it uses `apsarastack_route_table` to build a route table resource.

## Example Usage

Basic Usage

```
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name       = "vpc-example-name"
}

resource "apsarastack_route_table" "foo" {
  vpc_id      = "${apsarastack_vpc.foo.id}"
  name        = "route-table-example-name"
  description = "route-table-example-description"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The vpc_id of the route table, the field can't be changed.
* `name` - (Optional) The name of the route table.
* `description` - (Optional) The description of the route table instance.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the route table instance id.



