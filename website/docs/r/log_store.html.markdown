---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_store"
sidebar_current: "docs-alibabacloudstack-resource-log-store"
description: |-
  Provides a Alibabacloudstack log store resource.
---

# alibabacloudstack\_log\_store

The log store is a unit in Log Service to collect, store, and query the log data. Each log store belongs to a project,
and each project can create multiple Logstores. [Refer to details](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/sls/enterprise-ascm-developer-guide/CreateLogstore.html?spm=a2c4g.14484438.10001.307)

## Example Usage

Basic Usage

To invoke this resource, you need to set the provider parameter: sls_openapi_endpoint
```
provider "alibabacloudstack" {
  sls_openapi_endpoint = "var.sls_openapi_endpoint"
  ...
}

resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alibabacloudstack_log_store" "example" {
  project               = alibabacloudstack_log_project.example.name
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
```
Encrypt Usage
```
provider "alibabacloudstack" {
  sls_openapi_endpoint = "var.sls_openapi_endpoint"
  ...
}

resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alibabacloudstack_log_store" "example" {
    project                 = alibabacloudstack_log_project.example.name
    name                    = "tf-log-store"
    retention_period        = "30"
    shard_count             = 3
    enable_web_tracking     = false
    auto_split              = true
    max_split_shard_count   = "64"
    append_meta             = true
    encryption              = true
    encrypt_type            = "aes_gcm"
    arn                     = "acs:ram::0000000000000080:role/ascm-role-00-0-0000"
    cmk_key_id              = "your_cmk_key_id"
}
```


## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `name` - (Required, ForceNew) The log store, which is unique in the same project.
* `retention_period` - (Optional) The data retention time (in days). Valid values: [1-3650]. Default to `30`. Log store data will be stored permanently when the value is `3650`.
* `shard_count` - (Optional) The number of shards in this log store. Default to 2. You can modify it by "Split" or "Merge" operations. [Refer to details](https://www.alibabacloud.com/help/doc-detail/28976.htm)
* `auto_split` - (Optional) Determines whether to automatically split a shard. Default to `false`.
* `max_split_shard_count` - (Optional) The maximum number of shards for automatic split, which is in the range of 1 to 64. You must specify this parameter when autoSplit is true.
* `append_meta` - (Optional) Determines whether to append log meta automatically. The meta includes log receive time and client IP address. Default to `true`.
* `enable_web_tracking` - (Optional) Determines whether to enable Web Tracking. Default `false`.
* `encryption`(Optional) enable encryption. Default `false`
* `encrypt_type` (Optional) Supported encryption type, only supports `default(sm4_gcm)`,` aes_gcm`
* `cmk_key_id` (Optional) User master key id.
* `arn`   (Optional) role arn.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - Log store name.
* `retention_period` - The data retention time.
* `shard_count` - The number of shards.
* `auto_split` - Determines whether to automatically split a shard.
* `max_split_shard_count` - The maximum number of shards for automatic split.
* `append_meta` - Determines whether to append log meta automatically.
* `enable_web_tracking` - Determines whether to enable Web Tracking.
* `shards` - The shard attribute.
  * `id` - The ID of the shard.
  * `status` - Shard status, only two status of readwrite and readonly.
  * `begin_key` - The begin value of the shard range(MD5), included in the shard range.
  * `end_key` - The end value of the shard range(MD5), not included in shard range.

## Import

Log store can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_log_store.example tf-log:tf-log-store
```
