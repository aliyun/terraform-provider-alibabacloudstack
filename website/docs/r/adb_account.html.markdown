---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_account"
sidebar_current: "docs-Alibabacloudstack-adb-account"
description: |- 
  Provides a adb Account resource.
---

# alibabacloudstack_adb_account

Provides a adb Account resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccadbaccount96136"
}

variable "password" {
  description = "The password of adb account."
}

resource "alibabacloudstack_adb_account" "default" {
  db_cluster_id       = "am-bp1j43v9c35ef2cvf"
  account_name        = "nametest123"
  account_password    = var.password
  account_type        = "Normal"
  account_description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the ADB cluster where the account belongs. Once set, it cannot be changed.
* `account_name` - (Required, ForceNew) The name of the account. It must start with a letter and can consist of lowercase letters, numbers, and underscores (_). The length should not exceed 16 characters.
* `account_password` - (Optional) The password for the account. It must consist of letters, digits, or underscores, with a length between 6 and 32 characters. You must specify either `account_password` or `kms_encrypted_password`.
* `kms_encrypted_password` - (Optional) An KMS encrypted password used to create or update the database account. If `account_password` is provided, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating the database account. This is valid only when `kms_encrypted_password` is set.
* `account_type` - (Optional, ForceNew) The type of the database account. Default value: `Normal`. Valid values:
  * `Normal`: Standard account. Up to 256 standard accounts can be created for a cluster.
  * `Super`: Privileged account. Only one privileged account can be created for a cluster.
* `account_description` - (Optional) The description of the account. It must start with a Chinese character or an English letter and can include Chinese characters, English letters, underscores (_), hyphens (-), and numbers. The length should be between 2 and 256 characters.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the account. It is composed of the instance ID and the account name in the format `<instance_id>:<account_name>`.