---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-instances"
description: |-
  Provides a list of cen Instances owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_Instances

This data source provides a list of cen Instances in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccInstancesDatasource15452"
}

resource "alibabacloudstack_cen_instance" "default" {
    cen_instance_name = "${var.name}"
	description = "${var.name}"
}

data "alibabacloudstack_cen_instances" "default" {
	name_regex = "${alibabacloudstack_cen_instance.default.cen_instance_name}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - cen instance ids.
  * `name_regex` - (Optional) - cen instance name regex.
  * `description_regex` - (Optional) - cen instance description regex.
  * `transit_router_name_regex` - (Optional) -cen instance transit router name regex.
  * `transit_router_description_regex` - (Optional) - cen instance transit router description regex.
  * `cidr` - (Optional) - cen instance router cidr.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `cens` - cen instances list.
