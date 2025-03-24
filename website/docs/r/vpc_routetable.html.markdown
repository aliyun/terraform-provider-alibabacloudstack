---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetable"
sidebar_current: "docs-Alibabacloudstack-vpc-routetable"
description: |- 
  Provides a vpc Routetable resource.
---

# alibabacloudstack_vpc_routetable
-> **NOTE:** Alias name has: `alibabacloudstack_route_table`

Provides a vpc Routetable resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
	default = "tf-testaccvpcroute_table19406"
}

resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name       = "${var.name}"
}

resource "alibabacloudstack_vpc_routetable" "default" {
	vpc_id      = "${alibabacloudstack_vpc.default.id}"
	name        = "${var.name}"
	description = "A detailed description of the route table."
	tags = {
		Environment = "Test"
	}
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the route table belongs. This field cannot be modified after creation.
* `name` - (Optional) The name of the route table. If not specified, Terraform will automatically generate a unique name.
* `description` - (Optional) A detailed description of the route table. This helps in identifying the purpose or usage of the route table.
* `tags` - (Optional) A mapping of tags to assign to the route table. These tags can be used for categorization and cost allocation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the route table instance.
* `route_table_name` - The name of the route table as specified by the `name` argument or auto-generated if not provided.