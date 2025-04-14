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
* `specification_type` - (必填) MaxCompute 项目的规格类型，例如 `OdpsStandard`。
* `order_type` - (必填) MaxCompute 项目的计费类型，例如 `PayAsYouGo`（按量付费）。

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