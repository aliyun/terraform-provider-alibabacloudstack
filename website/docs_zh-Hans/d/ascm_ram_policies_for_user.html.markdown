---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_ram_policies_for_user"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-policies-for-user"
description: |-
    查询用户的RAM策略列表
---

# alibabacloudstack_ascm_ram_policies_for_user

根据指定过滤条件列出当前凭证权限可以访问的特定用户RAM策略列表。

## 示例用法

```
data "alibabacloudstack_ascm_ram_policies_for_user" "default" {
  login_name = "test_admin"
}
output "ramPolicy" {
  value = data.alibabacloudstack_ascm_ram_policies_for_user.default.*
}

```

## 参数说明

以下参数受支持：

* `ids` - (可选) RAM策略ID的列表。
* `login_name` - (必填，变更时重建) 用户的登录名。

## 属性参考

除了上述参数外，还导出以下属性：

* `policies` - 策略列表。每个元素包含以下属性：
  * `policy_name` - 策略名称。
  * `description` - 关于策略的描述。
  * `attach_date` - RAM策略的创建日期。
  * `policy_type` - 策略类型。
  * `default_version` - 默认版本。
  * `policy_document` - 策略文档。