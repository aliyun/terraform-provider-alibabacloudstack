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
  value = data.alibabacloudstack_ascm_ram_pategies_for_user.default.*
}

```

## 参数说明

以下参数受支持：

* `ids` - (可选) RAM策略ID的列表。可以通过此参数筛选出特定的策略。
* `login_name` - (必填，变更时重建) 用户的登录名。用于指定需要查询其RAM策略的用户。

## 属性参考

除了上述参数外，还导出以下属性：

* `policies` - 策略列表。每个元素包含以下属性：
  * `policy_name` - 策略名称。表示该RAM策略的唯一标识名称。
  * `description` - 关于策略的描述。提供有关策略功能或用途的详细信息。
  * `attach_date` - 策略绑定到用户的日期。格式为标准时间格式（如：2023-01-01T12:00:00Z）。
  * `policy_type` - 策略类型。表示策略是自定义策略还是系统预定义策略，例如“Custom”或“System”。
  * `default_version` - 默认版本。表示策略文档的默认版本号。
  * `policy_document` - 策略文档。以JSON格式表示的策略内容，定义了具体的权限规则和范围。