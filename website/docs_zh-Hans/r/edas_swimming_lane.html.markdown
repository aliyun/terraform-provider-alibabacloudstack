---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane"
sidebar_current: "docs-Alibabacloudstack-resource-edas-swimming-lane"
description: |-
  提供阿里云EDAS泳道（Swimming Lane）资源
---

# alibabacloudstack\_edas\_swimming\_lane

提供 EDAS 泳道（Swimming Lane）资源。

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

* `name` - (必选) 泳道名称。
* `group_id` - (必选) 泳道组 ID。
* `apps` - (必选) 应用 ID 列表。至少需指定一个应用 ID。
* `priority` - (必选) 泳道优先级，取值范围 1-100。
* `path` - (必选) 泳道匹配的请求路径。
* `condition` - (必选) 匹配条件之间的逻辑关系。有效值：`OR`（任一条件满足）、`ADD`（所有条件满足）。
* `rest_items` - (必选) 流量规则项列表。每个规则项包含以下参数：
  * `type` - (必选) 规则类型。有效值：`cookie`、`header`、`param`。
  * `name` - (必选) 规则键名。
  * `value` - (必选) 规则值。
  * `cond` - (必选) 规则匹配条件。有效值：`==`、`!=`、`>`、`<`、`>=`、`<=`、`list`。
  * `operator` - (可选) 值类型。有效值：`rawvalue`（原始值）、`mod`（取模）、`list`（列表值）。默认值：`rawvalue`。
* `logical_region_id` - (可选) EDAS 实例的逻辑区域 ID。格式为 `物理区域ID:自定义命名空间标识`，例如 `cn-hangzhou:test`。如果未指定，则使用 Provider 的默认区域 ID。
* `enabled` - (可选) 是否启用泳道。默认值：`false`。

## 属性导出

以下属性会被导出：

* `id` - 泳道 ID。格式为 `logical_region_id:group_id:lane_id`。
* `lane_id` - 泳道在 EDAS 中的唯一标识。

## 导入

EDAS 泳道可以使用 `logical_region_id:group_id:lane_id` 格式的 ID 导入，例如：

```bash
$ terraform import alibabacloudstack_edas_swimming_lane.example cn-hangzhou:test:12345:88
```