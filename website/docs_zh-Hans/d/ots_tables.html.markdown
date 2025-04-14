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

## 参数说明

以下是支持的参数：

* `instance_name` - （必需）OTS 实例名称。
* `ids` - （可选）表 ID 列表，用于筛选特定的表。
* `name_regex` - （可选）用于通过表名筛选结果的正则表达式字符串。
* `output_file` - （可选）将查询结果保存到文件的路径。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 表 ID 列表，每个 ID 唯一标识一个表。
* `names` - 表名称列表，包含所有匹配条件的表名。
* `tables` - 表详细信息列表。每个元素包含以下属性：
  * `id` - 表的唯一标识符，格式为 `<instance_name>:<table_name>`。
  * `instance_name` - OTS 实例名称，该表所属的实例。
  * `table_name` - OTS 表的名称，该名称一旦创建后无法更改。
  * `primary_key` - 表的主键结构信息，包含以下子属性：
    * `name` - 主键列的名称。
    * `type` - 主键列的数据类型，可用值为 {1 表示 int, 2 表示 string, 3 表示 binary}。
  * `time_to_live` - 数据在表中的保留时间（以秒为单位），设置为 -1 表示数据永不过期。
  * `max_version` - 表中存储的最大版本数，设置为 1 表示仅保留最新版本。