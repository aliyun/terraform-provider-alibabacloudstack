---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_disks"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-disks"
description: |- 
  Provides a list of ecs disks owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_disks
-> **NOTE:** Alias name has: `alibabacloudstack_disks`

This data source provides a list of ECS disks in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_disks" "disks_ds" {
  name_regex = "sample_disk"
  type       = "data"
  category   = "cloud_ssd"
  instance_id = "i-bp1234567890abcdefg"

  tags = {
    Environment = "Production"
    Owner      = "TeamA"
  }
}

output "first_disk_id" {
  value = "${data.alibabacloudstack_ecs_disks.disks_ds.disks.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of disk IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by disk name.
* `type` - (Optional, ForceNew) The method that you want to use to resize the disk. Valid values:
  * `offline` (default): Resizes the disk offline. After resizing the disk offline, you must restart the associated instance for the resizing operation to take effect.
  * `online`: Resizes the disk online. After resizing the disk online, the resizing operation immediately takes effect without restarting the associated instance.
* `category` - (Optional, ForceNew) The category of the disk. Valid values:
  * `all`: All disk categories.
  * `cloud`: Basic cloud disk.
  * `cloud_efficiency`: Ultra cloud disk.
  * `cloud_ssd`: Standard SSD cloud disk.
  * `cloud_essd`: Enterprise SSD (ESSD) cloud disk.
  * `cloud_auto`: ESSD AutoPL disk.
  * `local_ssd_pro`: I/O-intensive local disk.
  * `local_hdd_pro`: Throughput-intensive local disk.
  * `cloud_essd_entry`: ESSD Entry disk.
  * `elastic_ephemeral_disk_standard`: Standard elastic ephemeral disk.
  * `elastic_ephemeral_disk_premium`: Premium elastic ephemeral disk.
  * `ephemeral`: Retired local disk.
  * `ephemeral_ssd`: Retired local SSD.
  Default value: `all`.
* `instance_id` - (Optional, ForceNew) The ID of the subscription instance to which to attach the subscription disk.
  * If you specify an instance ID, the following parameters are ignored: `ResourceGroupId`, `Tag.N.Key`, `Tag.N.Value`, `ClientToken`, and `KMSKeyId`.
  * You cannot specify both `ZoneId` and `InstanceId` in a request.
  This parameter is empty by default, which indicates that a pay-as-you-go disk is created in the region and zone specified by `RegionId` and `ZoneId`.
* `tags` - (Optional) A map of tags assigned to the disks.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `disks` - A list of disks. Each element contains the following attributes:
  * `id` - The ID of the disk.
  * `name` - The name of the disk.
  * `description` - The description of the disk. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`. This parameter is left empty by default.
  * `region_id` - The ID of the region in which the disk resides.
  * `availability_zone` - The availability zone of the disk.
  * `status` - The lifecycle status of the disk. Possible values:
    * `In_use`: The disk is in use.
    * `Available`: The disk can be attached.
    * `Attaching`: The disk is being attached.
    * `Detaching`: The disk is being detached.
    * `Creating`: The disk is being created.
    * `ReIniting`: The disk is being initialized.
  * `type` - The method that you want to use to resize the disk. Valid values:
    * `offline` (default): Resizes the disk offline. After resizing the disk offline, you must restart the associated instance for the resizing operation to take effect.
    * `online`: Resizes the disk online. After resizing the disk online, the resizing operation immediately takes effect without restarting the associated instance.
  * `category` - The category of the disk. See the `category` argument for valid values.
  * `size` - The size of the disk in GiB.
  * `image_id` - The ID of the image from which the disk is created. It is null unless the disk is created using an image.
  * `snapshot_id` - The snapshot used to create the disk. It is null if no snapshot is used to create the disk.
  * `instance_id` - The ID of the related instance. It is `null` unless the `status` is `In_use`.
  * `kms_key_id` - The ID of the KMS key corresponding to the data disk.
  * `creation_time` - The creation time of the disk.
  * `attached_time` - The attachment time of the disk.
  * `detached_time` - The detachment time of the disk.
  * `storage_set_id` - The ID of the storage set to which the disk belongs.
  * `expiration_time` - The expiration time of the disk.
  * `tags` - A map of tags assigned to the disk.