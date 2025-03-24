---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_autosnapshotpolicy"
sidebar_current: "docs-Alibabacloudstack-ecs-autosnapshotpolicy"
description: |- 
  Provides a ecs Autosnapshotpolicy resource.
---

# alibabacloudstack_ecs_autosnapshotpolicy
-> **NOTE:** Alias name has: `alibabacloudstack_snapshot_policy`

Provides a ecs Autosnapshotpolicy resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccecsauto_snapshot_policy44210"
}

resource "alibabacloudstack_ecs_autosnapshotpolicy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["1", "2", "3"]
  retention_days            = -1
  time_points               = ["1", "22", "23"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_name` - (Optional) The name of the automatic snapshot policy. The name must be 2 to 128 characters in length. The name must start with a letter and cannot start with `http://` or `https://`. The name can contain letters, digits, colons (`:`), underscores (`_`), and hyphens (`-`). By default, this parameter is left empty.
* `repeat_weekdays` - (Required) The days of the week on which the automatic snapshots are created. Valid values are `"1"` to `"7"`, where `"1"` represents Monday and `"7"` represents Sunday. You can specify up to seven days. The format is a JSON array, such as `["1", "2", "3"]`.
* `retention_days` - (Optional) The retention period of the automatic snapshots. Unit: days. Valid values:
  - `-1`: The snapshots are retained permanently until manually deleted.
  - `1` to `65535`: The number of days the snapshots are retained. After the retention period expires, the snapshots are automatically deleted.
  Default value: `-1`.
* `time_points` - (Required) The time points at which the automatic snapshots are created. Valid values are `"0"` to `"23"`, representing the hours from `00:00` to `23:00`. You can specify up to 24 time points. The format is a JSON array, such as `["1", "22", "23"]`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the automatic snapshot policy.
* `name` - The name of the automatic snapshot policy.
* `auto_snapshot_policy_name` - The name of the automatic snapshot policy. This attribute mirrors the `auto_snapshot_policy_name` argument.
