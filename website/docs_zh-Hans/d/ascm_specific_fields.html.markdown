---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_specific_fields"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-specific-fields"
description: |-
    查询特定字段
---

# alibabacloudstack_ascm_specific_fields

根据指定过滤条件列出当前凭证权限可以访问的特定字段列表。

## 示例用法

```
data "alibabacloudstack_ascm_specific_fields" "specifields" {
  group_filed ="storageType"
  resource_type ="OSS"
  output_file = "fields"
}
output "specifields" {
  value = data.alibabacloudstack_ascm_specific_fields.specifields.*
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 特定字段ID的列表。此参数可用于筛选特定字段。
* `group_filed` - (必填) 要查询有效值的字段。例如，您可以查询与存储类型相关的字段。
* `resource_type` - (必填) 通过指定的资源类型过滤结果。有效值包括：`OSS`, `ADB`, `DRDS`, `SLB`, `NAT`, `MAXCOMPUTE`, `POSTGRESQL`, `ECS`, `RDS`, `IPSIX`, `REDIS`, `MONGODB`, 和 `HITSDB`。
* `label` - (可选) 指定是否对字段进行国际化处理。如果设置为 `true`，则返回的字段名称将被翻译为国际化的格式；如果设置为 `false` 或未提供，则返回原始字段名称。有效值为：`true` 和 `false`。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `specific_fields` - 特定字段的列表。此属性包含符合查询条件的所有字段信息。每个字段可能包含以下信息：
  * `id` - 字段的唯一标识符。
  * `name` - 字段名称。
  * `value` - 字段的有效值或默认值。
  * `label` - 如果启用了国际化，此字段将显示翻译后的名称。
  * `description` - 字段的描述信息（如果存在）。
  * `type` - 字段的数据类型（如字符串、整数等）。
  * `required` - 字段是否为必填项。有效值为：`true`（必填）或 `false`（非必填）。
  * `default_value` - 字段的默认值（如果有）。
  * `options` - 字段的可选项列表（如果适用）。
