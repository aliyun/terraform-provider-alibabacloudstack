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
resource "alibabacloudstack_ascm_resource_group" "default" {
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

* `name` - (Required) The name of the resource group. This name can have a string of 2 to 128 characters.

* `organization_id` - (Deprecated, Optional) This parameter has been deprecated. The resource group will be created under the organization to which the current user belongs. Modifying this parameter will force a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource group. The value is in format `<organization_id>:<resource_group_id>`.
* `name` - The name of the resource group.
* `rg_id` - The ID of the resource group (internal ID returned by API).
* `organization_id` - The ID of the organization to which the resource group belongs.

## Import

ASCM Resource Group can be imported using the organization ID and resource group ID separated by colon, e.g.

```
$ terraform import alibabacloudstack_ascm_resource_group.example 12345:67890
```