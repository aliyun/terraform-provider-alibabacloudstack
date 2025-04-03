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

## 参数参考

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多AZ配置中使用。默认为`false`。多AZ通常用于启动KVStore实例。
* `instance_charge_type` - (可选) 通过特定的实例计费类型过滤结果。有效值：`PrePaid`和`PostPaid`。默认为`PostPaid`。
* `ids` - (可选) 可用区ID列表。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区的ID。
  * `multi_zone_ids` - 多可用区中的可用区ID列表。
  * `zones` - 可用区列表。