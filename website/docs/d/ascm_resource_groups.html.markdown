---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_resource_groups"
sidebar_current: "docs-apsarastack-datasource-apsarastack-ascm-resource-groups"
description: |-
    Provides a list of Resource Groups owned by an Apsarastack Cloud account.
---

# apsarastack\_ascm\_resource\_groups

This data source provides a list of Resource Groups owned by an Apsarastack Cloud account.


## Example Usage

```
data "apsarastack_ascm_resource_groups" "default" {
  name_regex = "another resource" #Optional
}

output "resource_group" {
  value = data.apsarastack_ascm_resource_groups.default.groups
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Resource Groups IDs.
* `name_regex` - (Optional) A regex string to filter results by name of Resource Group.
* `organization_id` - (Optional) Organization ID Apsarastack Cloud account.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Resource Groups IDs.
* `names` - A list of Resource Groups names.
* `groups` - A list of Resource Groups. Each element contains the following attributes:
  * `id` - ID of the Resource Group.
  * `name` - Name of Resource Group.
  * `organization_id` - Organization ID for Apsarastack Cloud account.
  
