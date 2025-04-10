---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_password_policy"
sidebar_current: "docs-alibabacloudstack-resource-ascm-password-policy"
description: |-
   编排Ascm密码策略
---
# alibabacloudstack_ascm_password_policy

使用Provider配置的凭证在指定的资源集下编排Ascm密码策略。

## 示例用法

```
resource "alibabacloudstack_ascm_password_policy" "default"{
  minimum_password_length = 20
  max_login_attempts      = 8
  hard_expiry             = true
}
```

## 参数说明

支持以下参数：

* `hard_expiry` - (可选) 指定在密码过期后是否禁用登录。设置为 `true` 表示禁用登录，设置为 `false` 表示不禁用登录。
* `require_numbers` - (可选) 指定密码中是否需要包含数字。默认值为 `false`。
* `require_symbols` - (可选) 指定密码中是否需要包含特殊字符。默认值为 `false`。
* `require_lowercase_characters` - (可选) 指定密码中是否需要包含小写字母。默认值为 `false`。
* `require_uppercase_characters` - (可选) 指定密码中是否需要包含大写字母。默认值为 `false`。
* `max_login_attempts` - (可选) 允许的最大登录尝试次数。超过此次数后，账户将被锁定。有效值范围：[1-10]。
* `max_password_age` - (可选) 密码的有效期（以天为单位）。密码在此期限后将过期。有效值范围：[0-365]，其中 `0` 表示密码永不过期。
* `minimum_password_length` - (可选) 密码的最小长度。有效值范围：[8-32]。
* `password_reuse_prevention` - (可选) 允许的最大密码重用次数。有效值范围：[0-24]，其中 `0` 表示不允许重用任何之前的密码。

## 属性说明

导出以下属性：

* `id` - 资源的唯一标识符。
* `policy_name` - 密码策略名称。
* `policy_details` - 密码策略的具体配置信息。