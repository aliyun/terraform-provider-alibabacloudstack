---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_baseline"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-baseline"
description: |-
  Provides a DataWorks Baseline resource.
---

# alibabacloudstack_data_works_baseline

Provides a DataWorks Baseline resource.

-> **Note:** This resource can also be referred to by the following alias:
-> - `apsarastack_data_works_baseline`

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_baseline12345"
}

resource "alibabacloudstack_data_works_project" "default" {
  name           = var.name
  description    = "${var.name}_desc"
  task_auth_type = "PROJECT"
}

data "alibabacloudstack_account" "current" {}

data "alibabacloudstack_ascm_users" "default" {
  organization_id = data.alibabacloudstack_account.current.organization_id
}

resource "alibabacloudstack_data_works_user" "default" {
  project_id      = alibabacloudstack_data_works_project.default.id
  user_id         = data.alibabacloudstack_ascm_users.default.users.0.primary_key
  role_code       = ["role_project_admin"]
  lifecycle {
    ignore_changes = [role_code]
  }
}

resource "alibabacloudstack_data_works_baseline" "default" {
  baseline_name          = var.name
  project_id             = alibabacloudstack_data_works_project.default.id
  owner                  = alibabacloudstack_data_works_user.default.project_member_id
  priority               = 5
  baseline_type          = "DAILY"
  alert_margin_threshold = 30
  enabled                = true
  alert_enabled          = true
  overtime_settings {
    cycle = 1
    time  = "08:00"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Parameters

* `baseline_name` - (Required) The name of the baseline.

* `project_id` - (Required, ForceNew) The ID of the DataWorks project (workspace). Changing this parameter will force a new resource to be created.

* `owner` - (Required) The ID of the Alibaba Cloud account used by the baseline owner.

* `priority` - (Required) The priority of the baseline. Valid values: 1, 3, 5, 7, 8.

* `baseline_type` - (Required) The type of the baseline. Valid values: `DAILY`, `HOURLY`, `WEEKLY`, `MONTHLY`.

* `overtime_settings` - (Required) The settings of the committed completion time of the baseline. Structure is documented below.

### Optional Parameters

* `alert_margin_threshold` - (Optional) The alert margin threshold of the baseline. Unit: minutes.

* `enabled` - (Optional) Whether the baseline is enabled.

* `alert_enabled` - (Optional) Whether alert is enabled for the baseline.

---

### Nested Blocks

#### `overtime_settings`

The `overtime_settings` block supports:

* `cycle` - (Optional) The cycle number. For a day-level baseline, set this parameter to 1. For an hour-level baseline, set this parameter to a value that is no more than 24.

* `time` - (Optional) The committed completion time in the `hh:mm` format. Valid values of `hh`: [0, 47]. Valid values of `mm`: [0, 59].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource. The format is `<project_id>:<baseline_id>`.

* `baseline_id` - The ID of the created baseline.

## Import

DataWorks Baseline can be imported using the `project_id` and `baseline_id` separated by a colon, e.g.

```bash
$ terraform import alibabacloudstack_data_works_baseline.example 10000:100003
```
