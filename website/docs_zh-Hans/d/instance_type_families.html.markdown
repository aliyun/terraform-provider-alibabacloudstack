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

## 参数参考

支持以下参数：

* `zone_id` - (可选，变更时重建) 启动实例的可用区。
* `generation` - (可选) 实例类型族的代数。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 实例类型族ID列表。
* `families` - 实例类型族列表。每个元素包含以下属性：
  * `id` - 实例类型族的ID。
  * `generation` - 实例类型族的代数。
  * `zone_ids` - 启动实例的可用区列表。
* `families.zone_ids` - 启动实例的可用区列表。 