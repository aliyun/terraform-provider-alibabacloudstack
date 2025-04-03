---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_zones"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-zones"
description: |-
  Provides a list of polardb zones owned by an alibabacloudstack account.
---

# alibabacloudstack_polardb_zones

This data source provides the polardb zones on alibabacloudstack.

## Example Usage

```hcl
data "alibabacloudstack_polardb_zones" "example" {
  multi = false
}

output "zones" {
  value = data.alibabacloudstack_polardb_zones.example.zones
}
```

## Argument Reference
The following arguments are supported:

* `multi` - (Optional) Indicates whether to retrieve multi-zone IDs. Default is false.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the zones.
* `zones` - A list of zones. Each element contains the following attributes:
    * `id` - The unique identifier of the zone.
    * `multi_zone_ids` - A list of multi-zone IDs (only available if multi is set to true).
