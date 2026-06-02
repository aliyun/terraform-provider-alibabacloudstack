---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_filesystem_attachment"
description: |-
  Provides information about Alibaba Cloud NAS namespace filesystem attachments.
---

# alibabacloudstack_nas_namespace_filesystem_attachment

> This data source provides information about NAS namespace filesystem attachments, used to query mapped filesystem information in a namespace.

## Example Usage

```hcl
variable "name" {
  default = "tf-testnasfs1942"
}

data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_file_system" "default" {
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  encrypt_type  = "0"
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
}

resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

resource "alibabacloudstack_nas_namespace_filesystem_attachment" "default" {
  nas_namespace_id = alibabacloudstack_nas_namespace.default.id
  file_system_id   = alibabacloudstack_nas_file_system.default.id
  mapped_path      = var.name
}

data "alibabacloudstack_nas_namespace_filesystem_attachments" "default" {
  nas_namespace_id = alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id
  ids = [
    "${alibabacloudstack_nas_namespace_filesystem_attachment.default.id}"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `nas_namespace_id` - (String, Required) The ID of the namespace, used to specify the namespace to query.
* `file_system_id` - (String, Optional) The ID of the filesystem, used to filter mapping information for a specific filesystem.
* `ids` - (List, Optional) A list of attachment IDs in the format NasNamespaceId:FileSystemId, used to precisely match specific attachments.
* `mapped_path` - (String, Optional) The mapped path, used to filter filesystems with a specific mapped path.
* `name_regex` - (String, Optional) A regular expression used to filter results by mapped path.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The data source ID, generated from the hash of the attachment ID.
* `attachment_id` - (String) The ID of the attachment, in the format NasNamespaceId:FileSystemId.
* `create_time` - (String) The creation time of the attachment, in ISO 8601 standard format.
* `file_system_type` - (String) The type of the filesystem, for example, "standard" represents General-purpose NAS.
* `storage_type` - (String) The storage type of the filesystem, possible values include "Capacity" (Capacity-type) and "Performance" (Performance-type).
* `id` - (String, Deprecated) The field 'id' has been deprecated since provider version 1.200.0, please use 'attachment_id' instead.