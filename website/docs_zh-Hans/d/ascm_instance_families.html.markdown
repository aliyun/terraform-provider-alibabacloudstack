---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_instance_families"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-instance-families"
description: |-
    查询ascm实例族
---

# alibabacloudstack_ascm_instance_families

根据指定过滤条件列出当前凭证权限可以访问的实例族列表。

## 示例用法

```
data "alibabacloudstack_ascm_instance_families" "default" {
  output_file = "instance_families"
  resource_type = "DRDS"
  status = "Available"
}
output "instfam" {
  value = data.alibabacloudstack_ascm_instance_families.default.*
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 实例族ID列表，用于精确匹配和过滤结果。
* `name_regex` - (可选) 一个正则表达式字符串，用于通过实例族名称进行模糊匹配过滤。
* `status` - (可选) 指定状态以按实例族的可用性过滤结果。例如，`Available` 表示仅返回可用的实例族，`Unavailable` 表示返回不可用的实例族。
* `resource_type` - (可选) 指定资源类型以过滤结果。例如，`DRDS` 表示仅返回与 DRDS 相关的实例族，`ECS` 表示返回与 ECS 相关的实例族。
* `families` - (可选) 实例族列表，用户可以通过该参数预定义或验证返回的实例族信息。每个元素包含以下属性。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `families` - 实例族列表。每个元素包含以下属性：
    * `id` - 实例族的唯一标识符，用于区分不同的实例族。
    * `order_by_id` - 实例族的排序ID，用于对查询结果中的实例族进行排序。
    * `series_name` - 实例族所属的系列名称，表示该实例族属于哪个硬件或性能系列。
    * `modifier` - 最后修改该实例族信息的用户或系统名称，用于追踪变更来源。
    * `series_name_label` - 实例族系列名称的显示标签，通常用于前端界面展示，可能包含翻译或格式化后的名称。
    * `is_deleted` - 实例族是否已被删除的状态标志，值为 `"Y"` 表示已删除，`"N"` 表示未删除。
    * `resource_type` - 实例族关联的资源类型，例如 `DRDS`、`ECS` 等，用于标识实例族适用的云服务类型。
    * `computed_attribute` - 自动生成的计算属性，其具体含义可能依赖于上下文环境，例如性能指标或推荐配置。