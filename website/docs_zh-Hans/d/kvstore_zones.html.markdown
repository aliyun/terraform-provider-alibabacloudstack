---
subcategory: "云数据库 Redis 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_zones"
sidebar_current: "docs-Alibabacloudstack-datasource-kvstore-zones"
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

* `multi` - (可选) 指示这些可用区是否可以在多AZ配置中使用。默认值为`false`。多AZ通常用于启动KVStore实例。当设置为`true`时，将返回可用区ID为`MAZ`形式或包含逗号分隔可用区ID的可用区条目。
* `instance_charge_type` - (可选) 通过特定的实例计费类型过滤结果。有效值为`PrePaid`（包年包月）和`PostPaid`（按量付费）。默认值为`PostPaid`。
* `output_file` - (可选，已废弃) 'output_file'字段已被废弃，并计划在3.19.0版本中移除。如需将内容写入文件，请改用'local_file' provider。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配的可用区ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 可用区的唯一标识符。
  * `multi_zone_ids` - 组成多可用区ID的各个单独可用区ID列表。`MAZ`形式的多可用区ID和逗号分隔的可用区ID均会被识别并拆分为此列表。