---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_lifecyclepolicies"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-lifecyclepolicies"
description: |-
  Provides a list of nas lifecyclepolicies owned by an alibabacloudstack account.
---

# alibabacloudstack\_nas\_lifecyclepolicies

This data source provides a list of nas lifecyclepolicies in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testacc-nas-liecycle-datasource315726"
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
	lifecycle_policy_name = "${var.name}"
	file_system_id        = "${alibabacloudstack_nas_file_system.default.id}"
	path                  = "/"
	recursive             = "false"
	lifecycle_rule_name   = "DEFAULT_ATIME_14"
	oss_bucket            = "${alibabacloudstack_oss_bucket.default.id}"

}
data "alibabacloudstack_nas_lifecycle_policies" "default" {
	depends_on = [
		alibabacloudstack_nas_lifecycle_policy.default
	]
	file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of nas lifecyclepolicy IDs.
  * `file_system_id` - (Optional) - The ID of the file system.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `lifecycle_policies` - The list of nas lifecyclepolicies.
    * `id` - The ID of the nas lifecyclepolicy.
    * `create_time` - The creation time of the resource
    * `file_system_id` - The ID of the file system.
    * `lifecycle_policy_name` - The first ID of the resource
    * `lifecycle_rule_name` - Management rules associated with lifecycle management policies.Value:-DEFAULT_ATIME_14: files not accessed 14 days ago-DEFAULT_ATIME_30: files not accessed 30 days ago-DEFAULT_ATIME_60: files not accessed 60 days ago-DEFAULT_ATIME_90: files not accessed 90 days ago
    * `path` - Absolute path to the directory associated with the lifecycle management policy.Only associating a single directory is supported. It must start with a forward slash (/) and be the real path in the Mount point.> It is recommended that you configure Paths.N to associate multiple directories at the same time.
    * `recursive` - Whether to recursively apply sub-paths
    * `storage_type` - The type of storage after the data dump.Default value: InfrequentAccess (low frequency media storage)
    * `oss_bucket` - OSS bucket name
