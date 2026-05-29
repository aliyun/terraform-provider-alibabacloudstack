---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_parameter_group"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-parameter-group"
description: |-
  编排 PolarDB 参数模板资源。
---

# alibabacloudstack_polardb_parameter_group

使用 Provider 配置的凭证在指定的资源集编排 PolarDB 参数模板资源。

-> **注意:** 某些资源在创建参数时可能存在限制。更多详情请参见 [PolarDB 产品文档](https://help.aliyun.com/product/polardb.html)。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccparamgroup"
}

resource "alibabacloudstack_polardb_parameter_group" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  parameter_group_name = var.name
  parameter_group_desc = "Terraform 参数模板"
  parameters = {
    "loose_multi_blocks_ddl_count" = "11"
  }
}
```

## 参数说明

支持以下参数：

  * `engine` - (必填, 变更时重建) 参数模板的数据库引擎类型。取值：`MySQL`、`PostgreSQL`、`PPAS`。
  * `engine_version` - (必填, 变更时重建) 数据库引擎版本。有效值取决于引擎类型。对于 MySQL：`5.6`、`5.7`、`8.0`。
  * `parameter_group_name` - (必填) 参数模板名称。名称在单个地域内必须唯一。
  * `parameter_group_desc` - (选填, 可回读) 参数模板的描述信息。
  * `parameters` - (选填) 要配置在参数模板中的参数映射。键为参数名称，值为参数值。

## 属性参考

除了上述所有参数外，还导出了以下属性：

  * `id` - 参数模板 ID。与 `parameter_group_id` 相同。
  * `parameter_group_id` - 参数模板 ID。
  * `parameter_group_type` - 参数模板类型。
  * `param_counts` - 参数模板中的参数数量。
  * `force_restart` - 表示修改参数后是否需要重启实例。
  * `created` - 参数模板的创建时间。
  * `modified` - 参数模板的最后修改时间。

## Import

PolarDB 参数模板可以使用参数模板 ID 导入，例如：

```
$ terraform import alibabacloudstack_polardb_parameter_group.example pg-12345678
```
