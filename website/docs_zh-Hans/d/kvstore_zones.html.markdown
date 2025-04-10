---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_zones"
sidebar_current: "docs-alibabacloudstack-datasource-kvstore-zones"
description: |-
   查询KVStore可用区
---

# alibabacloudstack_kvstore_zones

根据指定过滤条件列出当前凭证权限可以访问的KVStore可用区列表。


## 示例用法

```
# 声明数据源
data "alibabacloudstack_kvstore_zones" "zones_ids" {}

output "kvstore_zones" {
  value = "${data.alibabacloudstack_kvstore_zones.zones_ids.zones}"
}
```

## 参数说明

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多AZ配置中使用。默认值为`false`。多AZ通常用于启动KVStore实例。
* `instance_charge_type` - (可选) 通过特定的实例计费类型过滤结果。有效值为`PrePaid`（包年包月）和`PostPaid`（按量付费）。默认值为`PostPaid`。
* `ids` - (可选) 可用区ID列表，用于精确匹配指定的可用区。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配的可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区的唯一标识符。
  * `multi_zone_ids` - 当`multi`参数设置为`true`时，返回的多可用区配置中的可用区ID列表。
  * `zones` - 可用区的详细信息列表，包含与该可用区相关的其他元数据。