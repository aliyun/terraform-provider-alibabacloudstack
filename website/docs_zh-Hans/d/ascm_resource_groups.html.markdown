---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_groups"
sidebar_current: "docs-alibabacloudstack-datasource-alibabacloudstack-ascm-resource-groups"
description: |-
   查询ascm资源集

---

# alibabacloudstack_ascm_resource_groups

根据指定过滤条件列出当前凭证权限可以访问的资源组列表。


## 示例用法

```
data "alibabacloudstack_ascm_resource_groups" "default" {
  name_regex = "another resource" # 可选
}

output "resource_group" {
  value = data.alibabacloudstack_ascm_resource_groups.default.groups
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 资源组 ID 列表，用于筛选特定的资源组。
* `name_regex` - (可选) 一个正则表达式字符串，用于按资源组名称过滤结果。
* `organization_id` - (可选) Alibabacloudstack Cloud 账户的组织 ID，用于限定查询范围。如果未指定，则默认使用当前账户的组织 ID。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 符合条件的资源组 ID 列表。
* `names` - 符合条件的资源组名称列表。
* `groups` - 符合条件的资源组详细信息列表。每个元素包含以下属性：
  * `id` - 资源组的唯一标识符。
  * `name` - 资源组的名称。
  * `organization_id` - 资源组所属的 Alibabacloudstack Cloud 账户的组织 ID。
  * `gmt_created` - 资源组的创建时间，格式为标准时间戳。
  * `rs_id` - 资源组所属的资源组集合的唯一标识符（如 rs-xxxxx）。
  * `creator` - 创建该资源组的用户的唯一标识符。
  * `resource_group_type` - 资源组的类型，取值及其含义如下：
    * 1: 默认资源组。
    * 0: 非默认资源组。
* `computed_organization_id` - 根据查询结果计算得出的与资源组关联的组织 ID。此值可能与输入的 `organization_id` 不同，尤其是在指定了过滤条件的情况下。
