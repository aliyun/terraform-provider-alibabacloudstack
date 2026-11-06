---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_secret"
description: |-
  Provides a Alibabacloudstack Cloud kms secret resource.
---

# alibabacloudstack_kms_secret

This resouce used to create a secret and store its initial version.

-> **Note:** Supported credentials: the AK/SK of the target organization (primary account), or an ASCM user under the target organization that uses a properly configured custom role. The ASCM user has the following requirements:

* Bare user AK/SK is not supported: KMS permissions are granted through the role instead of being attached directly to the user, so a properly configured custom role is required.
* Cross-organization accounts (including the environment administrator account, admin) are not supported: KMS requests are rejected with `Forbidden.NoPermission`.
* RAM roles are not supported: ASCM cannot resolve the identity of a RAM role session.
* Required permissions: `kms:CreateSecret`, `kms:DescribeSecret`, `kms:GetSecretValue`, `kms:UpdateSecret`, `kms:PutSecretValue`, `kms:DeleteSecret`, and, when `tags` is managed, `kms:TagResource` and `kms:UntagResource`.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_kms_secret" "default" {
  secret_name                   = "secret-foo"
  description                   = "from terraform"
  secret_data                   = "Secret data."
  version_id                    = "000000000001"
  force_delete_without_recovery = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the secret.
* `encryption_key_id` - (Optional, ForceNew) The ID of the KMS CMK that is used to encrypt the secret value. If you do not specify this parameter, Secrets Manager automatically creates an encryption key to encrypt the secret.
* `force_delete_without_recovery` - (Optional) Specifies whether to forcibly delete the secret. If this parameter is set to true, the secret cannot be recovered. Valid values: true, false. Default to: false.
* `recovery_window_in_days` - (Optional) Specifies the recovery period of the secret if you do not forcibly delete it. Default value: 30. It will be ignored when `force_delete_without_recovery` is true.
* `secret_data` - (Required) The value of the secret that you want to create. Secrets Manager encrypts the secret value and stores it in the initial version.
* `secret_data_type` - (Optional) The type of the secret value. Valid values: text, binary. Default to "text".
* `secret_name` - (Required, ForceNew) The name of the secret.
* `version_id` - (Required) The version number of the initial version. Version numbers are unique in each secret object.
* `version_stages` - (Optional, List(string)) The stage labels that mark the new secret version. If you do not specify this parameter, Secrets Manager marks it with "ACSCurrent".
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `encryption_key_id` - (Optional, ForceNew) The ID of the KMS CMK that is used to encrypt the secret value. If you do not specify this parameter, Secrets Manager automatically creates an encryption key to encrypt the secret. 

## Attributes Reference

* `id` - The ID of the secret. It same with `secret_name`.
* `arn` - The Alibabacloudstack Resource Name (ARN) of the secret.
* `planned_delete_time` - The time when the secret is scheduled to be deleted.

## Import

alibabacloudstack_kms_secret can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_kms_secret.example <id>
```
