---
subcategory: "Auto Scaling(ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_alarm"
sidebar_current: "docs-alibabacloudstack-resource-ess-alarm"
description: |-
    Provides a ESS alarm task resource.
---

# alibabacloudstack\_ess\_alarm

Provides a ESS alarm task resource.

## Example Usage
```
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

resource "alibabacloudstack_ess_scaling_group" "foo" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "tf-testAccEssAlarm_basic"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vswitch.foo.id, alibabacloudstack_vswitch.bar.id]
}

resource "alibabacloudstack_ess_scaling_rule" "foo" {
  scaling_rule_name = "tf-testAccEssAlarm_basic"
  scaling_group_id  = alibabacloudstack_ess_scaling_group.foo.id
  adjustment_type   = "TotalCapacity"
  adjustment_value  = 2
  cooldown          = 60
}

resource "alibabacloudstack_ess_alarm" "foo" {
  name                = "tf-testAccEssAlarm_basic"
  description         = "Acc alarm test"
  alarm_actions       = [alibabacloudstack_ess_scaling_rule.foo.ari]
  scaling_group_id    = alibabacloudstack_ess_scaling_group.foo.id
  metric_type         = "system"
  metric_name         = "CpuUtilization"
  period              = 300
  statistics          = "Average"
  threshold           = 200.3
  comparison_operator = ">="
  evaluation_count    = 2
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name for ess alarm.
* `description` - (Optional) The description for the alarm.
* `enable` - (Optional) Whether to enable specific ess alarm. Default to true.
* `alarm_actions` - (Required) The list of actions to execute when this alarm transition into an ALARM state. Each action is specified as ess scaling rule ari.
* `scaling_group_id` - (Required, ForceNew) The scaling group associated with this alarm.
* `metric_type` - (Optional, ForceNew) The type for the alarm's associated metric. Supported value: system, custom. "system" means the metric data is collected by Aliyun Cloud Monitor Service(CMS), "custom" means the metric data is upload to CMS by users. Defaults to system.
* `metric_name` - (Required) The name for the alarm's associated metric. See [Block_metricNames_and_dimensions](#block-metricnames_and_dimensions) below for details.
* `period` - (Optional, ForceNew) The period in seconds over which the specified statistic is applied. Supported value: 60, 120, 300, 900. Defaults to 300.
* `statistics` - (Optional) The statistic to apply to the alarm's associated metric. Supported value: Average, Minimum, Maximum. Defaults to Average.
* `threshold` - (Required) The value against which the specified statistics is compared.
* `comparison_operator` - (Optional) The arithmetic operation to use when comparing the specified Statistic and Threshold. The specified Statistic value is used as the first operand. Supported value: >=, <=, >, <. Defaults to >=.
* `evaluation_count` - (Optional) The number of times that needs to satisfies comparison condition before transition into ALARM state. Defaults to 3.
* `cloud_monitor_group_id` - (Optional) Defines the application group id defined by CMS which is assigned when you upload custom metric to CMS, only available for custom metirc.
* `dimensions` - (Optional) The dimension map for the alarm's associated metric (documented below). For all metrics, you can not set the dimension key as "scaling_group" or "userId", which is set by default, the second dimension for metric, such as "device" for "PackagesNetIn", need to be set by users.

## Block metricNames_and_dimensions

Supported metric names and dimensions :

| MetricName         | Dimensions                   |
| ------------------ | ---------------------------- |
| CpuUtilization     | user_id,scaling_group        |
| ClassicInternetRx  | user_id,scaling_group        |
| ClassicInternetTx  | user_id,scaling_group        |
| VpcInternetRx      | user_id,scaling_group        |
| VpcInternetTx      | user_id,scaling_group        |
| IntranetRx         | user_id,scaling_group        |
| IntranetTx         | user_id,scaling_group        |
| LoadAverage        | user_id,scaling_group        |
| MemoryUtilization  | user_id,scaling_group        |
| SystemDiskReadBps  | user_id,scaling_group        |
| SystemDiskWriteBps | user_id,scaling_group        |
| SystemDiskReadOps  | user_id,scaling_group        |
| SystemDiskWriteOps | user_id,scaling_group        |
| PackagesNetIn      | user_id,scaling_group,device |
| PackagesNetOut     | user_id,scaling_group,device |
| TcpConnection      | user_id,scaling_group,state  |

-> **NOTE:** Dimension `user_id` and `scaling_group` is automatically filled, which means you only need to care about dimension `device` and `state` when needed.

## Attribute Reference

The following attributes are exported:

* `id` - The id for ess alarm.
* `state` - The state of specified alarm.

## Import

Ess alarm can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_ess_alarm.example asg-2ze500_045efffe-4d05
```
