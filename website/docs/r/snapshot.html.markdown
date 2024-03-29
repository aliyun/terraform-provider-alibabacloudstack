---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_snapshot"
sidebar_current: "docs-alibabacloudstack-resource-snapshot"
description: |-
  Provides an ECS snapshot resource.
---

# alibabacloudstack\_snapshot

Provides an ECS snapshot resource.

## Example Usage

```
data "alibabacloudstack_disks" "system_disk" {
	instance_id = "i-dw906gn8yjmrrzb46dmu"
	type 		= "system"
}

resource "alibabacloudstack_snapshot" "snapshot" {
  disk_id     = data.alibabacloudstack_disks.system_disk.ids[0]
  name        = "test-snapshot"
  description = "this snapshot is created for testing"
  tags = {
    version = "1.2"
  }
}

```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The ID of the disk.
* `name` - (Optional, ForceNew) Name of the snapshot. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://.
* `description` - (Optional, ForceNew) Description of the snapshot. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### Timeouts

* `create` - (Defaults to 2 mins) Used when creating the snapshot (until it reaches the initial `SnapshotCreatingAccomplished` status). 
* `delete` - (Defaults to 2 mins) Used when terminating the snapshot. 

## Attributes Reference

The following attributes are exported:

* `id` - The snapshot ID.
