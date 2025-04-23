---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_alarmtask"
sidebar_current: "docs-Alibabacloudstack-autoscaling-alarmtask"
description: |- 
  编排弹性伸缩的告警任务
---

# alibabacloudstack_autoscaling_alarmtask
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_alarm`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩的告警任务。

## 示例用法

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
  metric_type       = "system"
  metric_name       = "CpuUtilization"
  period            = 300
  statistics        = "Average"
  threshold         = 200.3
  comparison_operator = ">="
  evaluation_count  = 2
}
```

## 参数参考

支持以下参数：

* `name` - (选填) - 报警任务的名称，与 `alarm_task_name` 功能相同。
* `alarm_task_name` - (选填) - 报警任务的名称。
* `description` - (选填) - 报警任务的描述。
* `enable` - (选填) - 是否启用特定的报警任务。默认为 `true`。
* `status` - (选填) - 报警任务的状态。
* `scaling_group_id` - (必填, 变更时重建) - 报警任务关联的伸缩组的 ID。
* `metric_type` - (选填, 变更时重建) - 监控项类型。支持值：`system`、`custom`。`"system"` 表示指标数据由阿里云监控服务(CMS)收集，`"custom"` 表示指标数据由用户上传到 CMS。默认为 `system`。
* `metric_name` - (必填) - 监控项名称。详见下方 [Block_metricNames_and_dimensions](#block-metricnames_and_dimensions)。
* `period` - (选填, 变更时重建) - 指定统计量应用的时间段(秒)。支持值：`60`、`120`、`300`、`900`。默认为 `300`。
* `statistics` - (选填) - 要应用于报警关联指标的统计量。支持值：`Average`、`Minimum`、`Maximum`。默认为 `Average`。
* `threshold` - (必填) - 要比较的指定统计量的值。
* `comparison_operator` - (选填) - 比较指定统计量和阈值时使用的算术运算符。指定的统计量值用作第一个操作数。支持值：`>=`、`<=`、`>`、`<`。默认为 `>=`。
* `evaluation_count` - (选填) - 在进入 ALARM 状态之前需要满足比较条件的次数。默认为 `3`。
* `cloud_monitor_group_id` - (选填) - CMS 定义的应用程序组 ID，在您将自定义指标上传到 CMS 时分配，仅适用于自定义指标。
* `dimensions` - (选填) - 报警关联指标的维度映射。对于所有指标，不能将维度键设置为 `scaling_group` 或 `userId`，这是默认设置的。某些指标的第二个维度(如 `PackagesNetIn` 的 `device`)需要由用户设置。

### Block metricNames_and_dimensions

支持的指标名称和维度：

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

**注意:** 维度 `user_id` 和 `scaling_group` 是自动填充的，这意味着您只需要关心在需要时设置维度 `device` 和 `state`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `name` - 报警任务的名称，与 `alarm_task_name` 功能相同。
* `alarm_task_name` - 报警任务的名称。
* `enable` - 是否启用特定的报警任务。
* `status` - 报警任务的状态。
* `dimensions` - 报警关联指标的维度。
* `state` - 指定报警的状态。
* `alarm_trigger_state` - 报警任务的触发状态。可能值：
  * `ALARM`: 报警，已满足 ALARM 条件。
  * `OK`: 正常，未满足报警条件。
  * `INSUFFICIENT_DATA`: 数据不足，无法确定是否满足报警条件。
