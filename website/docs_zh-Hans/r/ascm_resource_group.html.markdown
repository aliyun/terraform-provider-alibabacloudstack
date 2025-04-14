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
resource "alibabacloudstack_ascm_organization" "default" {
    name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_resource_group" "default" {
    organization_id = alibabacloudstack_ascm_organization.default.org_id
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

* `name` - (必填) 资源组的名称。该名称可以包含2到128个字符的字符串，必须仅包含字母数字字符或连字符（例如“-”、“.”、“_”），并且不能以连字符开头或结尾，也不能以`http://`或`https://`开头。默认值为`null`。
* `organization_id` - (必填) 组织的ID。
* `rg_id` - (可选) 资源组的ID。如果已知资源组ID，可以通过此参数直接引用现有资源组。

## 属性说明

导出以下属性：

* `id` - 资源组的名称和ID。格式为`Name:ID`。
* `name` - 资源组的名称。
* `rg_id` - 资源组的ID。