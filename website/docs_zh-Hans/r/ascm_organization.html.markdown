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

## 参数说明

以下是支持的参数：

* `name` - (必填) 组织的名称。该名称可以包含2到128个字符，必须仅包含字母数字字符或连字符（如“-”、“.”、“_”），并且不能以连字符开头或结尾，也不能以 `http://` 或 `https://` 开头。默认值为 `null`。
* `parent_id` - (可选) 父组织的ID。默认情况下，`parent_id` 的值为 `"1"`。对于普通用户（非管理员），`parent_id` 将是其所属组织的ID。
* `person_num` - (可选) 保留参数，目前暂无实际用途。
* `resource_group_num` - (可选) 保留参数，目前暂无实际用途。

## 属性说明

以下属性被导出：

* `id` - 组织的名称和ID，格式为 `Name:ID`。
* `org_id` - 组织的唯一标识符（ID）。