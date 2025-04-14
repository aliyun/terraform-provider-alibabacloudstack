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
data "alibabacloudstack_ascm_password_pategies" "default" {}  

output "families" {  
  value = data.alibabacloudstack_ascm_password_policies.default.*  
}
```

## 参数说明

支持以下参数：  

* `ids` - (可选) 实例族ID列表。此参数用于过滤结果，仅返回与指定实例族相关的密码策略。  
* `hard_expiry` - 指定密码过期后是否禁用登录。如果设置为 `true`，则密码过期后用户将无法登录。  
* `require_numbers` - 指定密码中是否需要包含数字。如果设置为 `true`，则密码必须包含至少一个数字。  
* `require_symbols` - 指定密码中是否需要包含特殊字符。如果设置为 `true`，则密码必须包含至少一个特殊字符。  
* `require_lowercase_characters` - 指定密码中是否需要包含小写字母。如果设置为 `true`，则密码必须包含至少一个小写字母。  
* `require_uppercase_characters` - 指定密码中是否需要包含大写字母。如果设置为 `true`，则密码必须包含至少一个大写字母。  
* `max_login_attempts` - 允许的最大登录尝试次数。超过此限制后，账户可能会被锁定。  
* `max_password_age` - 密码的有效期（单位：天）。在此期限之后，用户需要更改密码才能继续使用。  
* `minimum_password_length` - 密码的最小长度。密码必须至少达到此长度要求。  
* `password_reuse_prevention` - 允许的最大密码重用尝试次数。此参数用于防止用户重复使用最近使用过的密码。  
* `policies` - (可选) 密码策略列表。此参数允许用户直接指定需要查询的密码策略。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `ids` - 实例族ID列表。返回与查询条件匹配的实例族ID。  
* `policies` - 密码策略列表。每个元素包含以下属性：  
  * `hard_expiry` - 指定密码过期后是否禁用登录。如果为 `true`，则密码过期后用户将无法登录。  
  * `require_numbers` - 指定密码中是否需要包含数字。如果为 `true`，则密码必须包含至少一个数字。  
  * `require_symbols` - 指定密码中是否需要包含特殊字符。如果为 `true`，则密码必须包含至少一个特殊字符。  
  * `require_lowercase_characters` - 指定密码中是否需要包含小写字母。如果为 `true`，则密码必须包含至少一个小写字母。  
  * `require_uppercase_characters` - 指定密码中是否需要包含大写字母。如果为 `true`，则密码必须包含至少一个大写字母。  
  * `max_login_attempts` - 允许的最大登录尝试次数。超过此限制后，账户可能会被锁定。  
  * `max_password_age` - 密码的有效期（单位：天）。在此期限之后，用户需要更改密码才能继续使用。  
  * `minimum_password_length` - 密码的最小长度。密码必须至少达到此长度要求。  
  * `password_reuse_prevention` - 允许的最大密码重用尝试次数。此参数用于防止用户重复使用最近使用过的密码。  
