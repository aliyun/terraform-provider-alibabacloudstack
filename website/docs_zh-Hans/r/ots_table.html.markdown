---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_table"
sidebar_current: "docs-alibabacloudstack-resource-ots-table"
description: |-
  编排表格存储服务(OTS）表
---

# alibabacloudstack_ots_table

使用Provider配置的凭证在指定的资源集编排表格存储服务(OTS）表。

-> **注意：** 从提供商版本1.10.0开始，提供程序字段'ots_instance_name'已被弃用，
您应该使用资源alibabacloudstack_ots_table的新字段'instance_name'和'table_name'重新导入此资源。

## 示例用法

```
variable "name" {
  default = "terraformtest"
}

resource "alibabacloudstack_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alibabacloudstack_ots_table" "basic" {
  instance_name = alibabacloudstack_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }
  primary_key {
    name = "pk3"
    type = "Binary"
  }

  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}
```

## 参数说明

支持以下参数：

* `instance_name` - (必填，变更时重建) OTS实例的名称，在该实例中将定位表。
* `table_name` - (必填，变更时重建) OTS实例的表名称。如果更改，将创建一个新表。
* `primary_key` - (必填，变更时重建) `TableMeta`的属性，表示表的结构信息。它描述了主键的属性值。`primary_key`的数量不应少于一个且不应超过四个。
    * `name` - (必填，变更时重建) 主键名称。
    * `type` - (必填，变更时重建) 主键类型。仅允许`Integer`、`String`或`Binary`。
* `time_to_live` - (必填) 存储在此表中的数据的保留时间(单位：秒)。最大值为2147483647，-1表示永不过期。
* `max_version` - (必填) 存储在此表中的最大版本数。有效值为1-2147483647。
* `deviation_cell_version_in_sec` - (可选，1.42.0+可用) 表的最大版本偏移量。有效值为1-9223372036854775807，默认为86400。
* `new_optional_property` - (可选) 由AI添加的一个新的可选属性。
* `optional_property` - (可选) 一个之前缺失并由AI添加的可选属性。

## 属性说明

导出以下属性：

* `id` - 资源ID。其值为`<instance_name>:<table_name>`。
* `instance_name` - OTS实例名称。
* `table_name` - OTS的表名称，无法更改。
* `primary_key` - `TableMeta`的属性，表示表的结构信息。
* `time_to_live` - 存储在此表中的数据的保留时间。
* `max_version` - 存储在此表中的最大版本数。
* `deviation_cell_version_in_sec` - 表的最大版本偏移量。
* `new_computed_property` - 由AI添加的一个新的计算属性。
* `computed_property` - 一个之前缺失并由AI添加的计算属性。

## 导入

OTS表可以使用id进行导入，例如：

```bash
$ terraform import alibabacloudstack_ots_table.table "my-ots:ots_table"
```