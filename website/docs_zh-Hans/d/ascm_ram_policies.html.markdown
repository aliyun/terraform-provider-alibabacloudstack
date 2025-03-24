---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_ram_policies"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-policies"
description: |-
    查询RAM策略

---

# alibabacloudstack_ascm_ram_policies

根据指定过滤条件列出当前凭证权限可以访问的所有RAM策略列表。

## 示例用法

```
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = "TestPolicy"
  description = "Testing"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}
data "alibabacloudstack_ascm_ram_policies" "default" {
  name_regex = alibabacloudstack_ascm_ram_policy.default.name
}
output "ram_policies" {
  value = data.alibabacloudstack_ascm_ram_policies.default.*
}


```

## 参数参考

支持以下参数：

* `ids` - (可选) RAM策略ID列表。
* `name_regex` - (可选) 用于通过RAM策略名称过滤结果的正则表达式字符串。
* `region` - (可选) 策略所属的区域名称。

## 属性参考

除了上述参数外，还导出以下属性：

* `policies` - 策略列表。每个元素包含以下属性：
    * `id` - 策略的ID。
    * `name` - 策略名称。
    * `description` - 策略的描述。
    * `ctime` - RAM策略的创建时间。
    * `cuser_id` - 策略创建者的ID。
    * `region` - 策略所属的区域名称。
    * `policy_document` - 策略文档。
    * `output_file` - 保存数据源结果的文件名(在运行`terraform plan`之后)。