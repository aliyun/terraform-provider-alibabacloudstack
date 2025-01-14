---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket"
sidebar_current: "docs-alibabacloudstack-resource-oss-bucket"
description: |-
  Provides a resource to create an oss bucket.
---

# alibabacloudstack\_oss\_bucket

Provides a resource to create an oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.


## Example Usage

Private Bucket

```
resource "alibabacloudstack_oss_bucket" "demo" {
  bucket = "sample_bucket"
  acl    = "public-read"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional) It Can be "private", "public-read" and "public-read-write". Defaults to "private".
* `logging` - (Optional) The logging object supports the following:
  * `target_bucket` - (Required) The name of the bucket that will receive the log objects.
  * `target_prefix` - (Optional) To specify a key prefix for log objects. 
* `storage_class` - (Optional, ForceNew) Object storage type. Possible values: `Standard`, `IA` and `Archive`.
* `vpclist` - (Optional) List of accessible VPCs
* `storage_capacity` - (Optional) Sets a capacity limit on the bucket. If the capacity limit is reached, write operations. unit GB
* `sse_algorithm` - (Optional) Encrypts files uploaded to OSS. It Can be "", "AES256", "SM4", and "KMS". Defaults to "private".
* `kms_key_id` - (Optional) Required when sse_algorithm is KMS.



## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.

