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

## 参数参考

支持以下参数：

* `ids` - (可选) 资源组 ID 列表。
* `name_regex` - (可选) 用于按资源组名称过滤结果的正则表达式字符串。
* `organization_id` - (可选) Alibabacloudstack Cloud 账户的组织 ID。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 资源组 ID 列表。
* `names` - 资源组名称列表。
* `groups` - 资源组列表。每个元素包含以下属性：
  * `id` - 资源组的 ID。
  * `name` - 资源组的名称。
  * `organization_id` - Alibabacloudstack Cloud 账户的组织 ID。
  * `gmt_created` - 资源组的创建时间。
  * `rs_id` - 资源组集的 ID(如 rs-xxxxx)。
  * `creator` - 资源组创建者的 ID。
  * `resource_group_type` - 资源集类型。可用值及其含义：
    * 1: 默认资源集。
    * 0: 非默认资源集。
* `computed_organization_id` - 与资源组关联的计算得出的组织 ID。