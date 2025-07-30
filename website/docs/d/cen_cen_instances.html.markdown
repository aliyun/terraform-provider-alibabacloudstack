---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstances"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-ceninstances"
description: |-
  Provides a list of cen ceninstances owned by an alibabacloudstack account.
---

# alibabacloudstack\_cen\_ceninstances

This data source provides a list of cen ceninstances in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccCenInstancesDatasource15452"
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

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `cens` - cen instances.
