---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_zones"
sidebar_current: "docs-alibabacloudstack-datasource-db-zones"
description: |-
    查询RDS数据库服务可用区
---

# alibabacloudstack_db_zones

根据指定过滤条件列出当前凭证权限可以访问的RDS数据库服务可用区列表。


## 示例用法

```
# 声明数据源
data "alibabacloudstack_db_zones" "zones_ids" {}

output "db_zones" {
  value = data.alibabacloudstack_db_zones.zones_ids.zones.*
}

```

## 参数参考

支持以下参数：

* `multi` - (可选) 指定可用区是否可用于多可用区配置，默认为 `false`。 多可用区通常用于启动RDS实例。

* `ids` - (可选) 可用区ID列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区ID。
  * `multi_zone_ids` - 多可用区配置中可用的区ID列表。

* `multi` - 表示这些可用区是否可用于多可用区配置。