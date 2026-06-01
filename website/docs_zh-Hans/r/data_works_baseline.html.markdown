---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_baseline"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-baseline"
description: |-
  提供 DataWorks 基线资源。
---

# alibabacloudstack_data_works_baseline

提供 DataWorks 基线资源。

## 示例

基本用法

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

## 参数说明

以下参数适用于此资源：

### 必填参数

* `baseline_name` - （必填）基线名称。

* `project_id` - （必填，ForceNew）DataWorks 工作空间（项目）的 ID。更改此参数将强制创建新资源。

* `owner` - （必填）基线责任人的阿里云 UID。

* `priority` - （必填）基线的优先级。有效值：1、3、5、7、8。

* `baseline_type` - （必填）基线类型。有效值：`DAILY`（天基线）、`HOURLY`（小时基线）、`WEEKLY`（周基线）、`MONTHLY`（月基线）。

* `overtime_settings` - （必填）基线承诺时间配置。结构见下文。

### 可选参数

* `alert_margin_threshold` - （可选）基线预警余量。单位：分钟。

* `enabled` - （可选）是否启用基线。

* `alert_enabled` - （可选）是否启用基线告警。

---

### 嵌套块

#### `overtime_settings`

`overtime_settings` 块支持以下参数：

* `cycle` - （可选）周期编号。对于天基线，将此参数设置为 1；对于小时基线，将此参数设置为不超过 24 的值。

* `time` - （可选）承诺时间，格式为 `hh:mm`。`hh` 的有效值范围：[0, 47]。`mm` 的有效值范围：[0, 59]。

## 属性说明

除上述所有参数外，还导出以下属性：

* `id` - 资源的 ID。格式为 `<project_id>:<baseline_id>`。

* `baseline_id` - 创建的基线 ID。

## 导入

DataWorks 基线可以使用 `project_id` 和 `baseline_id`（以冒号分隔）导入，例如：

```bash
$ terraform import alibabacloudstack_data_works_baseline.example 10000:100003
```
