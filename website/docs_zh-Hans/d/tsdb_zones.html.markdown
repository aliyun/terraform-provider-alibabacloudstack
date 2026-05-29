---
subcategory: "时间序列数据库"
layout: "alibabacloudstack"
page_title: "阿里云：alibabacloudstack_tsdb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-tsdb-zones"
description: |-
  查询时间序列数据库(TSDB)实例区域
---

# alibabacloudstack_tsdb_zones

根据指定过滤条件列出当前凭证权限可以访问的时间序列数据库(TSDB)实例区域列表。


## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_tsdb_zones" "example" {}

output "first_tsdb_zones_id" {
  value = data.alibabacloudstack_tsdb_zones.example.zones.0.zone_id
}
```

## 参数说明

此数据源无需参数。

## 属性说明

导出以下属性：

* `ids` - TSDB 实例区域 ID 列表。
* `zones` - TSDB 实例区域列表。每个元素包含以下属性：
  * `id` - 区域的唯一标识符。
  * `zone_id` - 区域 ID。

> **注意：** `local_name` 属性在专有云 v3.16 及更高版本中不支持。