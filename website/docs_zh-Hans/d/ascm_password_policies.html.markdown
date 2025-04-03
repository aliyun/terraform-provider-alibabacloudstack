---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_password_policies"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-password-policies"
description: |-
    查询用户密码策略

---

# alibabacloudstack_ascm_password_policies

根据指定过滤条件列出当前凭证权限可以访问的密码策略列表。

## 示例用法

```
data "alibabacloudstack_ascm_password_policies" "default" {}

output "families" {
  value = data.alibabacloudstack_ascm_password_policies.default.*
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 实例族ID列表。
* `hard_expiry` - 指定密码过期后是否禁用登录。
* `require_numbers` - 指定是否需要数字。
* `require_symbols` - 指定是否需要特殊字符。
* `require_lowercase_characters` - 指定是否需要小写字母。
* `require_uppercase_characters` - 指定是否需要大写字母。
* `max_login_attempts` - 允许的最大登录尝试次数。
* `max_password_age` - 密码的有效期。单位：天。
* `minimum_password_length` - 密码的最小长度。
* `password_reuse_prevention` - 允许的最大密码重用尝试次数。
* `policies` - (可选) 密码策略列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 实例族ID列表。
* `policies` - 密码策略列表。每个元素包含以下属性：
  * `hard_expiry` - 指定密码过期后是否禁用登录。
  * `require_numbers` - 指定是否需要数字。
  * `require_symbols` - 指定是否需要特殊字符。
  * `require_lowercase_characters` - 指定是否需要小写字母。
  * `require_uppercase_characters` - 指定是否需要大写字母。
  * `max_login_attempts` - 允许的最大登录尝试次数。
  * `max_password_age` - 密码的有效期。单位：天。
  * `minimum_password_length` - 密码的最小长度。
  * `password_reuse_prevention` - 允许的最大密码重用尝试次数。