---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_project"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-project"
description: |-
  编排Max Compute项目
---

# alibabacloudstack_maxcompute_project

使用Provider配置的凭证在指定的资源集编排Max Compute项目


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_maxcompute_project" "example" {
  project_name       = "tf_maxcompute_project"
  specification_type = "OdpsStandard"
  order_type         = "PayAsYouGo"
}
```

## 参数说明

以下参数被支持：
* `project_name` - (必填，变更时重建) MaxCompute 项目的名称。
* `quota_id` - (必填) MaxCompute 项目的配额 ID。
* `disk` - (必填) MaxCompute 项目的磁盘大小。
* `account` - （必填，强制新建）MaxCompute 项目的账户。
* `account_pk` - （必填，强制新建）MaxCompute 项目的账户主键。
* `external_table` - （可选）是否启用联合计算。
* `vpc_ids` - （可选）MaxCompute 项目的 VPC ID 列表。
* `core_arch` - （可选）MaxCompute 项目的核心架构。
* `cpu_type` - （可选）MaxCompute 项目的 CPU 类型。

## 属性说明

以下属性会被导出：

* `id` - MaxCompute 项目的唯一标识符。它与 `project_name` 相同。
* `project_name` - MaxCompute 项目的名称。
* `quota_id` - MaxCompute 项目的配额 ID。
* `disk` - MaxCompute 项目的磁盘大小。

## 导入

MaxCompute 项目可以使用 *名称* 或 ID 导入，例如

```bash
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```