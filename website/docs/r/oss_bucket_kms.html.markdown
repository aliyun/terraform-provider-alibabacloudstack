---
subcategory: "Object Storage Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_kms"
description: |-
  Provides a OSS Bucket KMS resource.
---

# alibabacloudstack_oss_bucket_kms

**Note:** This resource is deprecated and is scheduled for removal in version 3.21.0. The `oss_bucket` resource already includes the corresponding functions for managing server-side encryption with KMS.

Provides a OSS Bucket KMS resource.

## Example Usage

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_kms" "default" {
  bucket          = alibabacloudstack_oss_bucket.default.bucket
  oss_cluster     = "your-oss-cluster-id"
  sse_algorithm   = "KMS"
  kms_master_key_id = "your_kms_master_key_id"
}
```

## Argument Reference
The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the OSS bucket.
* `oss_cluster` - (Optional, ForceNew) The ID of the OSS cluster.
* `sse_algorithm` - (Required, ForceNew) The server-side encryption algorithm. Valid value: `KMS`.
* `kms_master_key_id` - (Optional, ForceNew) The ID of the KMS master key used for encryption. This parameter is required when `sse_algorithm` is set to `KMS`.

## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `bucket` - The name of the OSS bucket.
* `oss_cluster` - The ID of the OSS cluster.
* `sse_algorithm` - The server-side encryption algorithm.
* `kms_master_key_id` - The ID of the KMS master key used for encryption.

## Import

OSS Bucket KMS can be imported using the `oss_cluster:bucket_name`, e.g.

```
$ terraform import alibabacloudstack_oss_bucket_kms.example oss-cluster-id:my-tf-test-bucket
```