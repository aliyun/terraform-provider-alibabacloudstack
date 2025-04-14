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

## 参数说明

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多AZ配置中使用。默认值为`false`。如果设置为`true`，则返回的可用区列表将仅包含支持多AZ部署的可用区。多AZ通常用于启动高可用的MongoDB实例。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 可用区ID列表。此列表按顺序包含了所有符合条件的可用区ID。
* `zones` - 可用区详细信息列表。每个元素包含以下属性：
  * `id` - 可用区的唯一标识符。
  * `multi_zone_ids` - 如果该可用区支持多AZ部署，则此字段将包含与之关联的其他可用区ID列表。如果不支持多AZ部署，则此字段可能为空或不返回。