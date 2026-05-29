---
subcategory: "Time Series Database (TSDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_tsdb_zones"
sidebar_current: "docs-Alibabacloudstack-datasource-tsdb-zones"
description: |-
  Provides a list of Time Series Database (TSDB) instance available zones to the user.
---

# alibabacloudstack_tsdb_zones

This data source provides the available zones with the Time Series Database (TSDB) Instance of the current Alibaba Cloud user.


## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_tsdb_zones" "example" {}

output "first_tsdb_zones_id" {
  value = data.alibabacloudstack_tsdb_zones.example.zones.0.zone_id
}
```

## Argument Reference

This data source requires no arguments.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of TSDB instance zone IDs.
* `zones` - A list of TSDB Instance zones. Each element contains the following attributes:
  * `id` - The ID of the zone.
  * `zone_id` - The zone ID.

> **NOTE:** The `local_name` attribute is not supported in Apsara Stack v3.16 and later versions.