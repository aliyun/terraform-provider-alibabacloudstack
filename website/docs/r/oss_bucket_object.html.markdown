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

-> **Note:** This resource can also be referred to by the following alias: `apsarastack_oss_bucket_object`

## Example Usage

### Adjusting the configuration in provider.tf

```hcl
provider "alibabacloudstack" {
  popgw_domain            = "inter.env205.shuguang.com"
  access_key              = "xxx"
  secret_key              = "xxx"
  region                  = "cn-center-a"
  insecure                = "true"
  resource_group_set_name = "ResourceSet(xxx)"
}
```

### Create a new bucket at first

```hcl
resource "alibabacloudstack_oss_bucket" "example" {
  bucket = "your_bucket_name"
  acl    = "public-read"
}
```

### Uploading a file to a bucket

```hcl
resource "alibabacloudstack_oss_bucket_object" "object-source" {
  bucket  = alibabacloudstack_oss_bucket.example.bucket
  key     = "new_object_key"
  source  = "path/to/file"
}
```

### Uploading content to a bucket

```hcl
resource "alibabacloudstack_oss_bucket_object" "object-content" {
  bucket  = alibabacloudstack_oss_bucket.example.bucket
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```

### Uploading a file with server-side encryption

```hcl
resource "alibabacloudstack_oss_bucket_object" "encrypted-object" {
  bucket                 = alibabacloudstack_oss_bucket.example.bucket
  key                    = "encrypted_object_key"
  source                 = "path/to/file"
  server_side_encryption = "AES256"
}
```

## Argument Reference

-> **Note:** If you specify `content_encoding` you are responsible for encoding the body appropriately (i.e. `source` and `content` both expect already encoded/compressed bytes)

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket to put the file in. Changing this parameter forces a new resource to be created.
* `key` - (Required, ForceNew) The name of the object once it is in the bucket. Changing this parameter forces a new resource to be created.
* `oss_cluster` - (Optional, ForceNew) The OSS cluster identifier. If not specified, the default cluster will be used. Changing this parameter forces a new resource to be created.
* `source` - (Optional) The path to the source file being uploaded to the bucket. Conflicts with `content`. Either `source` or `content` must be provided.
* `content` - (Optional) The literal content being uploaded to the bucket. Conflicts with `source`. Either `source` or `content` must be provided.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/52284.htm) to apply. Valid values: `private`, `public-read`, `public-read-write`. Defaults to `private`.
* `content_type` - (Optional, Computed) A standard MIME type describing the format of the object data, e.g. `application/octet-stream`. All valid MIME types are valid for this input.
* `cache_control` - (Optional) Specifies caching behavior along the request/reply chain. Read [RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_disposition` - (Optional) Specifies presentational information for the object. Read [RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_encoding` - (Optional) Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_md5` - (Optional, Computed) The MD5 value of the content. Read [MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm) for computing method.
* `expires` - (Optional) Specifies expire date for the request/response. Must be in RFC1123 format. Read [RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `server_side_encryption` - (Optional) Specifies server-side encryption of the object in OSS. Valid values: `AES256`, `KMS`.
* `kms_key_id` - (Optional, Computed) Specifies the primary key managed by KMS. This parameter is valid when the value of `server_side_encryption` is set to `KMS`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource formatted as `{oss_cluster}:{bucket}:{key}`.
* `version_id` - A unique version ID value for the object, if bucket versioning is enabled.
* `content_type` - The MIME type of the object data.
* `content_md5` - The MD5 hash of the object content.
* `kms_key_id` - The KMS key ID used for server-side encryption (when `server_side_encryption` is `KMS`).

## Import

OSS Bucket Object can be imported using the composite ID in the format `{oss_cluster}:{bucket}:{key}`, e.g.

```shell
$ terraform import alibabacloudstack_oss_bucket_object.example default:my-bucket:my-object-key
```