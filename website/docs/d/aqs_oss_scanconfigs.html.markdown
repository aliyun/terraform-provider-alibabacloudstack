---
subcategory: "Server Guard"
layout: "alibabacloudstack"
page_title: "Alibaba Cloud Stack: alibabacloudstack_aqs_oss_scanconfigs"
description: |-
  Provides a list of AQS OSS Scan Configs available to Alibaba Cloud Stack users.
---

# alibabacloudstack_aqs_oss_scanconfigs

This data source provides a list of AQS OSS Scan Configs in an Alibaba Cloud Stack account according to the specified filters.

> Note: Orchestrating this resource requires installing a specific version (Version 3.18.6-6.6.0.2 and later versions) of AQS.

## Example Usage

```hcl
variable "name" {
	default = "testacc19483"
}

data "alibabacloudstack_oss_clusters" "default" {
  provider = alibabacloudstack-common
}

resource "alibabacloudstack_oss_bucket" "default" {
  provider = alibabacloudstack-common
  bucket = "${var.name}"
  oss_cluster = "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}"
}

resource "alibabacloudstack_aqs_oss_scanconfig" "example" {
  enable                     = true
  start_time                 = "00:00:00"
  end_time                   = "23:59:59"
  scan_day_list              = [1, 2, 3, 4, 5, 6, 7]
  scan_mode                  = "1"
  bucket_name                = "${alibabacloudstack_oss_bucket.default.bucket}"
  decryption                 = "OSS"
  key_prefix                 = "test"
  key_suffix                 = ".py"
  last_modified_start_time   = "2023-01-01 00:00:00"
}

data "alibabacloudstack_aqs_oss_scanconfigs" "example" {
  ids = ["${alibabacloudstack_aqs_oss_scanconfig.example.id}"]
}

output "first_scan_config_id" {
  value = data.alibabacloudstack_aqs_oss_scanconfigs.example.scan_configs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of OssScanConfig IDs.
* `name_regex` - (Optional) A regex string to filter results by bucket name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of OssScanConfig IDs.
* `names` - A list of bucket names corresponding to the scan configs.
* `scan_configs` - A list of OssScanConfigs. Each element contains the following attributes:
  * `id` - The ID of the OssScanConfig.
  * `enable` - Whether scanning is enabled.
  * `start_time` - The start time of the scan task.
  * `end_time` - The end time of the scan task.
  * `scan_day_list` - The days of the week to scan.
  * `scan_mode` - The scan mode.
  * `bucket_name` - The name of the OSS bucket to scan.
  * `decryption` - The decryption method.
  * `key_suffix` - The suffix of the object keys to scan.
  * `key_prefix` - The prefix of the object keys to scan.
  * `last_modified_start_time` - The start time for scanning objects based on their modification time.
  * `bucket_count` - The count of buckets.
  * `real_time_incr` - Whether real-time increment scanning is enabled.
  * `all_key_prefix` - Whether all key prefixes are scanned.
  * `last_update_time` - The last update time.
