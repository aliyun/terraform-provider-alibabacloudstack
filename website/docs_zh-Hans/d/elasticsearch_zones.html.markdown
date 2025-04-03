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

## 参数参考

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以在多 AZ 配置中使用。默认为 `false`。多 AZ 通常用于启动 Elasticsearch 实例。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 区域 ID 列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 区域的 ID。
  * `multi_zone_ids` - 多区域中的区域 ID 列表。