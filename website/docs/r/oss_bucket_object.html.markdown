---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_bucket_object"
sidebar_current: "docs-alibabacloudstack-resource-oss-bucket-object"
description: |-
  Provides a resource to create a oss bucket object.
---

# alibabacloudstack_oss_bucket_object

Provides a resource to put a object(content or file) to a oss bucket.

## Example Usage

### Adjusting the configuration in provider.tf"
> Since operations on objects are facilitated by the OSS Data Gateway,
> it is necessary to configure a separate address for the Data Gateway within the provider configuration in provider.tf.

```
provider "alibabacloudstack" {
  #popgw_domain            = "inter.env205.shuguang.com"
  access_key              = "xxx"
  ...
  insecure                = "true"
  resource_group_set_name = "xxx"
  ossservice_domain = "<oss data endpoint>"
}
```

### Create a new bucket at first
```
resource "alibabacloudstack_oss_bucket" "example" {
  bucket = "your_bucket_name"
  acl    = "public-read"
}
```

### Granting OSS permissions to specified users
> currently only supported via the ascm.

### Uploading a file to a bucket

```
resource "alibabacloudstack_oss_bucket_object" "object-source" {
  bucket  = "${alibabacloudstack_oss_bucket.example.bucket}"
  key    = "new_object_key"
  source = "path/to/file"
}
```

### Uploading a content to a bucket

```
resource "alibabacloudstack_oss_bucket_object" "object-content" {
  bucket  = "${alibabacloudstack_oss_bucket.example.bucket}"
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```

## Argument Reference

-> **Note:** If you specify `content_encoding` you are responsible for encoding the body appropriately (i.e. `source` and `content` both expect already encoded/compressed bytes)

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to put the file in.
* `key` - (Required) The name of the object once it is in the bucket.
* `source` - (Optional) The path to the source file being uploaded to the bucket.
* `content` - (Optional unless `source` given) The literal content being uploaded to the bucket.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/52284.htm) to apply. Defaults to "private".
* `content_type` - (Optional) A standard MIME type describing the format of the object data, e.g. application/octet-stream. All Valid MIME Types are valid for this input.
* `cache_control` - (Optional) Specifies caching behavior along the request/reply chain. Read [RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_disposition` - (Optional) Specifies presentational information for the object. Read [RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_encoding` - (Optional) Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_md5` - (Optional) The MD5 value of the content. Read [MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm) for computing method.
* `expires` - (Optional) Specifies expire date for the the request/response. Read [RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `server_side_encryption` - (Optional) Specifies server-side encryption of the object in OSS. Valid values are `AES256`, `KMS`. Default value is `AES256`.
* `kms_key_id` - (Optional, Available in 1.62.1+) Specifies the primary key managed by KMS. This parameter is valid when the value of `server_side_encryption` is set to KMS.
* `version_id` - (Optional)  A unique version ID value for the object, if bucket versioning is enabled.

Either `source` or `content` must be provided to specify the bucket content.
These two arguments are mutually-exclusive.

## Attributes Reference

The following attributes are exported

* `id` - the `key` of the resource supplied above.
* `version_id` -  A unique version ID value for the object, if bucket versioning is enabled.