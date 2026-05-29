---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_images"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-images"
description: |-
  Provides a list of BMCP (Bare Metal Compute Platform) images owned by an alibabacloudstack account.
---

# alibabacloudstack_bmcp_images

This data source provides a list of BMCP images in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_bmcp_images" "default" {
  type = "public"
}
```

### Filter by name regex

```hcl
data "alibabacloudstack_bmcp_images" "ubuntu_images" {
  type       = "public"
  name_regex = "ubuntu"
}
```

### Filter by ids

```hcl
data "alibabacloudstack_bmcp_images" "filtered_images" {
  type = "public"
  ids  = ["i-d"]
}
```

### Combine filters

```hcl
data "alibabacloudstack_bmcp_images" "combined" {
  type       = "public"
  name_regex = "ubuntu"
  ids        = ["i-d"]
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of the image. Valid values: `customize`, `public`.
* `name_regex` - (Optional) A regex string to filter resulting images by name.
* `ids` - (Optional) A list of image IDs for filtering. Supports fuzzy matching.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `images` - A list of images. Each element contains the following attributes:
  * `id` - The ID of the image (same as `name`).
  * `name` - The name of the image.
  * `description` - The description of the image.
  * `type` - The type of the image.
  * `platform` - The platform of the image.
  * `zone` - The zone of the image.
  * `region` - The region of the image.
  * `status` - The status of the image.
  * `size` - The size of the image.
  * `create_time` - The creation time of the image.
  * `update_time` - The update time of the image.
  * `os_arch` - The OS architecture of the image.
  * `disk_format` - The disk format of the image.
  * `package_type` - The package type of the image.
  * `url` - The URL of the image.
  * `unique_key` - The unique key of the image.
  * `template` - The template of the image.
  * `resource_group` - The resource group of the image.
  * `file_name` - The file name of the image.
  * `snapshot_id` - The snapshot ID of the image.
  * `error_description` - The error description of the image.
  * `part` - The part of the image.
  * `reboot_script` - The reboot script of the image.
  * `signature_file_url` - The signature file URL of the image.
  * `owner` - The owner of the image.
  * `component_type` - The component type of the image.
  * `version` - The version of the image.
  * `custom_script` - The custom script of the image.
  * `deleted` - Whether the image is deleted.
  * `disable` - Whether the image is disabled.
  * `kernel_version` - The kernel version of the image.
  * `organization` - The organization of the image.
  * `check_sum` - The checksum of the image.
