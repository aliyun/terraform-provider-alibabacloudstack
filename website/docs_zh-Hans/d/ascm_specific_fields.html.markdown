---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_specific_fields"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-specific-fields"
description: |-
    查询特定字段
---

# alibabacloudstack_ascm_specific_fields

根据指定过滤条件列出当前凭证权限可以访问的特定字段列表。

## 示例用法

```
data "alibabacloudstack_ascm_specific_fields" "specifields" {
  group_filed ="storageType"
  resource_type ="OSS"
  output_file = "fields"
}
output "specifields" {
  value = data.alibabacloudstack_ascm_specific_fields.specifields.*
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 特定字段ID的列表。
* `group_filed` - (必填) 要查询有效值的字段。
* `resource_type` - (必填) 通过指定的资源类型过滤结果。有效值：OSS、ADB、DRDS、SLB、NAT、MAXCOMPUTE、POSTGRESQL、ECS、RDS、IPSIX、REDIS、MONGODB 和 HITSDB。
* `label` - (可选) 指定是否对字段进行国际化。有效值：true 和 false。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `specific_fields` - 特定字段的列表。