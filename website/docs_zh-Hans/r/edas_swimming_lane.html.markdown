---
subcategory: "Edas"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane"
sidebar_current: "docs-alibabacloudstack-resource-edas-swimming-lane"
description: |-
  提供阿里云EDAS游泳道资源
---

# alibabacloudstack\_edas\_swimming\_lane

提供EDAS游泳道资源。

## 示例用法

```terraform
resource "alibabacloudstack_edas_swimming_lane" "example" {
  name             = "example_value"
  group_id         = "12345"
  apps             = ["example_app_id"]
  priority         = 1
  path             = "/example"
  condition        = "OR"
  enabled          = true

  rest_items {
    type     = "header"
    name     = "example_header"
    value    = "example_value"
    cond     = "=="
    operator = "rawvalue"
  }
}
```

## 参数说明

以下参数是可支持的：

* `name` - (必选) 游泳道的名称。
* `group_id` - (必选) 游泳道组的ID。
* `apps` - (必选) 应用程序ID列表。至少需要指定一个应用程序ID。
* `priority` - (必选) 游泳道的优先级。
* `path` - (必选) 游泳道的路径。
* `condition` - (必选) 游泳道的条件。有效值：`OR`，`ADD`。
* `rest_items` - (必选) 游泳道的REST项。每项支持以下参数：
  * `type` - (必选) REST项的类型。有效值：`cookie`，`header`，`param`。
  * `name` - (必选) REST项的名称。
  * `value` - (必选) REST项的值。
  * `cond` - (必选) REST项的条件。有效值：`==`，`!=`，`>`，`<`，`>=`，`<=`，`list`。
  * `operator` - (可选) REST项的操作符。有效值：`rawvalue`，`mod`，`list`。默认值：`rawvalue`。
* `logical_region_id` - (可选) EDAS实例的逻辑区域ID。如果未指定，则使用提供商的默认区域ID。
* `enabled` - (可选) 游泳道是否启用。默认值：`false`。

## 属性导出

以下属性会被导出：

* `id` - 游泳道的ID。
* `lane_id` - 游泳道的道ID。

## 导入

EDAS游泳道可以使用id进行导入，例如

```bash
$ terraform import alibabacloudstack_edas_swimming_lane.example logical_region_id:group_id:lane_id
```