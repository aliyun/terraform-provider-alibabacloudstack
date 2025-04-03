---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ecs_instance_families"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ecs-instance-families"
description: |-
    查询ECS实例族
---

# alibabacloudstack_ascm_ecs_instance_families

根据指定过滤条件列出当前凭证权限可以访问的ECS实例族列表。

## 示例用法

```
data "alibabacloudstack_ascm_ecs_instance_families" "default" {
  status = "Available"
  output_file = "ecs_instance"
}
output "ecs_instance" {
  value = data.alibabacloudstack_ascm_ecs_instance_families.default.*
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) ECS实例族ID列表。
* `status` - (必填)  指定ECS实例族的状态来过滤结果。

## 属性参考

除了上述参数外，还导出以下属性：

* `families` - ECS实例族列表。每个元素包含以下属性：
    * `instance_type_family_id` - ECS实例族的ID。
    * `generation` - ECS实例族的代数。