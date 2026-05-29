---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_key"
sidebar_current: "docs-Alibabacloudstack-resource-kms-key"
description: |-
  Provides a Alibabacloudstack KMS key resource.
---

# alibabacloudstack_kms_key

> **Note:** This resource can also be referred to by the following aliases:
> - `apsarastack_kms_key`

A KMS key can help users to protect data security during transmission.



## Example Usage

Basic Usage

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
  key_state               = "Enabled"
}
```
## Argument Reference

The following arguments are supported:

* `key_usage` - (Optional, ForceNew) Specifies the usage of the CMK. Valid values: `ENCRYPT/DECRYPT`, `SIGN/VERIfy`. Default: `ENCRYPT/DECRYPT`.
* `origin` - (Optional, ForceNew) The source of the key material for the CMK. Valid values: `Aliyun_KMS`, `EXTERNAL`. Default: `Aliyun_KMS`.
* `protection_level` - (Optional, ForceNew) The protection level of the CMK. Valid values: `SOFTWARE`, `HSM`. Default: `SOFTWARE`.
* `description` - (Optional) The description of the key as viewed in Alibabacloudstack console.
* `automatic_rotation` - (Optional) Specifies whether to enable automatic key rotation. Valid values: `Enabled`, `Disabled`. Default: `Disabled`.
* `key_state` - (Optional) The status of CMK. Valid values: `Enabled`, `Disabled`, `PendingDeletion`. Default: `Enabled`.
* `pending_window_in_days` - (Optional) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days (inclusive). Default: `7`.
* `rotation_interval` - (Optional) The period of automatic key rotation. The format is a number followed by a time unit (d for days, h for hours, m for minutes, s for seconds), such as `7d` or `2678400s`. Only valid when `automatic_rotation` is set to `Enabled`.
* `is_enabled` - (Optional, Deprecated) Field `is_enabled` has been deprecated from provider version 1.85.0. Use `key_state` instead.
* `deletion_window_in_days` - (Optional, Deprecated) Field `deletion_window_in_days` has been deprecated from provider version 1.85.0. Use `pending_window_in_days` instead.

-> **NOTE:** When the pre-deletion window elapses, the key is permanently deleted and cannot be recovered.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the key.
* `arn` - The Alibabacloudstack Resource Name (ARN) of the key.
* `creation_date` - The date and time when the CMK was created. The time is displayed in UTC.
* `creator` - The creator of the CMK.
* `delete_date` - The scheduled date to delete CMK. The time is displayed in UTC. This value is returned only when the KeyState value is PendingDeletion.
* `key_state` - The status of the CMK.
* `key_usage` - The usage of the CMK.
* `origin` - The source of the key material for the CMK.
* `protection_level` - The protection level of the CMK.
* `primary_key_version` - The ID of the current primary key version of the symmetric CMK.
* `automatic_rotation` - Indicates whether automatic rotation is enabled for the key.
* `description` - The description of the key.
* `last_rotation_date` - The date and time the last rotation was performed. The time is displayed in UTC.
* `material_expire_time` - The time and date the key material for the CMK expires. The time is displayed in UTC. If the value is empty, the key material for the CMK does not expire.
* `next_rotation_date` - The time the next rotation is scheduled for execution.
* `rotation_interval` - The period of automatic key rotation.

## Import

KMS Key can be imported using the key ID, e.g.

```
$ terraform import alibabacloudstack_kms_key.example key-12345678
``` 