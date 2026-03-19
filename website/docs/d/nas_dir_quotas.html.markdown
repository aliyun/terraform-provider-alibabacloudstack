---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_dir_quota"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-dir-quota"
description: |-
  Query NAS directory quota information
---

# alibabacloudstack_nas_dir_quota

This data source retrieves directory quota information for Alibaba Cloud NAS file systems. It is used to query directory information with configured quotas in a specified file system, including directory paths, inode numbers, and quota limits and actual usage for each user.

## Example Usage

```hcl
variable "name" {
  default = "tf-testnasdirquotas2499023223732943598"
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
    quota_type       = "Enforcement"
    user_type        = "Uid"
    user_id          = "500"
    size_limit       = 100
    file_count_limit = 10000
  }
}

data "alibabacloudstack_nas_dir_quotas" "default" {
  file_system_id = alibabacloudstack_nas_dir_quota.default.file_system_id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` (Required): The ID of the file system, used to specify the NAS file system to query.
* `path` (Optional): The absolute path of the directory in the file system. If not specified, all directories with configured quotas in the file system will be returned.
* `ids` (Optional): Filter results by a list of IDs. The ID format is {FileSystemId:Path}.
* `name_regex` (Optional): Filter results by regular expression, matching the directory path.

## Attributes Reference

The following attributes are exported:

* `id` (String): The resource ID, in the format {FileSystemId:Path}.

The following attributes are exported in the `dir_quotas` list:

* `dir_inode` (String): The inode number of the directory.
* `file_system_id` (String): The ID of the file system.
* `path` (String): The absolute path of the directory in the file system.
* `user_quotas` (List): A list of user quota information.

The following attributes are exported in the `user_quotas` list:

* `file_count_limit` (Integer): The file count limit for the user in the directory. -1 indicates no limit.
* `file_count_real` (Integer): The actual file count for the user in the directory.
* `quota_type` (String): The quota type, including Accounting (statistical) and Enforcement (restrictive).
* `size_limit` (Integer): The total file capacity limit for the user in the directory, in GB. -1 indicates no limit.
* `size_real` (Integer): The actual total file capacity for the user in the directory, in GB.
* `user_id` (String): The uid or gid to be restricted, depending on the value of UserType. This value is empty when UserType is AllUsers.
* `user_type` (String): The type of UserId, including Uid, Gid, and AllUsers.