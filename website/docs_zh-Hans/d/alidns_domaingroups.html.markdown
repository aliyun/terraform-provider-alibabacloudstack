---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_groups"
sidebar_current: "docs-alibabacloudstack-datasource-dns-groups"
description: |-
    查询DNS分组

---

# alibabacloudstack_dns_groups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_alidns_domaingroups`

根据指定过滤条件列出当前凭证权限可以访问的DNS分组列表。

## 示例用法

```
data "alibabacloudstack_dns_groups" "groups_ds" {
  name_regex  = "^y[A-Za-z]+"
  output_file = "groups.txt"
}

output "first_group_name" {
  value = "${data.alibabacloudstack_dns_groups.groups_ds.groups.0.group_name}"
}
```

## 参数参考

以下参数被支持：

* `name_regex` - (可选) 用于通过分组名称过滤结果的正则表达式字符串。
* `ids` - (可选) 分组ID列表。
* `names` - (可选) 分组名称列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 分组ID列表。
* `names` - 分组名称列表。
* `groups` - 分组列表。每个元素包含以下属性：
  * `group_id` - 分组的ID。
  * `group_name` - 分组的名称。
