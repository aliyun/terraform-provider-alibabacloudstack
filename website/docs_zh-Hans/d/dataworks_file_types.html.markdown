---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dataworks_file_types"
sidebar_current: "docs-Alibabacloudstack-datasource-dataworks-file-types"
description: |-
  提供DataWorks文件类型列表。
---

# alibabacloudstack_dataworks_file_types

该数据源用于获取专有云中可用的DataWorks文件类型列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_dataworks_file_types" "example" {
  project_id = 12345
  name       = "Shell"
}

output "file_types" {
  value = data.alibabacloudstack_dataworks_file_types.example.file_types
}
```

## 参数说明

以下参数支持配置：

* `project_id` - (必选) DataWorks项目的ID。
* `name` - (可选) 用于按名称过滤文件类型的关键词。
* `ids` - (可选, Computed) 用于过滤结果的文件类型ID列表。

## 属性参考

以下属性会被导出：

* `ids` - 文件类型ID列表。
* `names` - 文件类型名称列表。
* `file_types` - DataWorks文件类型列表。每个元素包含以下属性：
  * `node_type_name` - 节点类型的显示名称。
  * `node_type_id` - 节点类型的数字标识符。
