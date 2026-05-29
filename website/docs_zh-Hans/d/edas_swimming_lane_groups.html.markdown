---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-swimming-lane-groups"
description: |-
  提供用户的 Edas 游泳道组列表。
---

# alibabacloudstack\_edas\_swimming\_lane\_groups

该数据源提供当前阿里云用户可用的 Edas 游泳道组。

## 示例用法

```terraform
data "alibabacloudstack_edas_swimming_lane_groups" "example" {
  logical_region_id = "cn-beijing:test"
}

output "first_swimming_lane_group_id" {
  value = data.alibabacloudstack_edas_swimming_lane_groups.example.groups.0.id
}
```

## 参数说明

以下参数被支持：

* `logical_region_id` - （可选）EDAS 实例的逻辑区域 ID。如果未指定，则使用提供商的默认区域 ID。
* `ids` - （可选）游泳道组 ID 列表。
* `name_regex` - （可选）用于按游泳道组名称过滤结果的正则表达式。

## 属性导出

以下属性会被导出：

* `ids` - 游泳道组 ID 列表。
* `groups` - 游泳道组列表。每个元素包含以下属性：
  * `id` - 游泳道组的 ID。
  * `name` - 游泳道组的名称。
  * `entry_app_id` - 入口应用程序的 ID。
  * `apps` - 与此游泳道组关联的应用程序 ID 列表。
  * `logical_region_id` - 游泳道组的逻辑区域 ID。
  * `strategy_type` - 游泳道组的策略类型。
  * `group_id` - 游泳道组的组 ID。