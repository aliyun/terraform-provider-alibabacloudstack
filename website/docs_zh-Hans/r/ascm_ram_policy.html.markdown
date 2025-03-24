---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_policy"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-policy"
description: |-
  编排ASCM的RAM策略
---

# alibabacloudstack_ascm_ram_policy

使用Provider配置的凭证在指定的资源集下编排ASCM的RAM策略。

## 示例用法

```
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = "TestPolicy"
  description = "Testing"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"
}
output "rampolicy" {
  value = alibabacloudstack_ascm_ram_policy.default.*
}
```

## 参数参考

以下参数是支持的：

* `name` - (必填)  RAM策略名称。
* `policy_document` - (必填)  策略文档。
* `description` - (可选) RAM策略的描述。
* `name` - (必填)   RAM策略名称，长度在3到64个字符之间。
* `policy_document` - (必填)   策略的文档。

## 属性参考

导出以下属性：

* `id` - RAM策略名称的用户。
* `ram_id` - RAM策略的ID。 导出属性，表示RAM策略的唯一标识符。