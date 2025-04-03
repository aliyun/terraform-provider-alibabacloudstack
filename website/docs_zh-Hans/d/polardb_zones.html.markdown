---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_zones"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-zones"
description: |-
  查询polardb可用区列表。
---

# alibabacloudstack_polardb_zones

根据指定过滤条件列出当前凭证权限可以访问的polardb可用区列表。

## 示例用法

```hcl
data "alibabacloudstack_polardb_zones" "example" {
  multi = false
}

output "zones" {
  value = data.alibabacloudstack_polardb_zones.example.zones
}
```

## 参数参考
支持以下参数：

* `multi` - (可选) 是否检索多可用区 ID。默认为 false。

## 属性参考
导出以下属性：

* `ids` - 可用区的 ID 列表。
* `zones` - 可用区列表。每个元素包含以下属性：
    * `id` - 可用区的唯一标识符。
    * `multi_zone_ids` - 多可用区 ID 列表(仅在 multi 设置为 true 时可用)。