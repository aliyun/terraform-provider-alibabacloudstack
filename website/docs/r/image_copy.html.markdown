---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_copy"
sidebar_current: "docs-alibabacloudstack-resource-image-copy"
description: |-
  Provides an ECS image copy resource.
---

# alibabacloudstack_image_copy

Copies a custom image from one region to another. You can use copied images to perform operations in the target region, such as creating instances (RunInstances) and replacing system disks (ReplaceSystemDisk).

-> **NOTE:** You can only copy the custom image when it is in the Available state.

-> **NOTE:** You can only copy the image belonging to your Alibabacloudstack Cloud account. Images cannot be copied from one account to another.

-> **NOTE:** If the copying is not completed, you cannot call DeleteImage to delete the image but you can call CancelCopyImage to cancel the copying.

## Example Usage

```
resource "alibabacloudstack_image_copy" "default" {
  source_image_id    = "m-bp1gxyhdswlsn18tu***"
  source_region_id   = "cn-hangzhou"
  image_name         = "test-image"
  description        = "test-image"
  tags               = {
         FinanceDept = "FinanceDeptJoshua"
     }
}
```

## Argument Reference

The following arguments are supported:

* `source_image_id` - (Required, ForceNew) The source image ID.
* `destination_region_id` - (Required, ForceNew) The Target Region ID to Copy.
* `name` - (Optional,Deprecated) Field 'name' has been deprecated. New field 'image_name' instead.
* `image_name` - (Optional) The image name. It must be 2 to 128 characters in length, and must begin with a letter or Chinese character (beginning with http:// or https:// is not allowed). It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of the image. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `kms_key_id` - (Optional, ForceNew) The KMS key ID used for encryption. 
* `encrypted` - (Optional, ForceNew) Indicates whether the image is encrypted. 
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference
 
The following attributes are exported:
 
* `id` - ID of the image.
* `name` - The name of the image. 
* `image_name` - The name of the image. 
* `description` - The description of the image.