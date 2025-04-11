---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_autosnapshotpolicy"
sidebar_current: "docs-Alibabacloudstack-ecs-autosnapshotpolicy"
description: |- 
  编排云服务器（Ecs）自动快照策略
---

# alibabacloudstack_ecs_autosnapshotpolicy
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_snapshot_policy`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）自动快照策略。

## 示例用法

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

## 参数说明

支持以下参数：

* `auto_snapshot_policy_name` - (可选) 自动快照策略的名称。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(`:`)、下划线(`_`)或者短划线(`-`)。默认值为空。
* `repeat_weekdays` - (必填) 指定自动快照的重复日期。选定周一到周日中需要创建快照的日期，参数为1~7的数字，如：`1`表示周一，`7`表示周日。允许选择多个日期。格式为JSON数组，例如：`["1", "2", "3"]`。
* `retention_days` - (可选) 自动快照的保留时间，单位为天。取值范围：
  - `-1`：永久保存。
  - `1`~`65536`：指定保存天数。默认值为`-1`。
* `time_points` - (必填) 指定自动快照的创建时间点。最小单位为小时，从`00:00`~`23:00`共24个时间点可选，参数为`0`~`23`的数字，如：`1`代表在`01:00`时间点。可以选定多个时间点。传递参数为一个带有格式的Json Array，例如：`["0", "1", ..., "23"]`，最多24个时间点，用半角逗号字符隔开。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 自动快照策略的ID。
* `name` - 自动快照策略的名称。
* `auto_snapshot_policy_name` - 自动快照策略的名称。该属性与`auto_snapshot_policy_name`参数相同。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(`:`)、下划线(`_`)或者短划线(`-`)。默认值为空。