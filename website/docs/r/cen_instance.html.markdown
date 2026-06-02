---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_instance"
sidebar_current: "docs-Alibabacloudstack-cen-instance"
description: |-
  Provides a cen Instance resource.
---

# alibabacloudstack\_cen\_instance

Provides a cen Instance resource.

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

* `cen_instance_name` - (Optional) The name of the CEN instance.
* `description` - (Optional) The description of the CEN instance.
* `protection_level` - (Optional) The protection level of the CEN instance. Default value: `REDUCED`.
* `transit_router_name` - (Optional) The name of the transit router.
* `transit_router_description` - (Optional) The description of the transit router.
* `transit_router_cidrs` - (Optional, Available in v3.18+) The CIDR blocks of the transit router. Up to 5 CIDR blocks can be configured. Each item supports:
  * `cidr` - (Required) The CIDR block.
  * `cidr_id` - (Computed) The ID of the CIDR block.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the CEN instance.
* `cen_id` - The ID of the CEN instance.
* `cen_bandwidth_package_ids` - The IDs of the bandwidth packages associated with the CEN instance.
* `create_time` - The time when the CEN instance was created.
* `status` - The status of the CEN instance.
* `transit_router_id` - The ID of the transit router.

## Import

CEN Instance can be imported using the CenId, e.g.

```
$ terraform import alibabacloudstack_cen_instance.example cen-xxxxxxxxxxx
```
