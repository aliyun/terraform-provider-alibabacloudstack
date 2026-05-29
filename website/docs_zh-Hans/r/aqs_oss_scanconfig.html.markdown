---
subcategory: "阿里云并行文件系统(APFS)"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_aqs_oss_scanconfig"
sidebar_current: "docs-alibabacloudstack-resource-aqs-oss-scanconfig"
description: |-
  提供阿里云专有云安骑士 OSS扫描配置资源。
---

# alibabacloudstack_aqs_oss_scanconfig

提供一个安骑士 OSS扫描配置资源。

> 注意：编排该资源需要安装特定版本（3.18.6-6.6.0.2及后续版本）的安骑士。

## 使用示例

基础用法

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

## 参数说明

以下参数是可设置的：

* `enable` - (必选) 是否启用扫描。可选值: `true` 和 `false`。
* `start_time` - (必选) 扫描任务的开始时间。格式: `HH:MM:SS`。例如: `00:00:00`。
* `end_time` - (必选) 扫描任务的结束时间。格式: `HH:MM:SS`。例如: `23:59:59`。
* `scan_day_list` - (可选) 扫描的星期几。有效值: 1到7，其中1代表星期一，7代表星期日。
* `scan_mode` - (必选) 扫描模式。有效值: `1` (常规扫描) 和 `2` (全量扫描)。
* `bucket_name` - (必选, 不可修改) 要扫描的OSS存储桶名称。
* `decryption` - (必选) 解密方式。有效值: `OSS` 和 `No`。
* `key_prefix` - (可选) 要扫描的对象键前缀。
* `key_suffix` - (可选) 要扫描的对象键后缀。
* `last_modified_start_time` - (可选) 根据修改时间扫描对象的起始时间。格式: `YYYY-MM-DD HH:MM:SS`。

## 属性说明

以下属性会被导出：

* `id` - 资源ID。

## 导入资源

AQS OSS扫描配置可以通过id导入，例如：

```shell
$ terraform import alibabacloudstack_aqs_oss_scanconfig.example <id>
```