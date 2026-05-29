---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_remind"
sidebar_current: "docs-Alibabacloudstack-data-works-remind"
description: |-
  Provides a DataWorks Remind resource.
---

# alibabacloudstack_data_works_remind

Provides a DataWorks Remind (custom alert rule) resource.

-> **Note:** The `node_ids`, `baseline_ids`, `project_id`, and `biz_process_ids` parameters are mutually exclusive. At least one of them must be specified based on the `remind_unit` type.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_remind" "default" {
  remind_name   = "tf_test_remind"
  remind_type   = "ERROR"
  remind_unit   = "NODE"
  node_ids      = ["12345", "12346"]
  alert_unit    = "OWNER"
  alert_methods = ["MAIL", "SMS"]
}
```

## Argument Reference

The following arguments are supported:

### Required Parameters

* `remind_name` - (Required) The name of the custom alert rule. The name cannot exceed 128 characters in length.

* `remind_type` - (Required) The conditions that trigger an alert. Valid values: `FINISHED` (completed), `UNFINISHED` (not completed), `ERROR` (error), `CYCLE_UNFINISHED` (cycle not completed), `TIMEOUT` (timeout).

* `remind_unit` - (Required) The type of the object to which the custom alert rule is applied. Valid values: `NODE` (node), `BASELINE` (baseline), `PROJECT` (workspace), `BIZPROCESS` (workflow).

* `alert_unit` - (Required) The recipient of the alert. Valid values: `OWNER` (node owner) and `OTHER` (specified user).

* `alert_methods` - (Required) The notification methods. Valid values: `MAIL`, `SMS`, `PHONE`, `WEBHOOKS`, `DINGROBOTS`.

### Optional Parameters

* `node_ids` - (Optional) The IDs of the nodes to which the custom alert rule is applied. This parameter takes effect when `remind_unit` is set to `NODE`. You can specify multiple IDs separated by commas. A maximum of 50 nodes can be specified. Conflicts with `baseline_ids`, `project_id`, and `biz_process_ids`.

* `baseline_ids` - (Optional) The IDs of the baselines to which the custom alert rule is applied. This parameter takes effect when `remind_unit` is set to `BASELINE`. You can specify multiple IDs separated by commas. A maximum of 5 baselines can be specified. Conflicts with `node_ids`, `project_id`, and `biz_process_ids`.

* `project_id` - (Optional) The ID of the workspace to which the custom alert rule is applied. This parameter takes effect when `remind_unit` is set to `PROJECT`. Only one workspace can be specified. Conflicts with `node_ids`, `baseline_ids`, and `biz_process_ids`.

* `biz_process_ids` - (Optional) The IDs of the workflows to which the custom alert rule is applied. This parameter takes effect when `remind_unit` is set to `BIZPROCESS`. You can specify multiple IDs separated by commas. A maximum of 5 workflows can be specified. Conflicts with `node_ids`, `baseline_ids`, and `project_id`.

* `max_alert_times` - (Optional) The maximum number of alerts. Valid values: 1 to 10. Default value: 3.

* `alert_interval` - (Optional) The minimum interval at which alerts are reported. Unit: seconds. Minimum value: 1200. Default value: 1800.

* `detail` - (Optional) The configuration for different trigger conditions.
  - If `remind_type` is `FINISHED` or `ERROR`, leave this parameter empty.
  - If `remind_type` is `UNFINISHED`, configure as JSON string: `{"hour":23,"minu":59}`. Valid values of `hour`: [0,47]. Valid values of `minu`: [0,59].
  - If `remind_type` is `CYCLE_UNFINISHED`, configure as JSON string with cycle IDs as keys and timeout as values: `{"1":"05:50","2":"06:50"}`. Valid values of cycle ID: [1,288].
  - If `remind_type` is `TIMEOUT`, set the timeout period in seconds (e.g., `1800`).

* `alert_targets` - (Optional) The specified users' Alibaba Cloud UIDs when `alert_unit` is set to `OTHER`. Multiple UIDs are separated by commas. This parameter is computed from API response when `alert_unit` is `OWNER`.

* `webhooks` - (Optional) The webhook URLs for enterprise WeChat or Feishu robot notifications. Takes effect only when `WEBHOOKS` is included in `alert_methods`.

* `robot_urls` - (Optional) The DingTalk robot webhook URLs. Takes effect only when `DINGROBOTS` is included in `alert_methods`.

* `dnd_end` - (Optional) The end time of the quiet hours (do-not-disturb). Format: `hh:mm`. Valid values of `hh`: [0,23]. Valid values of `mm`: [0,59]. Default value: `00:00`.

* `use_flag` - (Optional) Whether the alert rule is enabled. Valid values: `true`, `false`. This attribute is computed from the API response.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the custom alert rule (remind_id).

* `remind_id` - The ID of the created remind rule.

## Import

DataWorks Remind can be imported using the `remind_id`, e.g.

```bash
$ terraform import alibabacloudstack_data_works_remind.example 12345
```
