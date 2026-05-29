---
subcategory: "Elastic Compute Service(ECS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhostclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicatedhostclusters"
description: |-
  Provides a list of ecs dedicatedhostclusters owned by an alibabacloudstack account.
---

# alibabacloudstack\_ecs\_dedicatedhostclusters

This data source provides a list of ecs dedicatedhostclusters in an alibabacloudstack account according to the specified filters.

## Example Usage
```
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


		resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
		  dedicated_host_cluster_name = "tf_testAccEcsDedicatedHostsClusterDataSource_5238343"
          zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
		}
	

data "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
  ids = [
          "${alibabacloudstack_ecs_dedicated_host_cluster.default.id}"
        ]
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - Thte list of dedicated host cluster ids to filter results.
  * `zone_id` - (Optional) - zone id
  * `dedicated_host_cluster_name_regex` - (Optional) - A regex string to filter results by dedicated host cluster name.
  * `dedicated_host_cluster_name` - (Optional) - dedicated host cluster name

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `dedicated_host_clusters` - A list of all the dedicated host clusters that match the specified filters.
    * `id` - The id of the dedicated host cluster.
    * `dedicated_host_cluster_id` - dedicated host cluster id
    * `dedicated_host_cluster_name` - dedicated host cluster name
    * `description` - description
    * `zone_id` - zone id
