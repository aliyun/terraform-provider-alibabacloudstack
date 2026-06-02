---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_filesystem_attachment"
description: |-
  Manages the attachment of NAS file systems to unified namespaces.
---

# alibabacloudstack_nas_namespace_filesystem_attachment

Creates, manages, and configures the attachment of NAS file systems to unified namespaces.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testAccNasnpsFsAttachment26768"
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
  mapped_path      = var.name
  file_system_id   = alibabacloudstack_nas_file_system.default.id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) File system ID. The identifier of the file system to be mapped to the namespace.
* `nas_namespace_id` - (Required, ForceNew) Namespace ID. The identifier of the unified namespace.
* `mapped_path` - (Required) Mapped path. The mount path of the file system within the namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `{NasNamespaceId:FileSystemId}`.
* `create_time` - The creation time of the attachment, in ISO 8601 standard format.
* `file_system_type` - The type of the file system. Possible values:
  * `standard`: Standard NAS
* `storage_type` - The storage type. Possible values:
  * `Capacity`: Capacity type
  * `Performance`: Performance type