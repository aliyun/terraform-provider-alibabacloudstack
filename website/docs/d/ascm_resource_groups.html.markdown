---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_groups"
sidebar_current: "docs-alibabacloudstack-datasource-alibabacloudstack-ascm-resource-groups"
description: |-
    Provides a list of Resource Groups owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack\_ascm\_resource\_groups

This data source provides a list of Resource Groups owned by an Alibabacloudstack Cloud account.


## Example Usage

```
data "alibabacloudstack_ascm_resource_groups" "default" {
  name_regex = "another resource" #Optional
}

output "resource_group" {
  value = data.alibabacloudstack_ascm_resource_groups.default.groups
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Resource Groups IDs.
* `name_regex` - (Optional) A regex string to filter results by name of Resource Group.
* `organization_id` - (Optional) Organization ID Alibabacloudstack Cloud account.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Resource Groups IDs.
* `names` - A list of Resource Groups names.
* `groups` - A list of Resource Groups. Each element contains the following attributes:
  * `id` - ID of the Resource Group.
  * `name` - Name of Resource Group.
  * `organization_id` - Organization ID for Alibabacloudstack Cloud account.
  * `gmt_created` - Create time of the Resource Groups.
  * `rs_id` - ID of the Resource Group Set (like rs-xxxxx).
  * `creator` - ID of the Resource Group creator.
  * `resource_group_type` - Resource set type. Available values and meanings:
    * 1: Default resource set.
    * 0: Non default resource set.
  
