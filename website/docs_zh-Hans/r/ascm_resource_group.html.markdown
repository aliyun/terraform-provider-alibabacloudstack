---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_group"
sidebar_current: "docs-alibabacloudstack-resource-ascm-resource-group"
description: |-
  编排Ascm资源组
---

# alibabacloudstack_ascm_resource_group

使用Provider配置的凭证在指定的资源集下编排Ascm资源组。

-> **注意：** 如果您需要在一个模板中创建不同资源集中的不同资源，则需要参考 [多资源组](ascm_resource_group_mult.html.markdown) 中的方法。



## 示例用法

```
resource "alibabacloudstack_ascm_resource_group" "default" {
    name = "Resource_Group_Name"
}

data "alibabacloudstack_ascm_resource_groups" "default" {
    name_regex = alibabacloudstack_ascm_resource_group.default.name
}
output "rg" {
  value = data.alibabacloudstack_ascm_resource_groups.default.*
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 资源组的名称。该名称可以包含 2 到 128 个字符。

* `organization_id` - (已弃用, 可选) 该参数已弃用。资源组将创建在当前用户所属的组织下。修改此参数将强制创建新资源。

## 属性说明

导出以下属性：

* `id` - 资源组的 ID。格式为 `<organization_id>:<resource_group_id>`。
* `name` - 资源组的名称。
* `rg_id` - 资源组的内部 ID（API 返回的 ID）。
* `organization_id` - 资源组所属的组织 ID。

## Import

ASCM 资源组可以通过组织 ID 和资源组 ID（用冒号分隔）导入，例如：

```
$ terraform import alibabacloudstack_ascm_resource_group.example 12345:67890
```