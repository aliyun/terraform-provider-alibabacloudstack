---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image"
sidebar_current: "docs-alibabacloudstack-resource-image"
description: |-
  Provides an ECS image resource.
---

# alibabacloudstack\_image

Creates a custom image. You can then use a custom image to create ECS instances (RunInstances) or change the system disk for an existing instance (ReplaceSystemDisk).

## Example Usage

```
resource "alibabacloudstack_image" "default" {
  instance_id        = "i-bp1g6zv0ce8oghu7k***"
  image_name         = "test-image"
  description        = "test-image"
  tags = {
    FinanceDept = "FinanceDeptJoshua"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew, Conflict with `snapshot_id ` and `disk_device_mapping `) The instance ID.
* `image_name` - (Optional) The image name. It must be 2 to 128 characters in length, and must begin with a letter or Chinese character (beginning with http:// or https:// is not allowed). It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of the image. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `snapshot_id` - (Optional, ForceNew, Conflict with `instance_id ` and `disk_device_mapping `) Specifies a snapshot that is used to create a custom image.
* `tags` - (Optional) The tag value of an image. The value of N ranges from 1 to 20.
* `disk_device_mapping` - (Optional, ForceNew, Conflict with `snapshot_id ` and `instance_id `) Description of the system with disks and snapshots under the image.
  * `size` - (Optional, ForceNew) Specifies the size of a disk in the combined custom image, in GiB. Value range: 5 to 2000.
  * `snapshot_id` - (Optional, ForceNew) Specifies a snapshot that is used to create a combined custom image.
* `force` - (Optional) Indicates whether to force delete the custom image, Default is `false`. 
  - true：Force deletes the custom image, regardless of whether the image is currently being used by other instances.
  - false：Verifies that the image is not currently in use by any other instances before deleting the image.
   
### Timeouts

* `create` - (Defaults to 10 mins) Used when creating the image (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the image.
   
   
### Attributes Reference
 
 The following attributes are exported:
 
* `id` - ID of the image.
 