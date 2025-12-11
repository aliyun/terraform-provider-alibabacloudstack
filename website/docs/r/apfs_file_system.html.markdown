---
subcategory: "Storage"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apfs_file_system"
sidebar_current: "docs-Alibabacloudstack-storage-apfs_file_system"
description: |-
  Manage APFS file storage systems
---

# alibabacloudstack_apfs_file_system

Manage APFS file storage systems using credentials configured in the Provider within the specified resource set.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testacc47734"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_apfs_file_system" "default" {
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  cluster_id   = "EfsStorageCluster-A-20251125-0131"
  storage_type = "parastor_advance_300"
  volume_size  = "1024"
  description  = var.name
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, Forces new resource) The storage type. For example: `parastor_advance_300`. A valid storage type must be specified.
* `volume_size` - (Required) The capacity of the file system, in GB. Minimum value is 1, depending on the storage type.
* `cluster_id` - (Forces new resource) The cluster ID. Specifies the cluster to which the file system belongs, for example: `EfsStorageCluster-A-20251125-0131`.
* `file_system_type` - (Forces new resource) The file system type. Default value is `efs`, indicating a standard file system.
* `protocol_type` - (Forces new resource) The protocol type. Default value is `EFS`, indicating the use of EFS protocol.
* `zone_id` - (Forces new resource) The zone ID. Specifies the zone where the file system is created, for example: `cn-wulan-env17e-amtest17001-a`.
* `description` - (Optional) The description of the file system. Length must be between 1 and 256 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the file system, in the format like `efs-0187aa76`.
* `create_time` - The creation time of the file system, in ISO8601 format (for example: `2025-12-11T07:12:37Z`).
* `file_system_id` - The ID of the file system (same as the `id` attribute).
* `metered_size` - The metered storage size, in GB. Represents the current actual storage usage.
* `quota_size` - The quota size of the file system, in GB. Represents the maximum storage quota allocated by the system.
* `status` - The current status of the file system. Possible values include `Pending` (creating), `Running` (running), `Deleting` (deleting), etc.