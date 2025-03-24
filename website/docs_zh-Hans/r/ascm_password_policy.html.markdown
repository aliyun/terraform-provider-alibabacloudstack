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

## 参数参考

支持以下参数：

* `hard_expiry` - (可选) 指定在密码过期后是否禁用登录。
* `require_numbers` - (可选) 指定是否需要数字。
* `require_symbols` - (可选) 指定是否需要特殊字符。
* `require_lowercase_characters` - (可选) 指定是否需要小写字母。
* `require_uppercase_characters` - (可选) 指定是否需要大写字母。
* `max_login_attempts` - (可选) 允许的最大登录尝试次数。
* `max_password_age` - (可选) 密码的有效期。
* `minimum_password_length` - (可选) 密码的最小长度。有效值范围：[8-32]。
* `password_reuse_prevention` - (可选) 允许的最大密码重用尝试次数。

## 属性参考

导出以下属性：