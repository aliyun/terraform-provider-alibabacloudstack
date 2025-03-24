---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-adb-zones"
description: |-
    Provides a list of availability zones for ADB that can be used by an Alibaba Cloud account.
---

# alibabacloudstack_adb_zones

This data source provides availability zones for ADB that can be accessed by an Alibaba Cloud account within the region configured in the provider.

## Example Usage

```
# Declare the data source
data "alibabacloudstack_adb_zones" "zones_ids" {}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch ADB instances.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.