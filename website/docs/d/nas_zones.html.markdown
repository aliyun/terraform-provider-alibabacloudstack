---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_zones"
description: |-
    Provides a list of NAS available zones.
---

# alibabacloudstack_nas_zones

Provide a data source to retrieve the type of zone used to create NAS file system.


## Example Usage

```terraform
data "alibabacloudstack_nas_zones" "default" {}

output "alibabacloudstack_nas_zones_id" {
  value = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_type` - (Optional) The type of the file system. Valid values: `standard`, `extreme`, `cpfs`. Default value: `standard`.
* `zone_id` - (Optional) The zone ID to filter results.
* `protocol` - (Optional) The protocol type to filter results. Valid values: `NFS`, `SMB`, `cpfs`.
* `output_file` - (Optional, Deprecated) This field is deprecated and will be removed in version 3.19.0. To write content to a file, use the `local_file` provider instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zone information collection.
    * `zone_id` - The ID of the zone.
    * `protocols` - A list of protocol types supported in the zone.
    * `clusters` - A list of cluster information in the zone.
        * `cluster_id` - The ID of the cluster.
        * `cluster_type` - The type of the cluster.
        * `cluster_version` - The version of the cluster.
        * `instance_types` - A list of instance type information in the cluster.
            * `storage_type` - The storage type. Valid values vary by `file_system_type`:
                * When `file_system_type` is `standard`: `Performance`, `Capacity`.
                * When `file_system_type` is `extreme`: `Standard`, `Advance`.
                * When `file_system_type` is `cpfs`: `advance_100`, `advance_200`.
            * `protocol_type` - The file transfer protocol type. Valid values vary by `file_system_type`:
                * When `file_system_type` is `standard`: `NFS`, `SMB`.
                * When `file_system_type` is `extreme`: `NFS`.
                * When `file_system_type` is `cpfs`: `cpfs`.