---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lanes"
sidebar_current: "docs-alibabacloudstack-datasource-edas-swimming-lanes"
description: |-
  提供用户的 Edas 游泳道列表。
---

# alibabacloudstack\_edas\_swimming\_lanes

该数据源提供当前阿里云用户可用的 Edas 游泳道。

## 示例用法

```terraform
data "alibabacloudstack_edas_swimming_lanes" "example" {
  group_id = "12345"
  logical_region_id = "cn-beijing:test"
}

output "first_swimming_lane_id" {
  value = data.alibabacloudstack_edas_swimming_lanes.example.lanes.0.id
}
```

## 参数说明

以下参数被支持：

* `logical_region_id` - （可选）EDAS 实例的逻辑区域 ID。如果未指定，则使用提供商的默认区域 ID。
* `group_id` - （必选）游泳道组的 ID。
* `ids` - （可选）游泳道 ID 列表。
* `name_regex` - （可选）用于按游泳道名称过滤结果的正则表达式。

## 属性导出

以下属性会被导出：

* `ids` - 游泳道 ID 列表。
* `lanes` - 游泳道列表。每个元素包含以下属性：
  * `id` - 游泳道的 ID。
  * `name` - 游泳道的名称。
  * `group_id` - 游泳道的组 ID。
  * `logical_region_id` - 游泳道的逻辑区域 ID。
  * `apps` - 与此游泳道关联的应用程序 ID 列表。
  * `priority` - 游泳道的优先级。
  * `path` - 游泳道的路径。
  * `condition` - 游泳道的条件。
  * `rest_items` - 游泳道的 REST 项。
    * `type` - REST 项的类型。
    * `name` - REST 项的名称。
    * `value` - REST 项的值。
    * `cond` - REST 项的条件。
    * `operator` - REST 项的操作符。
  * `enabled` - 游泳道是否启用。
  * `lane_id` - 游泳道的道 ID。