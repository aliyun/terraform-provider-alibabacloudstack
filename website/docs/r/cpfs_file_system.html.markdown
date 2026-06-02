---
subcategory: "File Storage CPFS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cpfs_file_system"
description: |-
  Provides a CPFS File System resource.
---

# alibabacloudstack_cpfs_file_system

Provides a CPFS File System resource.

> **Note:** This resource can also be referred to by the following alias:
> - `alibabacloudstack_cpfs_file_system`

## Example Usage

Basic Usage

```terraform
variable "name" {
	default = "tfacctest124"
}

data "alibabacloudstack_nas_zones" "default" {
	file_system_type = "bmcpfs"
}

resource "alibabacloudstack_cpfs_file_system" "example" {
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = "${var.name}"
	capacity = 20480
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, ForceNew) The storage type for the CPFS file system. Valid values: `advance_100` (100 MB/s/TiB baseline), `advance_200` (200 MB/s/TiB baseline), `economic` (economy).
* `protocol_type` - (Optional, ForceNew) File transfer protocol type. Default value: `CPFS`. Valid value: `CPFS`.
* `description` - (Optional) File system description. It must be 2 to 256 characters in length.
* `encrypt_type` - (Optional, ForceNew, Deprecated) **This field is deprecated and will be removed in version 3.21.0.** Whether the file system is encrypted. Valid values: `0` (default, not encrypted), `1` (encrypted), `2` (encrypted with KMS key, requires `kms_key_id`).
* `file_system_type` - (Optional, ForceNew) File system type. Default value: `bmcpfs`. Valid value: `bmcpfs`.
* `capacity` - (Optional) The capacity of the file system. Unit: GiB. Valid values: 20480 to 1740800. This parameter is required when `file_system_type` is `bmcpfs`.
* `zone_id` - (Required, ForceNew) The zone ID.
* `cluster_id` - (Required, ForceNew) The cluster ID. This parameter is specific to ApsaraStack and maps to the `Location` field returned by the API.
* `kms_key_id` - (Optional, Deprecated) **This field is deprecated and will be removed in version 3.21.0.** The ID of the KMS key. Required when `encrypt_type` is set to `2`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CPFS file system.
* `capacity` - (Computed) The capacity of the file system.
* `zone_id` - (Computed) The zone ID.
* `kms_key_id` - (Computed, Deprecated) The ID of the KMS key.
* `storage_type` - (Computed) The storage type.
* `protocol_type` - (Computed) The protocol type.
* `file_system_type` - (Computed) The file system type.
* `encrypt_type` - (Computed, Deprecated) The encryption type.
* `description` - (Computed) The file system description.

## Import

CPFS File System can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cpfs_file_system.example 1337849c59
```
