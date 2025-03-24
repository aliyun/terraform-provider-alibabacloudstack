---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_remind"
sidebar_current: "docs-Alibabacloudstack-data-works-remind"
description: |- 
  Provides a data works Remind resource.
---

# alibabacloudstack_data_works_remind

Provides a data works Remind resource.

For information about Data Works Remind and how to use it, see [What is Remind](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/CreateRemind-1-2.html?spm=a2c4g.14484438.10001.638).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccdataworksremind93940"
}

resource "alibabacloudstack_data_works_remind" "default" {
  alert_methods = "SMS"
  alert_unit = "OWNER"
  remind_name = var.name
  remind_type = "FINISHED"
  remind_unit = "PROJECT"
  project_id = "10023"
  dnd_end = "23:59"
  node_ids = "node1,node2"
  baseline_ids = "baseline1,baseline2"
  biz_process_ids = "bizprocess1,bizprocess2"
  max_alert_times = 5
  alert_interval = 1200
  detail = "{\"hour\":23,\"minu\":59}"
  alert_targets = "uid1,uid2"
  robot_urls = "https://robot.url1,https://robot.url2"
  use_flag = true
}
```

## Argument Reference

The following arguments are supported:

* `alert_unit` - (Required) The granularity of the receiving object, including:
  * `OWNER`: Task owner
  * `OTHER`: Designated person

* `remind_name` - (Required) The name of the custom rule cannot exceed 128 characters.

* `remind_type` - (Required) Trigger conditions, including:
  * `FINISHED`: Task finished
  * `UNFINISHED`: Task not finished
  * `ERROR`: Task error
  * `CYCLE_UNFINISHED`: Cycle task not finished
  * `TIMEOUT`: Task timeout

* `remind_unit` - (Required) Types of objects, including:
  * `NODE`: Task node
  * `BASELINE`: Baseline
  * `PROJECT`: Workspace
  * `BIZPROCESS`: Business process

* `dnd_end` - (Optional) Do not disturb deadline, format HH:MM. The value range of hh is 0-23, and the value range of mm is 0-59. Default is `00:00`.

* `node_ids` - (Optional) The monitored task node id when the object type (`remind_unit`) is `NODE`. Multiple IDs are separated by commas (`,`), and a rule can monitor up to 50 nodes.

* `baseline_ids` - (Optional) The monitored baseline id when the object type (`remind_unit`) is `BASELINE`. Multiple IDs are separated by commas (`,`), and one rule can monitor up to 5 baselines.

* `project_id` - (Optional) The monitored workspace id when the object type (`remind_unit`) is `PROJECT`. A rule can only monitor one workspace.

* `biz_process_ids` - (Optional) The monitored business process id when the object type (`remind_unit`) is `BIZPROCESS`. Multiple business process ids are separated by commas (`,`), and a rule can monitor up to 5 business processes.

* `max_alert_times` - (Optional) Maximum number of alarms. The minimum value is 1, the maximum value is 10, and the default value is 3.

* `alert_interval` - (Optional) Minimum alarm interval, in seconds. The minimum value is 1200, and the default value is 1800.

* `detail` - (Optional) The descriptions of different trigger conditions are as follows:
  * When the `remind_type` is `FINISHED`, it will be blank.
  * When the `remind_type` is `UNFINISHED`, the parameter format is `{"hour":23,"minu":59}`. The value range of hour is 0-47, and that of minu is 0-59.
  * When the `remind_type` is `ERROR`, it is passed blank.
  * When the `remind_type` is `CYCLE_UNFINISHED`, the format of parameter passing is `{"1": "05:50", "2": "06:50", ...}`. The string JSON key is the period number, and its value range is 1-288. Value is the unfinished time corresponding to this cycle, and the format is HH:mm. The value range of hh is 0-47, and the value range of mm is 0-59.
  * When the `remind_type` is `TIMEOUT`, the parameter format is 1800 in seconds. That is, from the start of the operation, running for more than 30 minutes will trigger an alarm.

* `alert_methods` - (Optional) The alarm methods include:
  * `MAIL`: Email
  * `SMS`: Short message
  * `PHONE`: Telephone (only supported by DataWorks Professional and above)

  Multiple alarm modes are separated by English commas (`,`).

* `alert_targets` - (Optional)
  * When the `alert_unit` is `OWNER` (node task owner), it is blank.
  * When the `alert_unit` is `OTHER`, the Alibaba Cloud UID of the specified user is passed in. Multiple Alibaba Cloud UIDs are separated by English commas (`,`), and the maximum number is 10.

* `robot_urls` - (Optional) The webhook addresses of DingTalk robots, and multiple webhook addresses are separated by English commas (`,`).

* `use_flag` - (Optional) Open and close rules, including:
  * `true`: Enable
  * `false`: Disable

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `remind_id` - The resource ID of Remind. The value formats as `<remind_id>:<$.ProjectId>`.