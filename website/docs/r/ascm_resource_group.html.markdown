---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_group"
sidebar_current: "docs-alibabacloudstack-resource-ascm-resource-group"
description: |-
  Provides Ascm resource group resource.
---

# alibabacloudstack_ascm_resource_group

-> **NOTE:**  If you need to create different resources in different resource sets in a template, you need to refer to the method in [Mult ResourceGroup](ascm_resource_group_mult.html.markdown).



## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
    name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_resource_group" "default" {
    organization_id = alibabacloudstack_ascm_organization.default.org_id
    name = "Resource_Group_Name"
}

data "alibabacloudstack_ascm_resource_groups" "default" {
    name_regex = alibabacloudstack_ascm_resource_group.default.name
}
output "rg" {
  value = data.alibabacloudstack_ascm_resource_groups.default.*
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource group. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null. 
* `organization_id` - (Required) ID of an Organization.

## Attributes Reference

The following attributes are exported:

* `id` - Name and ID of the resource group. The value is in format `Name:ID`
* `rg_id` - The ID of the resource group. 