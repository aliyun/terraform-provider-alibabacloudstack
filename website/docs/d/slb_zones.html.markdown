---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-slb-zones"
description: |-
    Provides a list of availability zones for SLB that can be used by an AlibabacloudStack Cloud account.
---

# alibabacloudstack_slb_zones

This data source provides availability zones for SLB that can be accessed by an AlibabacloudStack Cloud account within the region configured in the provider.


## Example Usage

```
# Declare the data source
data "alibabacloudstack_slb_zones" "zones_ids" {}

output "slb_zones" {
  value = data.alibabacloudstack_slb_zones.zones_ids.*
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to false and only output `id` in the `zones` block. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `slb_slave_zone_ids` - A list of slb slave zone ids in which the slb master zone.
  * `local_name` - The name of the secondary zone. 
  * `computed_attribute_example` - Example of a computed attribute. 