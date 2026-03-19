---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane_groups"
sidebar_current: "docs-alibabacloudstack-datasource-edas-swimming-lane-groups"
description: |-
  Provides a list of Edas Swimming Lane Groups to the user.
---

# alibabacloudstack\_edas\_swimming\_lane\_groups

This data source provides the Edas Swimming Lane Groups of the current Alibaba Cloud user.

## Example Usage

```terraform
data "alibabacloudstack_edas_swimming_lane_groups" "example" {
  logical_region_id = "cn-beijing:test"
}

output "first_swimming_lane_group_id" {
  value = data.alibabacloudstack_edas_swimming_lane_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `logical_region_id` - (Optional) The logical region ID of the EDAS instance. If not specified, the default region ID of the provider will be used.
* `ids` - (Optional) A list of Swimming Lane Group IDs.
* `name_regex` - (Optional) A regex string to filter results by Swimming Lane Group name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Swimming Lane Group IDs.
* `groups` - A list of Swimming Lane Groups. Each element contains the following attributes:
  * `id` - The ID of the Swimming Lane Group.
  * `name` - The name of the Swimming Lane Group.
  * `entry_app_id` - The ID of the entry application.
  * `apps` - A list of application IDs associated with this Swimming Lane Group.
  * `logical_region_id` - The logical region ID of the Swimming Lane Group.
  * `strategy_type` - The strategy type of the Swimming Lane Group.
  * `group_id` - The group ID of the Swimming Lane Group.