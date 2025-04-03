---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organizations"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-organizations"
description: |-
    查询ascm组织
---

# alibabacloudstack_ascm_organizations

根据指定过滤条件列出当前凭证权限可以访问的组织列表。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Test_org"
}
output "orgres" {
  value = alibabacloudstack_ascm_organization.default.*
}
data "alibabacloudstack_ascm_organizations" "default" {
    name_regex = alibabacloudstack_ascm_organization.default.name
    parent_id = alibabacloudstack_ascm_organization.default.parent_id
}
output "orgs" {
  value = data.alibabacloudstack_ascm_organizations.default.*
}
```

## 参数参考

以下是支持的参数：

* `ids` - (可选) 组织ID列表。
* `name_regex` - (可选) 用于按组织名称过滤结果的正则表达式字符串。
* `parent_id` - (可选) 通过指定的组织父级ID过滤结果。
* `organizations` - (可选) 组织列表。


## 属性参考

除了上述列出的参数外，还导出以下属性：

* `organizations` - 组织列表。每个元素包含以下属性：
  * `id` - 组织的ID。
  * `name` - 组织名称。
  * `cuser_id` - Cuser的ID。
  * `muser_id` - Muser的ID。
  * `alias` - 组织的别名。
  * `parent_id` - 组织的父级ID。
  * `internal` - 组织类型，是否为内部组织。
  * `name_regex` - 按组织名称过滤结果的正则表达式字符串。