---
subcategory: "OSS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_oss_bucket"
sidebar_current: "docs-apsarastack-resource-oss-bucket"
description: |-
  Provides a resource to create an oss bucket.
---

# apsarastack\_oss\_bucket

Provides a resource to create an oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.


## Example Usage

Private Bucket

```
resource "apsarastack_oss_bucket" "demo" {
  bucket = "sample_bucket"
  acl    = "public-read"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional) It Can be "private", "public-read" and "public-read-write". Defaults to "private".
* `logging` - (Optional) The logging object supports the following:
    - `target_bucket` - (Required) The name of the bucket that will receive the log objects.
    - `target_prefix` - (Optional) To specify a key prefix for log objects. 

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.


