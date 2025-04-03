---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshots"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-snapshots"
description: |- 
  Provides a list of ecs snapshots owned by an alibabacloudstack account according to the specified filters.
---

# alibabacloudstack_ecs_snapshots
-> **NOTE:** Alias name has: `alibabacloudstack_snapshots`

This data source provides a list of ECS snapshots in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_snapshots" "example" {
  instance_id   = "i-abc123456"
  disk_id       = "d-abc123456"
  name_regex    = "my-snapshot-*"
  status        = "accomplished"
  type          = "user"
  source_disk_type = "system"
  usage         = "image"
}

output "snapshot_ids" {
  value = data.alibabacloudstack_ecs_snapshots.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew) The ID of the instance associated with the snapshot.
* `disk_id` - (Optional, ForceNew) The ID of the disk associated with the snapshot.
* `ids` - (Optional, ForceNew) A list of snapshot IDs. If specified, the data source will return snapshots that match these IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by snapshot name.
* `status` - (Optional, ForceNew) The status of the snapshot. Valid values:
  * `progressing`: The snapshot is being created.
  * `accomplished`: The snapshot is created and ready to use.
  * `failed`: The snapshot creation failed.
  * `all` (default): This value indicates all snapshot states.
* `type` - (Optional, ForceNew) The category of the snapshot. Valid values:
  * `auto`: Auto snapshots.
  * `user`: Manual snapshots.
  * `all` (default): Both auto and manual snapshots.
* `source_disk_type` - (Optional, ForceNew) The type of source disk for the snapshot. Valid values:
  * `system`: System disk.
  * `data`: Data disk.
  > Note: The value of this parameter is case-insensitive.
* `usage` - (Optional, ForceNew) Specifies whether the snapshot has been used to create custom images or disks. Valid values:
  * `image`: The snapshot has been used to create custom images.
  * `disk`: The snapshot has been used to create disks.
  * `image_disk`: The snapshot has been used to create both custom images and data disks.
  * `none`: The snapshot has not been used to create custom images or disks.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of snapshot IDs.
* `names` - A list of snapshot names.
* `snapshots` - A list of snapshots. Each element contains the following attributes:
  * `id` - The ID of the snapshot.
  * `name` - The name of the snapshot.
  * `description` - The description of the snapshot. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
  * `progress` - Snapshot creation progress, presented in percentage.
  * `source_disk_id` - The ID of the source disk.
  * `source_disk_size` - The size of the source disk, measured in GB.
  * `source_disk_type` - The type of the source disk. Valid values:
    * `system`: System disk.
    * `data`: Data disk.
  * `product_code` - The product code inherited from the image market place.
  * `remain_time` - The remaining time for the snapshot creation task, in seconds.
  * `creation_time` - The creation time of the snapshot. It follows the ISO8601 standard and uses UTC time. Format: `YYYY-MM-DDThh:mmZ`.
  * `status` - The status of the snapshot. Valid values:
    * `progressing`: The snapshot is being created.
    * `accomplished`: The snapshot is created and ready to use.
    * `failed`: The snapshot creation failed.
    * `all`: This value indicates all snapshot states.
  * `usage` - Whether the snapshot has been used to create resources or not. Valid values:
    * `image`: The snapshot has been used to create custom images.
    * `disk`: The snapshot has been used to create disks.
    * `image_disk`: The snapshot has been used to create both custom images and data disks.
    * `none`: The snapshot has not been used to create custom images or disks.
* `tags` - A map of tags assigned to the snapshot.