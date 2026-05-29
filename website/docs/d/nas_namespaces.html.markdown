---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespaces"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-namespaces"
description: |-
  Query Alibaba Cloud NAS namespace information
---

# alibabacloudstack_nas_namespaces

The NAS unified namespace list data source is used to query Alibaba Cloud NAS namespace information.

## Example Usage

```hcl
variable "name" {
  default = "tf-testnasfs3606"
}

data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

data "alibabacloudstack_nas_namespaces" "default" {
  ids = [
    "${alibabacloudstack_nas_namespace.default.id}"
  ]
}
```

## Argument Reference

The following arguments are supported as filter criteria:

* `file_system_type` - (String, Optional) The type of the file system. Default value: `standard`.
* `ids` - (List, Optional) A list of namespace IDs used to filter results.
* `name_regex` - (String, Optional) A regular expression used to filter results by namespace description.
* `protocol_type` - (String, Optional) The protocol type. Possible values: `NFS`, `SMB`.
* `storage_type` - (String, Optional) The storage type. Possible values: `Performance`, `Capacity`.
* `zone_id` - (String, Optional) The ID of the zone to which the namespace belongs.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of matched namespace IDs.
* `namespaces` - A list of matched namespaces. Each element contains the following attributes:
  * `id` - The namespace ID, equivalent to `nas_namespace_id`.
  * `create_time` - The creation time of the namespace, following the ISO 8601 standard in the format `yyyy-MM-ddTHH:mm:ssZ`.
  * `description` - The description of the namespace.
  * `encrypt_type` - The encryption type of the namespace. Valid values: `0` (no encryption), `1` (encryption).
  * `file_system_type` - The type of the file system. Default value: `standard`, indicating General-purpose NAS.
  * `mount_target_count` - The number of mount targets.
  * `nas_namespace_id` - The namespace ID.
  * `protocol_type` - The protocol type. Possible values: `NFS` (NFS file transfer protocol), `SMB` (SMB file transfer protocol).
  * `storage_type` - The storage type. Possible values: `Performance` (Performance type), `Capacity` (Capacity type).
  * `status` - The status of the namespace. Possible values: `Initializing` (initializing), `Normal` (normal).
  * `zone_id` - The ID of the zone.