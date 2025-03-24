---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_kms"
sidebar_current: "docs-Alibabacloudstack-oss-bucket-kms"
description: |-
  Provides a OSS Bucket KMS resource.
---

# alibabacloudstack_oss_bucket_kms

**Note:** This resource is deprecated. The `oss_bucket` resource already includes the corresponding functions for managing server-side encryption with KMS.

Provides a OSS Bucket KMS resource.

## Example Usage

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_kms" "default" {
  bucket          = alibabacloudstack_oss_bucket.default.bucket
  sse_algorithm   = "KMS"
  kms_data_encryption = "AES256"
  kms_master_key_id = "your_kms_master_key_id"
}
```

## Argument Reference
The following arguments are supported:

* `bucket` - (Required, ForceNew) - The name of the OSS bucket.
* `sse_algorithm` - (Required, ForceNew) - The server-side encryption algorithm. Valid values: KMS.
* `kms_data_encryption` - (Optional, ForceNew) - The data encryption algorithm used by KMS. Valid values: AES256.
* `kms_master_key_id` - (Optional, ForceNew) - The ID of the KMS master key used for encryption.

## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `bucket` - The name of the OSS bucket.
* `sse_algorithm` - The server-side encryption algorithm.
* `kms_data_encryption` - The data encryption algorithm used by KMS.
* `kms_master_key_id` - The ID of the KMS master key used for encryption.