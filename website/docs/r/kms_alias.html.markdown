---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_alias"
description: |-
  Provides a AlibabacloudStack KMS Alias resource.
---

# alibabacloudstack_kms_alias

Create an alias for the master key (CMK).

-> **Note:** Supported credentials: the AK/SK of the target organization (primary account), or an ASCM user under the target organization that uses a properly configured custom role. The ASCM user has the following requirements:

* Bare user AK/SK is not supported: KMS permissions are granted through the role instead of being attached directly to the user, so a properly configured custom role is required.
* Cross-organization accounts (including the environment administrator account, admin) are not supported: KMS requests are rejected with `Forbidden.NoPermission`.
* RAM roles are not supported: ASCM cannot resolve the identity of a RAM role session.
* Required permissions: `kms:CreateAlias`, `kms:DescribeKey`, `kms:UpdateAlias`, `kms:DeleteAlias`.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_kms_key" "key" {}

resource "alibabacloudstack_kms_alias" "alias" {
  alias_name = "alias/test_kms_alias"
  key_id     = alibabacloudstack_kms_key.key.id
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) The alias of CMK. `Encrypt`, `GenerateDataKey`, `DescribeKey` can be called using aliases. Length of characters other than prefixes: minimum length of 1 character and maximum length of 255 characters. Must contain prefix `alias/`.
* `key_id` - (Required) The id of the key.


-> **NOTE:** Each alias represents only one master key(CMK).

-> **NOTE:** Within an area of the same user, alias is not reproducible.

-> **NOTE:** UpdateAlias can be used to update the mapping relationship between alias and master key(CMK).


## Attributes Reference

* `id` - The ID of the alias.

## Import

KMS Alias can be imported using the alias_name, e.g.

```
$ terraform import alibabacloudstack_kms_alias.example alias/test_kms_alias
```