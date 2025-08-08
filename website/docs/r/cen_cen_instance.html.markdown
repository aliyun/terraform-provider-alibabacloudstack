---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstance"
sidebar_current: "docs-Alibabacloudstack-cen-ceninstance"
description: |-
  Provides a cen Ceninstance resource.
---

# alibabacloudstack\_cen\_ceninstance

Provides a cen Ceninstance resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccceninstance48958"
}


resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
}
```

## Argument Reference

The following arguments are supported:
  * `cen_instance_name` - (Optional) - the cen instance name.
  * `description` - (Optional) - the cen instance description.
  * `protection_level` - (Optional) - the protection level.
  * `status` - (Computed) - the status of cen instance.
  * `transit_router_name` - (Optional) - the transit router name.
  * `transit_router_description` - (Optional) - the transit router description.
  * `transit_router_cidrs` - (Optional) - the transit router cidrs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `cen_id` - the cen instance id.
  * `create_time` - the create time.
  * `status` - the cen instance status.
  * `transit_router_id` - the transit router id.
