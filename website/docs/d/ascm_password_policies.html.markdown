---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_password_policies"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-password-policies"
description: |-
    Provides a list of password policies to the user.
---

# alibabacloudstack_ascm_password_policies

This data source provides the password policies of the current Apsara Stack Cloud user.

## Example Usage

```
data "alibabacloudstack_ascm_password_policies" "default" {}

output "families" {
  value = data.alibabacloudstack_ascm_password_policies.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance family IDs.
* `hard_expiry` - Specifies whether to disable logon after the password expires.
* `require_numbers` - Specifies whether digits are required.
* `require_symbols` - Specifies whether special characters are required.
* `require_lowercase_characters` - Specifies whether lowercase letters are required.
* `require_uppercase_characters` - Specifies whether uppercase letters are required.
* `max_login_attempts` - The maximum number of allowed logon attempts.
* `max_password_age` - The validity period of the password. Unit: days.
* `minimum_password_length` - The minimum length of the password.
* `password_reuse_prevention` - The maximum number of allowed password reuse attempts.
* `policies` - (Optional) A list of password policies.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance family IDs.
* `policies` - A list of password policies. Each element contains the following attributes:
  * `hard_expiry` - Specifies whether to disable logon after the password expires.
  * `require_numbers` - Specifies whether digits are required.
  * `require_symbols` - Specifies whether special characters are required.
  * `require_lowercase_characters` - Specifies whether lowercase letters are required.
  * `require_uppercase_characters` - Specifies whether uppercase letters are required.
  * `max_login_attempts` - The maximum number of allowed logon attempts.
  * `max_password_age` - The validity period of the password. Unit: days.
  * `minimum_password_length` - The minimum length of the password.
  * `password_reuse_prevention` - The maximum number of allowed password reuse attempts.