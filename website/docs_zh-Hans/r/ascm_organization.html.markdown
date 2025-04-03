---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organization"
sidebar_current: "docs-alibabacloudstack-resource-ascm-organization"
description: |-
  编排Ascm组织
---

# alibabacloudstack_ascm_organization

使用Provider配置的凭证编排Ascm组织。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "apsara_Organization"
  parent_id = "19"
}
output "org" {
  value = alibabacloudstack_ascm_organization.default.*
}
```

## 参数参考

以下是支持的参数：

* `name` - (必填) 组织的名称。该名称可以包含2到128个字符，必须仅包含字母数字字符或连字符，例如“-”、“.”、“_”，并且不能以连字符开头或结尾，也不能以http://或https://开头。默认值为null。
* `parent_id` - (可选) 其父组织的ID。父级ID的默认值为“1”。对于普通用户(非管理员)，parent_id将是其组织ID。
* `person_num` - (可选) 保留参数。
* `resource_group_num` - (可选) 保留参数。

## 属性参考

以下属性被导出：

* `id` - 组织的名称和ID。格式为`Name:ID`
* `org_id` - 组织的ID。