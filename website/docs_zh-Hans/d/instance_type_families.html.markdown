---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance_type_families"
sidebar_current: "docs-alibabacloudstack-datasource-instance-type-families"
description: |-
  查询云服务器实例类型族
---

# alibabacloudstack_instance_type_families

根据指定过滤条件列出当前凭证权限可以访问的云服务器实例类型族列表。

## 示例用法

```
data "alibabacloudstack_instance_type_families" "default" {
  
}

output "first_instance_type_family_id" {
  value = "${data.alibabacloudstack_instance_type_families.default.instance_type_families.0.id}"
}

output "instance_ids" {
  value = "${data.alibabacloudstack_instance_type_families.default.ids}"
}
```

## 参数说明

支持以下参数：

* `zone_id` - (可选，变更时重建) 指定启动实例所在的可用区。如果未指定，则返回所有可用区的实例类型族。
* `generation` - (可选) 指定实例类型族的代数。例如，可以选择特定一代或新一代的实例类型族。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 包含所有匹配条件的实例类型族ID的列表。
* `families` - 实例类型族的详细信息列表。每个元素包含以下属性：
  * `id` - 实例类型族的唯一标识符。
  * `generation` - 实例类型族所属的代数。
  * `zone_ids` - 支持该实例类型族的可用区列表。