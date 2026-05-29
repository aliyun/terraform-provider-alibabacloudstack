---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_alarmtask"
sidebar_current: "docs-Alibabacloudstack-autoscaling-alarmtask"
description: |- 
  Provides a autoscaling Alarmtask resource.
---

# alibabacloudstack_autoscaling_alarmtask
-> **NOTE:** Alias name has: `alibabacloudstack_ess_alarm`

Provides a autoscaling Alarmtask resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_images" "ecs_image" {
  most_recent = true
  name_regex  = "^centos_6\\w{1,5}[64].*"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alibabacloudstack_vpc" "foo" {
  name       = "tf-testAccEssAlarm_basic"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "foo" {
  vswitch_name      = "tf-testAccEssAlarm_basic_foo"
  vpc_id            = alibabacloudstack_vpc.foo.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_vswitch" "bar" {
  vswitch_name      = "tf-testAccEssAlarm_basic_bar"
  vpc_id            = alibabacloudstack_vpc.foo.id
  cidr_block        = "172.16.1.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_autoscaling_group" "foo" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "tf-testAccEssAlarm_basic"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vswitch.foo.id, alibabacloudstack_vswitch.bar.id]
}

resource "alibabacloudstack_autoscaling_rule" "foo" {
  scaling_rule_name = "tf-testAccEssAlarm_basic"
  scaling_group_id  = alibabacloudstack_autoscaling_group.foo.id
  adjustment_type   = "TotalCapacity"
  adjustment_value  = 2
  cooldown          = 60
}

resource "alibabacloudstack_autoscaling_alarmtask" "foo" {
  alarm_task_name   = "tf-testAccEssAlarm_basic"
  description       = "Acc alarm test"
  scaling_group_id  = alibabacloudstack_autoscaling_group.foo.id
  alarm_actions     = [alibabacloudstack_autoscaling_rule.foo.scaling_rule_arn]
  metric_type       = "system"
  metric_name       = "CpuUtilization"
  period            = 300
  statistics        = "Average"
  threshold         = "200.3"
  comparison_operator = ">="
  evaluation_count  = 2
}
```

## Argument Reference

The following arguments are supported:

* `alarm_actions` - (Required) List of actions to execute when alarm is triggered. Max 5 items, min 1 item.
* `scaling_group_id` - (Required, ForceNew) The ID of the scaling group associated with the alarm task.
* `metric_name` - (Required) The name for the alarm's associated metric. See [Block_metricNames_and_dimensions](#block-metricnames_and_dimensions) below for details.
* `threshold` - (Required) The value against which the specified statistics is compared.
* `alarm_task_name` - (Optional) The name of the alarm task.
* `description` - (Optional) Description of the alarm task.
* `status` - (Optional) The status of the alarm task.
* `metric_type` - (Optional, ForceNew) The type for the alarm's associated metric. Supported values: `system`, `custom`. `"system"` means the metric data is collected by Aliyun Cloud Monitor Service(CMS), `"custom"` means the metric data is uploaded to CMS by users. Defaults to `system`.
* `period` - (Optional, ForceNew) The period in seconds over which the specified statistic is applied. Supported values: `60`, `120`, `300`, `900`. Defaults to `300`.
* `statistics` - (Optional) The statistic to apply to the alarm's associated metric. Supported values: `Average`, `Minimum`, `Maximum`. Defaults to `Average`.
* `comparison_operator` - (Optional) The arithmetic operation to use when comparing the specified Statistic and Threshold. The specified Statistic value is used as the first operand. Supported values: `>=`, `<=`, `>`, `<`. Defaults to `>=`.
* `evaluation_count` - (Optional) The number of times that needs to satisfy the comparison condition before transitioning into ALARM state. Defaults to `3`.

### Deprecated Arguments

> **WARNING:** The following arguments are deprecated and will be removed in a future release.

* `name` - (Deprecated) Use `alarm_task_name` instead. Conflicts with `alarm_task_name`.
* `enable` - (Deprecated) Use `status` instead. Conflicts with `status`.
* `cloud_monitor_group_id` - (Deprecated, Remove in 3.21.0) Defines the application group ID defined by CMS which is assigned when you upload custom metrics to CMS, only available for custom metrics.
* `dimensions` - (Deprecated, Remove in 3.21.0) The dimension map for the alarm's associated metric. For all metrics, you cannot set the dimension key as `scaling_group` or `userId`, which is set by default.

### Block metricNames_and_dimensions

Supported metric names and dimensions:

| MetricName         | Dimensions                   |
|--------------------|-----------------------------|
| CpuUtilization     | user_id,scaling_group      |
| ClassicInternetRx  | user_id,scaling_group      |
| ClassicInternetTx  | user_id,scaling_group      |
| VpcInternetRx      | user_id,scaling_group      |
| VpcInternetTx      | user_id,scaling_group      |
| IntranetRx         | user_id,scaling_group      |
| IntranetTx         | user_id,scaling_group      |
| LoadAverage        | user_id,scaling_group      |
| MemoryUtilization  | user_id,scaling_group      |
| SystemDiskReadBps  | user_id,scaling_group      |
| SystemDiskWriteBps | user_id,scaling_group      |
| SystemDiskReadOps  | user_id,scaling_group      |
| SystemDiskWriteOps | user_id,scaling_group      |
| PackagesNetIn      | user_id,scaling_group,device |
| PackagesNetOut     | user_id,scaling_group,device |
| TcpConnection      | user_id,scaling_group,state |

**NOTE:** Dimensions `user_id` and `scaling_group` are automatically filled, which means you only need to care about dimensions `device` and `state` when needed.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the alarm task.
* `alarm_task_name` - The name of the alarm task.
* `status` - The status of the alarm task.
* `alarm_trigger_state` - The trigger status of the alarm task. Possible values:
  * `ALARM`: ALARM, the ALARM condition has been met.
  * `OK`: Normal, the alarm condition has not been met.
  * `INSUFFICIENT_DATA`: Insufficient data to determine whether the alarm condition is met.

### Deprecated Attributes

> **WARNING:** The following attributes are deprecated and will be removed in a future release.

* `state` - (Deprecated) Use `alarm_trigger_state` instead.
* `dimensions` - (Deprecated) The dimensions for the alarm's associated metric.

## Import

ESS Alarm can be imported using the alarm task ID, e.g.

```
$ terraform import alibabacloudstack_ess_alarm.example at-12345678
```