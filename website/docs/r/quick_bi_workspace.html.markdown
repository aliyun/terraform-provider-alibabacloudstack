---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_workspace"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-workspace"
description: |-
  Provides a Alibabacloudstack Quick BI Workspace resource.
---

# alibabacloudstack\_quick\_bi\_workspace

Provides a Quick BI Workspace resource.

## Example Usage

Basic Usage

```terraform

resource "alibabacloudstack_quick_bi_workspace" "default" {
  workspace_name = "example_value"
  workspace_desc = "example_value"
  use_comment = "false"
  allow_share = "false"
  allow_publish = "false"
}

```

## Argument Reference

The following arguments are supported:

* `workspace_name` - (Required) Workspace name.
* `workspace_desc` - (Optional) Workspace description.
* `use_comment` - (Optional) Do you want to use table comments when creating data sets (corresponding to preferences). Valid values: `true` and `false`.
* `allow_share` - (Optional) Whether the report is allowed to be shared (corresponding function permission-works can be authorized). Valid values: `false`, `true`.
* `allow_publish` - (Optional) Whether the report is allowed to be made public (corresponding function permission-works can be made public).Valid values: `false`, `true`.

## Attributes Reference

The following attributes are exported:

* `workspace_id` - The resource ID in terraform of Workspace.

## Import

Quick BI Workspace can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_quick_bi_workspace.example <id>
```
