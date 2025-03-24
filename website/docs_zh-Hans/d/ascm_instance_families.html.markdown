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

## 参数参考

支持以下参数：

* `ids` - (可选) 实例族ID列表。
* `name_regex` - (可选) 用于通过轨迹名称过滤结果的正则表达式字符串。
* `status` - (可选) 指定状态以按实例族的可用性过滤结果。
* `resource_type` - (可选) 通过指定的资源类型过滤结果。
* `families` - (可选) 实例族列表。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `families` - 实例族列表。每个元素包含以下属性：
    * `id` - 实例族的ID。
    * `order_by_id` - 实例族的排序ID。
    * `series_name` - 实例族的系列名称。
    * `modifier` - 修改者名称。
    * `series_name_label` - 实例族系列名称的标签。
    * `is_deleted` - 以“Y”或“N”形式指定状态。
    * `resource_type` - 指定的资源类型。
