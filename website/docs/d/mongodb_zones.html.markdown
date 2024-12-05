---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-mongodb-zones"
description: |-
    Provides a list of availability zones for mongoDB that can be used by an Alibaba Cloud account.
---

# alibabacloudstack\_mongodb\_zones

This data source provides availability zones for mongoDB that can be accessed by an Alibaba Cloud account within the region configured in the provider.

## Example Usage

```
# Declare the data source
data "alibabacloudstack_mongodb_zones" "zones_ids" {}

# Create an mongoDB instance with the first matched zone
resource "alibabacloudstack_mongodb_instance" "mongodb" {
    zone_id = data.alibabacloudstack_mongodb_zones.zones_ids.zones[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch MongoDB instances.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.

