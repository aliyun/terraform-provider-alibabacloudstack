---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_instance_families"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-instance-families"
description: |-
    Provides a list of instance families to the user.
---

# alibabacloudstack_ascm_instance_families

This data source provides the instance families of the current Apsara Stack Cloud user.

## Example Usage

```
data "alibabacloudstack_ascm_instance_families" "default" {
  output_file = "instance_families"
  resource_type = "DRDS"
  status = "Available"
}
output "instfam" {
  value = data.alibabacloudstack_ascm_instance_families.default.*
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance family IDs.
* `name_regex` - (Optional) A regex string to filter results by trail name.
* `status` - (Optional) Specify Status to filter the resulting instance families by their availability.
* `resource_type` - (Optional) Filter the results by the specified resource type.
* `families` - (Optional) A list of instance families. Each element contains the following attributes:

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `families` - A list of instance families. Each element contains the following attributes:
    * `id` - ID of the instance families.
    * `order_by_id` - Sorted ID of the instacne families.
    * `series_name` - Series name for instance families.
    * `modifier` - Modifier name.
    * `series_name_label` - label of Series name for instance families.
    * `is_deleted` - Specify the state in "Y" or "N" form.
    * `resource_type` - Specified resource type.
    * `computed_attribute` - An example computed attribute that is automatically generated.