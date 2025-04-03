---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_zones"
sidebar_current: "docs-alibabacloudstack-datasource-elasticsearch-zones"
description: |-
    Provides a list of availability zones for Elasticsearch that can be used by an Alibaba Cloud account.
---

# alibabacloudstack_elasticsearch_zones

This data source provides availability zones for Elasticsearch that can be accessed by an Alibaba Cloud account within the region configured in the provider.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_elasticsearch_zones" "zones_ids" {}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch Elasticsearch instances.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.