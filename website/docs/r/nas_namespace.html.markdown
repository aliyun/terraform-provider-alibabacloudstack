---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace"
description: |-
  Create and manage NAS unified namespace
---

# alibabacloudstack_nas_namespace

Create and manage Alibaba Cloud NAS unified namespace resources.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-testAccNasNamespace62833"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}



data "alibabacloudstack_nas_zones" "default" {
}



resource "alibabacloudstack_nas_namespace" "default" {
  encrypt_type  = "0"
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = "tf-testAccNasNamespace62833"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, Forces new resource) The ID of the zone. A zone is a physical area within the same region where power and network are independent of each other. File systems and ECS instances in different zones within the same region are interconnected. It is recommended that the file system and ECS instances belong to the same zone to avoid latency caused by cross-zone operations.
* `cluster_id` - (Required, Forces new resource) The ID of the cluster. Default value: "StandardNasCluster".
* `description` - (Required, Forces new resource) The description of the namespace. It must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), and hyphens (-). It must start with a letter or Chinese character and cannot start with `http://` or `https://`.
* `storage_type` - (Required, Forces new resource) The storage type. Valid values: `Capacity` (Capacity type) or `Performance` (Performance type).
* `protocol_type` - (Required, Forces new resource) The file transfer protocol type. Valid values: `NFS` (NFS file transfer protocol) or `SMB` (SMB file transfer protocol).
* `encrypt_type` - (Required, Forces new resource) Specifies whether the file system in the namespace is encrypted. Uses KMS service to manage keys for encrypting file system data at rest. No decryption is required when reading and writing encrypted data. Valid values: `0` (not encrypted) or `1` (encrypted).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the namespace.
* `create_time` - The time when the namespace was created. Follows the ISO 8601 standard, in the format: `yyyy-MM-ddTHH:mm:ssZ`.
* `file_system_type` - The type of the file system. Default value: `standard`, which indicates General-purpose NAS.
* `mount_target_count` - The number of mount targets for the namespace file system.