---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organization"
sidebar_current: "docs-alibabacloudstack-resource-ascm-organization"
description: |-
  Provides an Ascm organization resource.
---

# alibabacloudstack\_ascm_organization

Provides an Ascm organization resource.

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "apsara_Organization"
  parent_id = "19"
}
output "org" {
  value = alibabacloudstack_ascm_organization.default.*
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `parent_id` - (Optional) ID of its Parent organization. Default value for a parent_id is "1". For normal user (not an admin) parent_id will be its organization ID.
* `person_num` - (Optional) A reserved parameter.
* `resource_group_num` - (Optional) A reserved parameter.

## Attributes Reference

The following attributes are exported:

* `id` - Name and ID of the organization. The value is in format `Name:ID`
* `org_id` - The ID of the organization.
