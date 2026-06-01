---
subcategory: "弹性伸缩 ESS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduled_task"
sidebar_current: "docs-Alibabacloudstack-resource-autoscaling-scheduled-task"
description: |-
  提供弹性伸缩定时任务资源。
---

# alibabacloudstack_autoscaling_scheduled_task

提供弹性伸缩定时任务资源，允许您创建、修改和删除伸缩组的定时任务。

## 示例用法

### 基础用法

```hcl
resource "alibabacloudstack_autoscaling_scheduled_task" "default" {
  scheduled_action       = "acs:ess:::scalingGroup/${alibabacloudstack_autoscaling_scaling_group.default.id}:changeInCapacity:1"
  launch_time            = "2026-06-01T15:04Z"
  scheduled_task_name    = "tf-testAccScheduledTask"
  description            = "测试定时任务"
  launch_expiration_time = 600
  recurrence_type        = "Daily"
  recurrence_value       = "1"
  task_enabled           = true
  scaling_group_id       = alibabacloudstack_autoscaling_scaling_group.default.id
}
```

## 参数说明

支持以下参数：

* `scheduled_action` - (必填) 定时任务触发时执行的操作。格式：`acs:ess:::scalingGroup/<scalingGroupId>:changeInCapacity:<value>`。
* `launch_time` - (必填) 定时任务的触发时间。格式：`YYYY-MM-DDThh:mmZ`（UTC时间）。
* `scheduled_task_name` - (可选) 定时任务名称。长度为2-64个字符，以字母或数字开头，可包含字母、数字、下划线(_)、连字符(-)和句点(.)。
* `description` - (可选, 由API返回) 定时任务描述。长度为2-200个字符。
* `launch_expiration_time` - (可选) 定时任务触发后的有效期限。单位：秒。取值范围：0-21600。默认值：600。
* `recurrence_type` - (可选, 由API返回) 定时任务的重复类型。取值：`Daily`（每天）、`Weekly`（每周）、`Monthly`（每月）。
* `recurrence_value` - (可选, 由API返回) 重复取值。格式因重复类型而异：
  - Daily：正整数，表示执行间隔。
  - Weekly：用逗号分隔的星期几（0-6，其中0表示星期日）。
  - Monthly：用逗号分隔的月份中的日期（1-31）。
* `recurrence_end_time` - (可选, 由API返回) 重复定时任务的结束时间。格式：`YYYY-MM-DDThh:mmZ`（UTC时间）。
* `task_enabled` - (可选) 是否启用定时任务。取值：`true`、`false`。默认值：`true`。
* `scaling_group_id` - (可选, 由API返回) 定时任务所属的伸缩组ID。

## 属性说明

导出以下属性：

* `id` - 定时任务ID。

## Import

弹性伸缩定时任务可以通过定时任务ID导入，例如：

```
$ terraform import alibabacloudstack_autoscaling_scheduled_task.example st-12345678
```
