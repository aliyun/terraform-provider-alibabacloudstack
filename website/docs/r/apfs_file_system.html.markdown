---
subcategory: "APFS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apfs_file_system"
sidebar_current: "docs-Alibabacloudstack-resource-apfs-file-system"
description: |-
  Manage APFS file storage systems
---

# alibabacloudstack_apfs_file_system

Manage APFS (Alibaba Parallel File System) file storage systems using credentials configured in the Provider within the specified resource set.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testacc47734"
}

data "alibabacloudstack_apfs_zones" "default" {
}

resource "alibabacloudstack_apfs_file_system" "default" {
  zone_id      = data.alibabacloudstack_apfs_zones.default.zones.0.zone_id
  cluster_id   = data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.cluster_id
  storage_type = data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.storage_type
  volume_size  = 256
  description  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, Forces new resource) The storage type. A valid storage type must be specified, for example: `parastor_advance_300`.
* `volume_size` - (Required) The capacity of the file system, in GB.
* `cluster_id` - (Optional, Forces new resource) The cluster ID. Specifies the cluster to which the file system belongs.
* `file_system_type` - (Optional, Forces new resource) The file system type. Default value is `efs`, indicating a standard file system.
* `protocol_type` - (Optional, Forces new resource) The protocol type. Default value is `EFS`, indicating the use of EFS protocol.
* `zone_id` - (Optional, Forces new resource) The zone ID. Specifies the zone where the file system is created.
* `description` - (Optional) The description of the file system. Length must be between 0 and 256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the file system, for example: `efs-0187aa76`.
* `file_system_id` - The ID of the file system (same as the `id` attribute).
* `create_time` - The creation time of the file system, in ISO8601 format (for example: `2025-12-11T07:12:37Z`).
* `status` - The current status of the file system. Possible values include `Pending` (creating), `Running` (running), `Deleting` (deleting), etc.
* `metered_size` - The metered storage size, in GB. Represents the current actual storage usage.
* `quota_size` - The quota size of the file system, in GB. Represents the maximum storage quota allocated by the system.

## Import

APFS File System can be imported using the FileSystemId, e.g.

```
$ terraform import alibabacloudstack_apfs_file_system.example efs-xxxxxxxx
```