---
subcategory: "时间序列数据库 (TSDB)"
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

支持以下参数：

* `ids` - (必填) 需要查询的 TSDB 实例区域 ID 列表。此参数用于筛选特定的区域。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 筛选后返回的 TSDB 实例区域 ID 列表。
* `zones` - 筛选后返回的 TSDB 实例区域列表。每个元素包含以下属性：
  * `id` - 区域的唯一标识符，与 `zone_id` 相同。
  * `zone_id` - 区域 ID，表示该区域的唯一标识。
  * `local_name` - 区域的本地名称，通常为该区域的中文或英文描述。
  * `computed_attribute` - 计算属性，表示由系统计算得出的额外信息。