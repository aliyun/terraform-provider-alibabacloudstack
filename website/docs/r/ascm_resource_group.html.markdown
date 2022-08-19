---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_resource_group"
sidebar_current: "docs-apsarastack-resource-ascm-resource-group"
description: |-
  Provides Ascm resource group resource.
---

# apsarastack\_ascm_resource_group

## Example Usage

```
resource "apsarastack_ascm_organization" "default" {
    name = "Dummy_Test_1"
}

resource "apsarastack_ascm_resource_group" "default" {
    organization_id = apsarastack_ascm_organization.default.org_id
    name = "Resource_Group_Name"
}

data "apsarastack_ascm_resource_groups" "default" {
    name_regex = apsarastack_ascm_resource_group.default.name
}
output "rg" {
  value = data.apsarastack_ascm_resource_groups.default.*
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
