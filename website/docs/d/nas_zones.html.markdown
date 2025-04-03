---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_zones"
sidebar_current: "docs-alibabacloudstack-datasource-nas-zones"
description: |-
    Provides a list of FileType owned by an Alibaba Cloud account.
---

# alibabacloudstack_nas_zones

Provide  a data source to retrieve the type of zone used to create NAS file system.


## Example Usage

```terraform
data "alibabacloudstack_nas_zones" "default" {}

output "alibabacloudstack_nas_zones_id" {
  value = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_type` - (Optional, ForceNew, Available in v1.152.0+) The type of the file system.  Valid values: `standard`, `extreme`, `cpfs`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zone information collection.
    * `zone_id` - String to filter results by zone id.
    * `instance_types` - A list of instance type information collection
        * `storage_type` - The storage type of the nas zones. Valid values:
          * `standard` - When FileSystemType is standard. Valid values: `Performance` and `Capacity`.
          * `extreme` - When FileSystemType is extreme. Valid values: `Standard` and `Advance`.
          * `cpfs` - When FileSystemType is cpfs. Valid values: `advance_100` and `advance_200` .
        * `protocol_type` - File transfer protocol type. Valid values:
          * `standard` - When FileSystemType is standard. Valid values: `NFS` and `SMB`.
          * `extreme` - When FileSystemType is extreme. Valid values: `NFS`.
          * `cpfs` - When FileSystemType is cpfs. Valid values: `cpfs`.