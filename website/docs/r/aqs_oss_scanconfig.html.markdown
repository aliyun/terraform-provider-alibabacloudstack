---
subcategory: "Server Guard"
layout: "alibabacloudstack"
page_title: "Alibaba Cloud Stack: alibabacloudstack_aqs_oss_scanconfig"
sidebar_current: "docs-Alibabacloudstack-resource-aqs-oss-scanconfig"
description: |-
  Provides a Alibaba Cloud Stack AQS OSS Scan Config resource.
---

# alibabacloudstack_aqs_oss_scanconfig

Provides an AQS OSS Scan Config resource.

For information about AQS OSS Scan Config and how to use it, see [What is Oss Scan Config](https://www.alibabacloud.com/help/en/security-center/latest/createossconfig).

> Note: Orchestrating this resource requires installing a specific version (Version 3.18.6-6.6.0.2 and later versions) of AQS.


## Example Usage

Basic Usage

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
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required) Whether to enable scanning. Valid values: `true` and `false`.
* `start_time` - (Required) The start time of the scan task. Format: `HH:MM:SS`. Example: `00:00:00`.
* `end_time` - (Required) The end time of the scan task. Format: `HH:MM:SS`. Example: `23:59:59`.
* `scan_day_list` - (Optional) The days of the week to scan. Valid values: 1 to 7, where 1 represents Monday and 7 represents Sunday.
* `scan_mode` - (Required) The scan mode. Valid values: `1` (regular scan) and `2` (full scan).
* `bucket_name` - (Required, ForceNew) The name of the OSS bucket to scan.
* `decryption` - (Required) The decryption method. Valid values: `OSS` and `No`.
* `key_prefix` - (Optional) The prefix of the object keys to scan.
* `key_suffix` - (Optional) The suffix of the object keys to scan.
* `last_modified_start_time` - (Optional) The start time for scanning objects based on their modification time. Format: `YYYY-MM-DD HH:MM:SS`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID.

## Import

AQS OSS Scan Config can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_aqs_oss_scanconfig.example <id>
```
