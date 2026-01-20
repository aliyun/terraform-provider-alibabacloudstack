---
subcategory: "ApsaraDB for Redis (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_parameter_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-kvstore-parameter-groups"
description: |-
  查询阿里云Redis（KVStore）参数模板列表。
---

# alibabacloudstack_kvstore_parameter_groups

查询阿里云Redis（KVStore）参数模板列表。该数据源用于检索已创建的参数模板信息，包括参数组ID、名称、描述、引擎版本等。

## 示例用法

```hcl

variable "name" {
  default = "tf-kvparamgroup34466"
}

resource "alibabacloudstack_kvstore_parameter_group" "default" {
  character_type       = "logic"
  parameter_group_name = var.name
  engine_version       = "7.0"
  parameter_group_desc = var.name
  parameters {
    param_name = "resp_version"
    value      = "3"
  }
  parameters {
    param_name = "rt_threshold_ms"
    value      = "400"
  }
  parameters {
    param_name = "#no_loose_check-whitelist-always"
    value      = "yes"
  }

}



data "alibabacloudstack_kvstore_parameter_groups" "default" {
  name_regex = alibabacloudstack_kvstore_parameter_group.default.parameter_group_name
}
```

## 参数说明

以下参数支持过滤查询结果：

* `character_type` (字符串)：参数模板的字符类型。可选值：`logic`（逻辑参数模板）、`normal`（物理参数模板）。

* `engine_version` (字符串)：参数模板的引擎版本。例如：`7.0`。

* `ids` (列表)：参数模板ID列表，用于过滤结果。列表中的每个元素都是字符串类型。

* `name_regex` (字符串)：参数模板名称的正则表达式，用于过滤结果。

## 属性说明

以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，由过滤后的参数模板ID哈希生成。

* `character_type` (字符串)：参数模板的字符类型。

* `create_time` (字符串)：参数模板的创建时间，格式为ISO 8601标准时间格式。

* `engine_version` (字符串)：参数模板的引擎版本。

* `groups` (列表)：匹配的参数模板列表。每个元素包含以下属性：
  * `character_type` (字符串)：参数模板的字符类型。
  * `create_time` (字符串)：参数模板的创建时间。
  * `engine_version` (字符串)：参数模板的引擎版本。
  * `id` (字符串)：参数模板ID，与parameter_group_id相同。
  * `is_dynamic` (整数)：参数是否为动态参数。1表示是，0表示否。
  * `parameter_group_desc` (字符串)：参数模板的描述信息。
  * `parameter_group_id` (字符串)：参数模板ID。
  * `parameter_group_name` (字符串)：参数模板名称。
  * `parameters` (列表)：参数模板中的参数列表。每个元素包含：
    * `param_name` (字符串)：参数名称。
    * `value` (字符串)：参数值。
  * `type` (整数)：参数模板类型。

* `names` (列表)：匹配的参数模板名称列表。

* `parameter_group_desc` (字符串)：参数模板的描述信息。

* `parameter_group_id` (字符串)：参数模板ID。

* `parameter_group_name` (字符串)：参数模板名称。

* `type` (整数)：参数模板类型。