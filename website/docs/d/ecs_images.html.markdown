---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_images"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-images"
description: |- 
  Provides a list of ecs images owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_images
-> **NOTE:** Alias name has: `alibabacloudstack_images`

This data source provides a list of ECS images in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_images" "images_ds" {
  owners     = "system"
  name_regex = "^centos_6"
}

output "first_image_id" {
  value = "${data.alibabacloudstack_ecs_images.images_ds.images.0.image_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting images by name.
* `most_recent` - (Optional, ForceNew, type: bool) If more than one result is returned, select the most recent one.
* `owners` - (Optional, ForceNew) Filter results by a specific image owner. Valid items are `system`, `self`, `others`, `marketplace`.

-> **NOTE:** At least one of the `name_regex`, `most_recent`, and `owners` must be set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of image IDs.
* `images` - A list of images. Each element contains the following attributes:
  * `id` - The ID of the image (same as `image_id`).
  * `image_id` - The unique identifier for the image.
  * `architecture` - Platform type of the image system: `i386` or `x86_64`.
  * `creation_time` - The time when the image was created.
  * `description` - The description of the image.
  * `image_owner_alias` - Alias of the image owner.
  * `os_name` - The Chinese display name of the operating system.
  * `os_name_en` - The English display name of the operating system.
  * `os_type` - The operating system type of the image. Valid values: `windows` and `linux`.
  * `platform` - The operating system platform of the image.
  * `status` - Status of the image. Possible values: `UnAvailable`, `Available`, `Creating`, and `CreateFailed`.
  * `state` - State of the image (same as `status`).
  * `size` - The size of the image in GiB.
  * `disk_device_mappings` - Snapshot information for the image. Each mapping includes:
    * `device` - Device information of the created disk, such as `/dev/xvdb`.
    * `size` - Size of the created disk in GiB.
    * `snapshot_id` - The snapshot ID associated with the disk.
  * `product_code` - Product code of the image on the image market.
  * `is_subscribed` - Whether the user has subscribed to the terms of service for the image product corresponding to the `product_code`.
  * `is_copied` - Indicates whether it is a copied image. Valid values: `true` or `false`.
  * `is_self_shared` - Indicates whether the custom image has been shared with other users. Valid values: `true` or `false`.
  * `image_version` - The version of the image.
  * `progress` - The progress of image creation, presented in percentages.
  * `usage` - Specifies whether to check the validity of the request without actually making the request.
  * `is_support_io_optimized` - Indicates whether the image can be used on I/O optimized instances.
  * `tags` - The tags of the resource.
  * `name` - Name of the image.