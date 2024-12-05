---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_images"
sidebar_current: "docs-alibabacloudstack-datasource-images"
description: |-
    Provides a list of images available to the user.
---

# alibabacloudstack\_images

This data source provides available image resources. It contains user's private images, system images provided by Alibabacloudstack Cloud, 
other public images and the ones available on the image market. 

## Example Usage

```
data "alibabacloudstack_images" "images_ds" {
  owners     = "system"
  name_regex = "^centos_6"
}

output "first_image_id" {
  value = "${data.alibabacloudstack_images.images_ds.images.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting images by name. 
* `most_recent` - (Optional, type: bool) If more than one result are returned, select the most recent one.
* `owners` - (Optional) Filter results by a specific image owner. Valid items are `system`, `self`, `others`, `marketplace`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

-> **NOTE:** At least one of the `name_regex`, `most_recent` and `owners` must be set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of image IDs.
* `images` - A list of images. Each element contains the following attributes:
  * `id` - ID of the image.
  * `name` - name of the image.
  * `image_id` - Alias of the key `id`.
  * `architecture` - Platform type of the image system: i386 or x86_64.
  * `state` - state of the image.
  * `creation_time` - Time of creation.
  * `description` - Description of the image.
  * `image_owner_alias` - Alias of the image owner.
  * `os_name` - Display Chinese name of the OS.
  * `os_name_en` - Display English name of the OS.
  * `os_type` - The operating system type of the image. Valid values: windows and linux.
  * `platform` - the Operating system platform of the image.
  * `status` - Status of the image. Possible values: `UnAvailable`, `Available`, `Creating` and `CreateFailed`.
  * `size` - Size of the image.
  * `disk_device_mappings` - Description of the system with disks and snapshots under the image.
    * `device` - Device information of the created disk: such as /dev/xvdb.
    * `size` - Size of the created disk.
    * `snapshot_id` - Snapshot ID.
  * `product_code` - Product code of the image on the image market.
  * `is_subscribed` - Whether the user has subscribed to the terms of service for the image product corresponding to the ProductCode.
  * `is_copied` - Is it a copied image. Available values: 'true' or 'false'。 
  * `is_self_shared` - Have you shared this custom image with other users. Available values: 'true' or 'false'。
  * `image_version` - Version of the image.
  * `usage` - Specifies whether to check the validity of the request without actually making the request. Valid values:
  * `progress` - Progress of image creation, presented in percentages.
  * `is_support_io_optimized` - Specifies whether the image can be used on I/O optimized instances.
  * `tags` - The tag of the resource.