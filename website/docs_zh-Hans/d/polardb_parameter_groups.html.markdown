---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_parameter_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-parameter-groups"
description: |-
  提供阿里云账号下拥有的 PolarDB 参数模板列表。
---

# alibabacloudstack\_polardb\_parameter\_groups

此数据源提供根据指定过滤条件列出的阿里云账号下的 PolarDB 参数模板资源列表。

## 示例用法

```hcl
data "alibabacloudstack_polardb_parameter_groups" "default" {
  engine         = "MySQL"
  engine_version = "8.0"
}

# 按名称正则过滤
data "alibabacloudstack_polardb_parameter_groups" "example" {
  name_regex = "^my-parameter-group"
}

# 按 ID 过滤
data "alibabacloudstack_polardb_parameter_groups" "ids" {
  ids = ["pg-xxxxxxxxxxxxx"]
}
```

## 参数参考

以下参数是支持的：

  * `ids` - (选填) 用于过滤结果的参数模板 ID 列表。
  * `name_regex` - (选填) 用于过滤结果的参数模板名称，支持正则表达式。
  * `engine` - (选填) 数据库引擎类型。有效值：`MySQL`、`PostgreSQL`、`Oracle`。
  * `engine_version` - (选填) 数据库引擎版本。
  * `parameter_group_type` - (选填) 参数模板类型。

## Attributes Reference

除了上述参数外，还导出以下属性：

  * `ids` - 参数模板 ID 列表。
  * `names` - 参数模板名称列表。
  * `groups` - 参数模板列表。每个元素包含以下属性：
    * `id` - 参数模板 ID。
    * `parameter_group_id` - 参数模板 ID。
    * `parameter_group_name` - 参数模板名称。
    * `parameter_group_desc` - 参数模板描述。
    * `engine` - 数据库引擎类型。
    * `engine_version` - 数据库引擎版本。
    * `parameter_group_type` - 参数模板类型。
    * `force_restart` - 表示参数生效是否需要重启。
    * `param_counts` - 参数模板中的参数数量。
    * `created` - 参数模板创建时间。
    * `modified` - 参数模板最后修改时间。
    * `parameters` - 参数模板中的参数列表。每个元素包含：
      * `param_name` - 参数名称。
      * `param_value` - 参数值。
