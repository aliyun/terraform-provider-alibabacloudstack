---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-mongodb-zones"
description: |-
    查询MongoDB可用区
---

# alibabacloudstack_mongodb_zones

根据指定过滤条件列出当前凭证权限可以访问的MongoDB可用区列表。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_mongodb_zones" "zones_ids" {}

# 使用第一个匹配的可用区创建一个MongoDB实例
resource "alibabacloudstack_mongodb_instance" "mongodb" {
    zone_id = data.alibabacloudstack_mongodb_zones.zones_ids.zones[0].id

  # 其他属性...
}
```

## 参数参考

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多AZ配置中使用。默认为`false`。多AZ通常用于启动MongoDB实例。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区的ID。
  * `multi_zone_ids` - 多可用区中的可用区ID列表。