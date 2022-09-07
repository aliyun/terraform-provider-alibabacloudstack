---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_password_policy"
sidebar_current: "docs-alibabacloudstack-resource-ascm-password-policy"
description: |-
   Provides a Ascm password policy configuration.
---
# alibabacloudstack\_ascm_password_policy

Provides an Ascm Password Policy resource.

## Example Usage

```
resource "alibabacloudstack_ascm_password_policy" "default"{
  minimum_password_length = 20
  max_login_attempts      = 8
  hard_expiry             = true
}

```
## Argument Reference

The following arguments are supported:

* `hard_expiry` - (Optional) Specifies whether to disable logon after the password expires.
* `require_numbers` - (Optional) Specifies whether digits are required.
* `require_symbols` - (Optional) Specifies whether special characters are required.
* `require_lowercase_characters` - (Optional)  Specifies whether lowercase letters are required.
* `require_uppercase_characters` - (Optional)  Specifies whether uppercase letters are required.
* `max_login_attempts` - (Optional) The maximum number of allowed logon attempts
* `max_password_age` - (Optional) The validity period of the password.
* `minimum_password_length` - (Optional) The minimum length of the password.Valid value range: [8-32].
* `password_reuse_prevention` - (Optional) The maximum number of allowed password reuse attempts.

