---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_lifecycle_policy"
description: |-
  Provides a nas Lifecyclepolicy resource.
---

# alibabacloudstack\_nas\_lifecyclepolicy

Provides a nas Lifecyclepolicy resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccnaslifecycle32504"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_file_system" "default" {
  protocol_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}"
  storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
  encrypt_type = "0"
  zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
  cluster_id ="${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
  description = "${var.name}"
}


resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "${var.name}"
  acl    = "public-read"
}



resource "alibabacloudstack_nas_lifecycle_policy" "default" {
  file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
  path = "/"
  recursive = "false"
  lifecycle_rule_name = "DEFAULT_ATIME_14"
  oss_bucket = "${alibabacloudstack_oss_bucket.default.id}"
  lifecycle_policy_name = "tf-testaccnaslifecycle32504"
}
```

## Argument Reference

The following arguments are supported:
  * `lifecycle_policy_name` - (Required, ForceNew) - The first ID of the resource
  * `storage_type` - (Optional, ForceNew) - The type of storage after the data dump.Default value: InfrequentAccess (low frequency media storage)
  * `file_system_id` - (Required, ForceNew) - The ID of the file system.
  * `path` - (Required, ForceNew) - Absolute path to the directory associated with the lifecycle management policy.Only associating a single directory is supported. It must start with a forward slash (/) and be the real path in the Mount point.> It is recommended that you configure Paths.N to associate multiple directories at the same time.
  * `recursive` - (Optional) - Whether to recursively apply sub-paths.
  * `lifecycle_rule_name` - (Required) - Management rules associated with lifecycle management policies.Value:-DEFAULT_ATIME_14: files not accessed 14 days ago-DEFAULT_ATIME_30: files not accessed 30 days ago-DEFAULT_ATIME_60: files not accessed 60 days ago-DEFAULT_ATIME_90: files not accessed 90 days ago
  * `oss_bucket` - (Required, ForceNew) - OSS bucket name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `create_time` - The creation time of the resource
