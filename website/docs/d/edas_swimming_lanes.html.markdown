---
subcategory: "Enterprise Distributed Application Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lanes"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-swimming-lanes"
description: |-
  Provides a list of Edas Swimming Lanes to the user.
---

# alibabacloudstack\_edas\_swimming\_lanes

This data source provides the Edas Swimming Lanes of the current Alibaba Cloud user.

## Example Usage

```terraform
data "alibabacloudstack_edas_swimming_lanes" "example" {
  group_id = "12345"
  logical_region_id = "cn-beijing:test"
}

output "first_swimming_lane_id" {
  value = data.alibabacloudstack_edas_swimming_lanes.example.lanes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `logical_region_id` - (Optional) The logical region ID of the EDAS instance. If not specified, the default region ID of the provider will be used.
* `group_id` - (Required) The ID of the swimming lane group.
* `ids` - (Optional) A list of Swimming Lane IDs.
* `name_regex` - (Optional) A regex string to filter results by Swimming Lane name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Swimming Lane IDs.
* `lanes` - A list of Swimming Lanes. Each element contains the following attributes:
  * `id` - The ID of the Swimming Lane.
  * `name` - The name of the Swimming Lane.
  * `group_id` - The group ID of the Swimming Lane.
  * `logical_region_id` - The logical region ID of the Swimming Lane.
  * `apps` - A list of application IDs associated with this Swimming Lane.
  * `priority` - The priority of the Swimming Lane.
  * `path` - The path of the Swimming Lane.
  * `condition` - The condition of the Swimming Lane.
  * `rest_items` - The REST items of the Swimming Lane.
    * `type` - The type of the REST item.
    * `name` - The name of the REST item.
    * `value` - The value of the REST item.
    * `cond` - The condition of the REST item.
    * `operator` - The operator of the REST item.
  * `enabled` - Whether the Swimming Lane is enabled.
  * `lane_id` - The lane ID of the Swimming Lane.