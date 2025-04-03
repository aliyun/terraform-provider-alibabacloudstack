---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_buckets"
sidebar_current: "docs-alibabacloudstack-datasource-oss-buckets"
description: |-
    Provides a list of OSS buckets to the user.
---

# alibabacloudstack_oss_buckets

This data source provides the OSS buckets of the current AlibabacloudStack Cloud user.

## Example Usage

```
data "alibabacloudstack_oss_buckets" "oss_buckets_ds" {
  name_regex = "sample_oss_bucket"
}

output "first_oss_bucket_name" {
  value = "${data.alibabacloudstack_oss_buckets.oss_buckets_ds.buckets.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by bucket name.
* `ids` - (Optional) A list of Bucket IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of bucket names. 
* `buckets` - A list of buckets. Each element contains the following attributes:
  * `name` - Bucket name.
  * `acl` - Bucket access control list. Possible values: `private`, `public-read` and `public-read-write`.
  * `extranet_endpoint` - Internet domain name for accessing the bucket from outside.
  * `intranet_endpoint` - Intranet domain name for accessing the bucket from an ECS instance in the same region.
  * `location` - Region of the data center where the bucket is located.
  * `owner` - Bucket owner.
  * `storage_class` - Object storage type. Possible values: `Standard`, `IA` and `Archive`.
  * `creation_date` - Bucket creation date.
  * `id` - Bucket ID.