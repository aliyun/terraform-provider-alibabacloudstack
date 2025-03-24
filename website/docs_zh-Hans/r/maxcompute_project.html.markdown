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

## 参数参考

以下参数被支持：
* `name` - (必填，变更时重建) 已从提供程序版本 1.110.0 开始弃用，取而代之的是 `project_name`。
* `quota_id` - (必填) MaxCompute 项目的配额 ID。
* `disk` - (必填) MaxCompute 项目的磁盘大小。

## 属性参考

以下属性会被导出：

* `id` - MaxCompute 项目的 ID。它与名称相同。
* `name` - MaxCompute 项目的名称。

## 导入

MaxCompute 项目可以使用 *名称* 或 ID 导入，例如

```bash
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```