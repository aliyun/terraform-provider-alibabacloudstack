---
subcategory: "Object Storage Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket"
description: |-
  Provides a resource to create an oss bucket.
---

# alibabacloudstack_oss_bucket

Provides a resource to create an oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.

-> **Note:** Supported credentials: the AK/SK of the target organization (primary account), or an ASCM user under the target organization that uses a properly configured custom role. The ASCM user has the following requirements:

* Bare user AK/SK is not supported: OSS permissions are granted through the role instead of being attached directly to the user, so a properly configured custom role is required.
* Cross-organization accounts (including the environment administrator account, admin) are not supported: OSS requests are rejected with `AccessDenied` ("Put bucket request is not resource owner.", EC 0003-00000202).
* RAM roles are not supported: ASCM cannot resolve the identity of a RAM role session.
* Required permissions: `oss:PutBucket`, `oss:GetBucketAcl`, `oss:PutBucketAcl`, `oss:GetBucketTags`, `oss:PutBucketTags`, `oss:DeleteBucketTags`, `oss:GetBucketLogging`, `oss:PutBucketLogging`, `oss:DeleteBucketLogging`, `oss:GetBucketEncryption`, `oss:PutBucketEncryption`, `oss:DeleteBucketEncryption`, `oss:GetBucketSync`, `oss:PutBucketSync`, `oss:DeleteBucketSync`, `oss:GetBucketStorageCapacity`, `oss:SetBucketStorageCapacity`, `oss:DeleteBucket`, and `ascm:UpdateInstanceBelong` (right after the bucket is created, the provider moves it into the resource set through the ASCM management plane; if that call fails, the created bucket is deleted and the create fails).

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
* `bucket_sync` - (Optional)  Enables or disables bucket synchronization. Defaults to `true`.
* `tags` - (Optional) A mapping of tags assigned to the bucket. The tag `ascm:tag/rg` is managed by the provider and must not be set in `tags`.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.
* `storage_class` -  The storage class of the bucket.
* `vpclist` -  The list of accessible VPCs for the bucket.

## Import

OSS bucket can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_oss_bucket.example <oss_cluster>:<bucket>
```

The `oss_cluster` is the ID of the OSS cluster and `bucket` is the name of the storage bucket.
