---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_parameter_group"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-parameter-group"
description: |-
  编排 PolarDB 集群版参数模板资源。
---

# alibabacloudstack_polardb_cluster_parameter_group

使用 Provider 配置的凭证编排 PolarDB 集群版（Cluster）参数模板资源，参数模板可应用于集群版 PolarDB 集群。

本资源调用 PolarDB `2017-08-01` 版本的 `CreateParameterGroup`、`DescribeParameterGroups` 和 `DeleteParameterGroup` 接口。`CreateParameterGroup` 接口直接返回参数模板 ID，并将其作为 Terraform 资源 ID。

-> **注意:** 某些资源在创建参数时可能存在限制。更多详情请参见 [PolarDB 产品文档](https://help.aliyun.com/product/polardb.html)。

-> **注意:** `parameters` 仅在创建参数模板时提交给接口。创建后修改该映射不会更新已有参数模板，Provider 也不会回读参数值，因此创建后 `parameters` 出现差异时通常需要重建资源。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccclusterparamgroup"
}

resource "alibabacloudstack_polardb_cluster_parameter_group" "default" {
  engine               = "MySQL"
  engine_version       = "8.0"
  parameter_group_name = var.name
  parameter_group_desc = "Terraform 集群参数模板"
  parameters = {
    "wait_timeout" = "3600"
  }
}
```

## 参数说明

支持以下参数：

  * `engine` - (必填, 变更时重建) 集群参数模板的数据库引擎类型。PolarDB 集群版参数模板仅支持 `MySQL`，引擎取值由接口校验。
  * `engine_version` - (必填, 变更时重建) 数据库引擎版本。有效值取决于引擎类型。对于 MySQL：`5.6`、`5.7`、`8.0`。
  * `parameter_group_name` - (必填) 参数模板名称。
  * `parameter_group_desc` - (选填, 可回读) 参数模板的描述信息。
  * `parameters` - (选填) 要配置在参数模板中的参数映射。键为参数名称，值为参数值，键与值均为字符串。

## 属性参考

除了上述所有参数外，还导出了以下属性：

  * `id` - 参数模板 ID。与 `parameter_group_id` 相同。
  * `parameter_group_id` - 参数模板 ID，由 `CreateParameterGroup` 接口返回。
  * `parameter_group_type` - 参数模板类型。
  * `param_counts` - 参数模板中的参数数量。
  * `force_restart` - 表示修改参数后是否需要重启集群。
  * `created` - 参数模板的创建时间。
  * `modified` - 参数模板的最后修改时间。

## Import

PolarDB 集群版参数模板可以使用参数模板 ID 导入，例如：

```
$ terraform import alibabacloudstack_polardb_cluster_parameter_group.example pg-12345678
```

-> **注意:** 读取参数模板时接口不返回 `parameters`，因此导入后无法恢复该字段。
