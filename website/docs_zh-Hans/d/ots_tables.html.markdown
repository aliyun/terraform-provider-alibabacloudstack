---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_tables"
sidebar_current: "docs-alibabacloudstack-datasource-ots-tables"
description: |-
    查询表格存储（OTS）数据表
---

# alibabacloudstack_ots_tables

根据指定过滤条件列出当前凭证权限可以访问的表格存储（OTS）数据表列表。


## 示例用法

``` terraform
data "alibabacloudstack_ots_tables" "tables_ds" {
  instance_name = "sample-instance"
  name_regex    = "sample-table"
  output_file   = "tables.txt"
}

output "first_table_id" {
  value = "${data.alibabacloudstack_ots_tables.tables_ds.tables.0.id}"
}
```

## 参数参考

以下是支持的参数：

* `instance_name` - OTS 实例名称。
* `ids` - (可选) 表 ID 列表。
* `name_regex` - (可选) 用于通过表名筛选结果的正则表达式字符串。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 表 ID 列表。
* `names` - 表名称列表。
* `tables` - 表列表。每个元素包含以下属性：
  * `id` - 表的 ID。值为 `<instance_name>:<table_name>`。
  * `instance_name` - OTS 实例名称。
  * `table_name` - OTS 表的名称，该名称一旦创建后无法更改。
  * `primary_key` - 表示 `TableMeta` 的属性，指示表的结构信息。
    * `name` - 属性名称。
    * `type` - 属性类型。可用值为 {1 表示 int, 2 表示 string, 3 表示 binary}。
  * `time_to_live` - 此表中存储的数据保留时间。
  * `max_version` - 此表中存储的最大版本数。