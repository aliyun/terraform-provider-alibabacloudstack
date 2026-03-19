---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_dir_quota"
sidebar_current: "docs-Alibabacloudstack-nas-dir-quota"
description: |-
  Manage NAS directory quotas
---

# alibabacloudstack_nas_dir_quota

Manages directory quotas for Alibaba Cloud NAS file systems, allowing capacity and file count limits to be set for specified directories for users or user groups.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testaccnasdirquota60008"
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

resource "alibabacloudstack_nas_dir_quota" "default" {
  file_system_id = alibabacloudstack_nas_file_system.default.id
  path           = "/"
  quotas {
    file_count_limit = 10000
    quota_type       = "Enforcement"
    user_type        = "Uid"
    user_id          = "500"
    size_limit       = 100
  }
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system, used to specify the NAS file system for which the quota is set.
* `path` - (Required, ForceNew) The absolute path of the directory in the file system, such as "/" or "/data/sub1".
* `quotas` - (Optional) A collection of quota configurations. At least one quota item must be configured. Each quota item contains the following parameters:
  * `quota_type` - (Required) The type of quota. Valid values:
    * `Enforcement`: Enforcement quota. When usage exceeds the limit, operations such as creating files or directories and appending writes will fail.
    * `Accounting`: Accounting quota. Only usage is counted; operations are not restricted.
  * `user_type` - (Required) The type of user. Valid values:
    * `Uid`: User ID.
    * `Gid`: Group ID of the user.
    * `AllUsers`: All users.
  * `user_id` - (Optional) The Uid or Gid to be restricted. Required when `user_type` is `Uid` or `Gid`.
  * `size_limit` - (Optional) The total capacity limit for the user's files under the directory, in GB. When `quota_type` is `Enforcement`, at least one of `size_limit` or `file_count_limit` must be specified.
  * `file_count_limit` - (Optional) The file count limit for the user under the directory, including files, directories, and special files. When `quota_type` is `Enforcement`, at least one of `size_limit` or `file_count_limit` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `{FileSystemId:Path}`.
* `status` - The statistical status of the directory. Possible values:
  * `Initializing`: Initializing.
  * `Normal`: Normal status.