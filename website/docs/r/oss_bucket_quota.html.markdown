---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_quota"
sidebar_current: "docs-Alibabacloudstack-oss-bucket-quota"
description: |-
  Provides a OSS Bucket Quota resource.
---

# alibabacloudstack_oss_bucket_quota

Provides a OSS Bucket Quota resource.

## Example Usage

```hcl
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "my-tf-test-bucket"
}

resource "alibabacloudstack_oss_bucket_quota" "default" {
  bucket      = alibabacloudstack_oss_bucket.default.bucket
  oss_cluster = "your-oss-cluster"  # Optional: OSS cluster name
  quota       = 10240  # Quota in MB
}
```

## Argument Reference
The following arguments are supported:

* `bucket` - (Required, ForceNew) - The name of the OSS bucket.
* `oss_cluster` - (Optional, ForceNew) - The name of the OSS cluster. If not specified, the default cluster will be used.
* `quota` - (Required, ForceNew) - The storage quota for the OSS bucket in megabytes (MB).


## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `bucket` - The name of the OSS bucket.
* `oss_cluster` - The name of the OSS cluster.
* `quota` - The storage quota for the OSS bucket in megabytes (MB).
