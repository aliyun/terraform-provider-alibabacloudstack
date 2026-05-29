---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane_group"
sidebar_current: "docs-Alibabacloudstack-resource-edas-swimming-lane-group"
description: |-
  提供阿里云EDAS游泳道组资源
---

# alibabacloudstack\_edas\_swimming\_lane\_group

提供EDAS游泳道组资源。


## 示例用法

```terraform
resource "alibabacloudstack_edas_swimming_lane_group" "example" {
  name             = "example_value"
  entry_app_id     = "example_app_id"
  apps             = ["example_app_id"]
  logical_region_id = "cn-beijing:test"
  strategy_type    = "CONTENT"
}
```

## 参数说明

* `name` - (必选) 游泳道组的名称。
* `entry_app_id` - (必选) 入口应用程序的ID。
* `apps` - (必选) 应用程序ID列表。
* `logical_region_id` - (可选) EDAS实例的逻辑区域ID。如果未指定，则使用提供商的默认区域ID。
* `strategy_type` - (可选) 游泳道组的策略类型。有效值：`CONTENT`，`PERCENT`。默认值：`CONTENT`。

## 属性导出

以下属性会被导出：

* `id` - 游泳道组的ID。
* `group_id` - 游泳道组的组ID。

## 导入

EDAS游泳道组可以使用id进行导入，例如

```bash
$ terraform import alibabacloudstack_edas_swimming_lane_group.example logical_region_id:group_id
```