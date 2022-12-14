---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_snapshots"
sidebar_current: "docs-alibabacloudstack-datasource-snapshots"
description: |-
  Provides a data source to get a list of snapshot according to the specified filters.
---

# alibabacloudstack\_snapshots

Use this data source to get a list of snapshot according to the specified filters in an Alibabacloudstack Cloud account.

## Example Usage

```
data "alibabacloudstack_snapshots" "snapshots" {
  name_regex = "${var.name_regex}"
}

output "snapshots" {
  value = data.alibabacloudstack_snapshots.snapshots
}
```

##  Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) The specified instance ID.
* `disk_id` - (Optional) The specified disk ID.
* `ids` - (Optional)  A list of snapshot IDs.
* `name_regex` - (Optional) A regex string to filter results by snapshot name.
* `status` - (Optional) The specified snapshot status.
  * The snapshot status. Optional values:
  * progressing: The snapshots are being created.
  * accomplished: The snapshots are ready to use.
  * failed: The snapshot creation failed.
  * all: All status.
  
  Default value: all.

* `type` - (Optional) The snapshot category. Optional values:
  * auto: Auto snapshots.
  * user: Manual snapshots.
  * all: Auto and manual snapshots.
  
  Default value: all.
* `source_disk_type` - (Optional) The type of source disk:
  * System: The snapshots are created for system disks.
  * Data: The snapshots are created for data disks.
  
* `usage` - (Optional) The usage of the snapshot:
  * image: The snapshots are used to create custom images.
  * disk: The snapshots are used to CreateDisk.
  * mage_disk: The snapshots are used to create custom images and data disks.
  * none: The snapshots are not used yet.
* `tags` - (Optional) A map of tags assigned to snapshots.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of snapshot IDs.
* `names` - A list of snapshots names.
* `snapshots` - A list of snapshots. Each element contains the following attributes:
    * `id` - ID of the snapshot.
    * `name` - Name of the snapshot.
    * `description` - Description of the snapshot.
    * `progress` - Progress of snapshot creation, presented in percentage.
    * `source_disk_id` - Source disk ID, which is retained after the source disk of the snapshot is deleted.
    * `source_disk_size` - Size of the source disk, measured in GB.
    * `source_disk_type` - Source disk attribute. Value range:
      * System
      * Data
    * `product_code` - Product code on the image market place.
    * `remain_time` - The remaining time of a snapshot creation task, in seconds.
    * `creation_time` - Creation time. Time of creation. It is represented according to ISO8601, and UTC time is used. Format: YYYY-MM-DDThh:mmZ.
    * `status` - The snapshot status. Value range:
      * progressing
      * accomplished
      * failed
    * `usage` - Whether the snapshots are used to create resources or not. Value range:
      * image
      * disk
      * image_disk
      * none
* `tags` - A map of tags assigned to the snapshot.
