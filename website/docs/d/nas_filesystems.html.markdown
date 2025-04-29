---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_filesystems"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-filesystems"
description: |- 
  Provides a list of nas filesystems owned by an alibabacloudstack account.
---

# alibabacloudstack_nas_filesystems
-> **NOTE:** Alias name has: `alibabacloudstack_nas_file_systems`

This data source provides a list of NAS file systems in an AlibabacloudStack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_nas_filesystems" "example" {
  storage_type   = "Performance"
  protocol_type  = "NFS"
  description_regex = "example-file-system"
}

output "file_system_id" {
  value = "${data.alibabacloudstack_nas_filesystems.example.systems.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Optional, ForceNew) The storage type:
  * When `file_system_type = standard`, the values are `Performance`, `Capacity`, and `Premium`.
  * When `file_system_type = extreme`, the values are `standard` or `advance`.
  * When `file_system_type = cpfs`, the values are `advance_100` (100MB/s/TiB baseline) and `advance_200` (200MB/s/TiB baseline).
  
* `protocol_type` - (Optional, ForceNew) File transfer protocol type:
  * When `file_system_type = standard`, the values are `NFS` and `SMB`.
  * When `file_system_type = extreme`, the value is `NFS`.
  * When `file_system_type = cpfs`, the value is `cpfs`.

* `description_regex` - (Optional, ForceNew) A regex string used to filter results by file system description.

* `ids` - (Optional) A list of file system IDs.


* `file_system_type` - (Optional) File system type:
  * `standard` (default): Universal NAS.
  * `extreme`: Extreme NAS.
  * `cpfs`: File Storage CPFS.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `descriptions` - A list of file system descriptions.

* `systems` - A list of file systems. Each element contains the following attributes:
  * `id` - ID of the file system.
  * `region_id` - Region ID where the file system is located.
  * `create_time` - Time of creation.
  * `description` - Description of the file system. Restrictions:
    * Length between 2~128 English or Chinese characters.
    * Must start with upper and lower case letters or Chinese, and cannot start with 'http://' or 'https://'.
    * Can contain numbers, colons (:), underscores (_), or dashes (-).
  * `protocol_type` - File transfer protocol type:
    * When `file_system_type = standard`, the values are `NFS` and `SMB`.
    * When `file_system_type = extreme`, the value is `NFS`.
    * When `file_system_type = cpfs`, the value is `cpfs`.
  * `storage_type` - The storage type:
    * When `file_system_type = standard`, the values are `Performance`, `Capacity`, and `Premium`.
    * When `file_system_type = extreme`, the value is `standard` or `advance`.
    * When `file_system_type = cpfs`, the values are `advance_100` (100MB/s/TiB baseline) and `advance_200` (200MB/s/TiB baseline).
  * `metered_size` - Metered size of the file system.
  * `encrypt_type` - Whether the file system is encrypted:
    * `0` (default): Not encrypted.
    * `1`: Encrypted with NAS managed key. Supported when `file_system_type = standard` or `extreme`.
    * `2`: Encrypted with user-managed key. Supported when `file_system_type = extreme`.
  * `file_system_type` - File system type:
    * `standard` (default): Universal NAS.
    * `extreme`: Extreme NAS.
    * `cpfs`: File Storage CPFS.
  * `capacity` - Capacity of the file system.
  * `kms_key_id` - The ID of the KMS key.
  * `zone_id` - The zone ID. The usable area refers to the physical area where power and network are independent of each other in the same region:
    * Optional when `file_system_type = standard`. By default, a zone that meets the conditions is randomly selected based on the `protocol_type` and `storage_type` configurations.
    * Required when `file_system_type = extreme` or `file_system_type = cpfs`. We recommend that the file system and the ECS instance belong to the same zone to avoid cross-zone latency.