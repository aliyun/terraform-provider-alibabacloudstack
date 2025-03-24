---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_zones"
sidebar_current: "docs-alibabacloudstack-datasource-db-zones"
description: |-
    Provides a list of availability zones for RDS that can be used by an Alibabacloudstack Cloud account.
---

# alibabacloudstack_db_zones

This data source provides availability zones for RDS that can be accessed by an Alibabacloudstack Cloud account within the region configured in the provider.


## Example Usage

```
# Declare the data source
data "alibabacloudstack_db_zones" "zones_ids" {}

output "db_zones" {
  value = data.alibabacloudstack_db_zones.zones_ids.zones.*
}

```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch RDS instances.

* `ids` - (Optional) A list of zone IDs. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.

* `multi` - Indicates if the zones can be used in a multi AZ configuration. 