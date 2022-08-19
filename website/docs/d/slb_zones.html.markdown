---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_zones"
sidebar_current: "docs-apsarastack-datasource-slb-zones"
description: |-
    Provides a list of availability zones for SLB that can be used by an ApsaraStack Cloud account.
---

# apsarastack\_slb\_zones

This data source provides availability zones for SLB that can be accessed by an ApsaraStack Cloud account within the region configured in the provider.


## Example Usage

```
# Declare the data source
data "apsarastack_slb_zones" "zones_ids" {}

output "slb_zones" {
  value = data.apsarastack_slb_zones.zones_ids.*
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to false and only output `id` in the `zones` block. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `slb_slave_zone_ids` - A list of slb slave zone ids in which the slb master zone.
   * `local_name` - The name of the secondary zone.
