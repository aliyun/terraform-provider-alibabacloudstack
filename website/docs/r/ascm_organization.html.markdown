---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organization"
sidebar_current: "docs-alibabacloudstack-resource-ascm-organization"
description: |-
  Provides an Ascm organization resource.
---

# alibabacloudstack_ascm_organization

Provides an Ascm organization resource.

-> **Note:** This resource can also be referred to by the following aliases:
- `apsarastack_ascm_organization`

## Example Usage

```hcl
resource "alibabacloudstack_ascm_organization" "default" {
  name        = "apsara_Organization"
  parent_id   = "1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization. This name can have a length of 2 to 128 characters.
* `parent_id` - (Optional) The ID of the parent organization. Default value: `"1"`.
* `person_num` - (Optional) A reserved parameter.
* `resource_group_num` - (Optional) A reserved parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization.
* `org_id` - The UUID of the organization.
* `primary_key` - The primary key associated with the organization.
* `aliyunid` - The Aliyun ID associated with the organization.

## Import

Ascm Organization can be imported using the organization ID, e.g.

```
$ terraform import alibabacloudstack_ascm_organization.example 12345
```