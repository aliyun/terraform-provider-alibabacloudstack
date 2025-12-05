---
subcategory: "File Storage CPFS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cpfs_filesystem"
sidebar_current: "docs-alibabacloudstack-resource-cpfs-filesystem"
description: |-
  Provides a CPFS File System resource.
---

# alibabacloudstack_cpfs_filesystem

Provides a CPFS File System resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
	default = "tfacctest124"
}

data "alibabacloudstack_nas_zones" "default" {
	file_system_type = "bmcpfs"
}

resource "alibabacloudstack_cpfs_filesystem" "example" {
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = "${var.name}"
	capacity = 20480
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, ForceNew) The storage type for the CPFS file system.
* `protocol_type` - (Optional, ForceNew) File transfer protocol type. Default value: `CPFS`. Valid value: `CPFS`.
* `description` - (Optional) File system description. It must be 2 to 256 characters in length.
* `encrypt_type` - (Optional, ForceNew) Whether the file system is encrypted. Valid values:
  * `0` (default): Not encrypted.
  * `1`: Encrypted.
* `file_system_type` - (Optional, ForceNew) File system type. Default value: `bmcpfs`. Valid value: `bmcpfs`.
* `capacity` - (Optional) The capacity of the file system. Unit: GiB. Valid values: 20480 to 1740800.
* `zone_id` - (Required, ForceNew) The zone ID.
* `cluster_id` - (Required, ForceNew) The cluster ID.
* `kms_key_id` - (Optional) The ID of the KMS key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CPFS file system.
* `capacity` - The capacity of the file system.
* `zone_id` - The zone ID.
* `kms_key_id` - The ID of the KMS key.

## Import

CPFS File System can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cpfs_filesystem.example 1337849c59
```
