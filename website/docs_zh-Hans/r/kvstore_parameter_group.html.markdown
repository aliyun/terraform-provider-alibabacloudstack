---
subcategory: "云数据库 Redis 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_parameter_group"
sidebar_current: "docs-Alibabacloudstack-resource-kvstore-parameter-group"
description: |-
  管理R-kvstore参数模板
---

# alibabacloudstack_kvstore_parameter_group

管理R-kvstore参数模板，用于创建和管理Redis等数据库引擎的参数配置集合。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-kvparmgroup31846"
}


resource "alibabacloudstack_kvstore_parameter_group" "default" {
  character_type       = "logic"
  engine               = "Redis"
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
```

## 参数说明

支持以下参数：

* `character_type` - (必填, 变更时重建) 参数模板的字符类型，例如`logic`表示逻辑参数模板，`normal`（物理参数模板）。
* `engine_version` - (必填, 变更时重建) 数据库引擎版本，例如"7.0"。
* `parameter_group_desc` - (必填, 变更时重建) 参数模板的描述信息，长度1-256字符。
* `parameter_group_name` - (必填, 变更时重建) 参数模板的名称，长度1-128字符。
* `engine` - (可选, 变更时重建) 数据库引擎类型，默认值为"Redis"。
* `parameters` - (可选, 变更时重建) 自定义参数列表，结构为集合。每个元素包含：
  * `param_name` - (必填, 变更时重建) 参数名称
  * `value` - (必填, 变更时重建) 参数值

## 属性说明

导出以下属性：

* `id` - 参数模板的ID。
* `create_time` - 参数模板的创建时间，格式为ISO 8601标准时间（如"2026-01-20T07:27Z"）。
* `is_dynamic` - 是否动态参数模板（1表示是，0表示否）。
* `type` - 参数模板的类型标识。