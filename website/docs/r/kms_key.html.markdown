---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_key"
sidebar_current: "docs-alibabacloudstack-resource-kms-key"
description: |-
  Provides a Alibabacloudstackms key resource.
---

# alibabacloudstack_kms_key

A kms key can help user to protect data security in the transmission process.



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

* `description` - (Optional) The description of the key as viewed in Alibabacloudstack console.
* `key_usage` - (Optional, ForceNew) Specifies the usage of CMK. Currently, default to 'ENCRYPT/DECRYPT', indicating that CMK is used for encryption and decryption.
* `automatic_rotation` - (Optional) Specifies whether to enable automatic key rotation. Default:"Disabled".
* `is_enabled` - (Optional,Deprecated) Field 'is_enabled' has been deprecated. New field 'key_state' instead.
* `key_state` - (Optional) The status of CMK. Defaults to Enabled.
* `origin` - (Optional, ForceNew) The source of the key material for the CMK. Defaults to "Aliyun_KMS".
* `deletion_window_in_days` - (Optional,Deprecated) Field 'deletion_window_in_days' has been deprecated. New field 'pending_window_in_days' instead.
* `pending_window_in_days` - (Optional) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 30 days.
* `protection_level` - (Optional, ForceNew) The protection level of the CMK. Defaults to "SOFTWARE".
* `rotation_interval` - (Optional) The period of automatic key rotation. Unit: seconds. 
* `rotation_interval` - (Optional) The period of automatic key rotation. Unit: seconds.

-> **NOTE:** When the pre-deletion days elapses, the key is permanently deleted and cannot be recovered.


## Attributes Reference

* `id` - The ID of the key.
* `arn` - The Alibabacloudstack Resource Name (ARN) of the key.
* `creation_date` -The date and time when the CMK was created. The time is displayed in UTC.
* `creator` -The creator of the CMK.
* `delete_date` -The scheduled date to delete CMK. The time is displayed in UTC. This value is returned only when the KeyState value is PendingDeletion.
* `last_rotation_date` - The date and time the last rotation was performed. The time is displayed in UTC. 
* `material_expire_time` - The time and date the key material for the CMK expires. The time is displayed in UTC. If the value is empty, the key material for the CMK does not expire.
* `next_rotation_date` - The time the next rotation is scheduled for execution. 
* `primary_key_version` - The ID of the current primary key version of the symmetric CMK.