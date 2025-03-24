---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_keys"
sidebar_current: "docs-alibabacloudstack-datasource-kms-keys"
description: |-
    Provides a list of available KMS Keys.
---

# alibabacloudstack_kms_keys

This data source provides a list of KMS keys in an Alibabacloudstack Cloud account according to the specified filters.

## Example Usage

```
# Declare the data source
data "alibabacloudstack_kms_keys" "kms_keys_ds" {
  description_regex = "Hello KMS"
}

output "first_key_id" {
  value = "${data.alibabacloudstack_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of KMS key IDs.
* `description_regex` - (Optional) A regex string to filter the results by the KMS key description.
* `status` - (Optional) Filter the results by status of the KMS keys. Valid values: `Enabled`, `Disabled`, `PendingDeletion`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of KMS key IDs.
* `keys` - A list of KMS keys. Each element contains the following attributes:
  * `id` - ID of the key.
  * `arn` - The alibabacloudstack Cloud Resource Name (ARN) of the key.
  * `description` - Description of the key.
  * `status` - Status of the key. Possible values: `Enabled`, `Disabled` and `PendingDeletion`.
  * `creation_date` - Creation date of key.
  * `delete_date` - Deletion date of key.
  * `creator` - The owner of the key.
  * `computed_property` - Indicates a computed property of the key.