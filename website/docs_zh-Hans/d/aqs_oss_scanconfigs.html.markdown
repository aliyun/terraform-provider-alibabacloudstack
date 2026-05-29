---
subcategory: "防暴力破解安全服务"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_aqs_oss_scanconfigs"
sidebar_current: "docs-Alibabacloudstack-datasource-aqs-oss-scanconfigs"
description: |-
  提供阿里云专有云用户可用的AQS OSS扫描配置列表。
---

# alibabacloudstack_aqs_oss_scanconfigs

本数据源根据指定的过滤条件，提供阿里云专有云账户中AQS OSS扫描配置的列表。

> 注意：编排该资源需要安装特定版本（3.18.6-6.6.0.2及后续版本）的安骑士。

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

## 参数说明

以下参数是可设置的：

* `ids` - (可选) OSS扫描配置ID列表。
* `name_regex` - (可选) 用于按存储桶名称过滤结果的正则表达式。

## 属性说明

以下属性会被导出：

* `ids` - OSS扫描配置ID列表。
* `names` - 与扫描配置对应的存储桶名称列表。
* `scan_configs` - OSS扫描配置列表。每个元素包含以下属性：
  * `id` - OSS扫描配置的ID。
  * `enable` - 是否启用了扫描。
  * `start_time` - 扫描任务的开始时间。
  * `end_time` - 扫描任务的结束时间。
  * `scan_day_list` - 扫描的星期几。
  * `scan_mode` - 扫描模式。
  * `bucket_name` - 要扫描的OSS存储桶名称。
  * `decryption` - 解密方式。
  * `key_suffix` - 要扫描的对象键后缀。
  * `key_prefix` - 要扫描的对象键前缀。
  * `last_modified_start_time` - 根据修改时间扫描对象的起始时间。
  * `bucket_count` - 存储桶数量。
  * `real_time_incr` - 是否启用了实时增量扫描。
  * `all_key_prefix` - 是否扫描所有键前缀。
  * `last_update_time` - 最后更新时间。