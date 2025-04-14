---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_zones"
sidebar_current: "docs-alibabacloudstack-datasource-elasticsearch-zones"
description: |-
  查询Elasticsearch 可用区
---

# alibabacloudstack_elasticsearch_zones

根据指定过滤条件列出当前凭证权限可以访问的Elasticsearch可用区列表



## 示例用法

```
# 声明数据源
data "alibabacloudstack_elasticsearch_zones" "zones_ids" {}
```

## 参数说明

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多 AZ 配置中使用。默认值为 `false`。如果设置为 `true`，则返回的可用区列表将仅包含支持多 AZ 配置的可用区。多 AZ 通常用于启动高可用的 Elasticsearch 实例。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 区域 ID 列表。此列表包含了所有符合条件的可用区的 ID。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 区域的 ID。表示该可用区的唯一标识符。
  * `multi_zone_ids` - 多区域中的区域 ID 列表。当 `multi` 参数设置为 `true` 时，此字段会列出支持多 AZ 配置的可用区 ID。