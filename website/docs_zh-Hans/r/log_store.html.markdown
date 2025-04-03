---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_store"
sidebar_current: "docs-alibabacloudstack-resource-log-store"
description: |-
  编排日志告警的日志库
---

# alibabacloudstack_log_store

使用Provider配置的凭证在指定的资源集编排日志告警的日志库。
日志库(Log Store)是日志服务中用于收集、存储和查询日志数据的单元。每个日志库属于一个项目，每个项目可以创建多个日志库。[更多详情](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/sls/enterprise-ascm-developer-guide/CreateLogstore.html?spm=a2c4g.14484438.10001.307)

## 示例用法

### 基础用法

要调用此资源，您需要在provider参数中设置sls的endpoint地址
```
provider "alibabacloudstack" {
  endpoints {
    sls_endpoint = "var.sls_openapi_endpoint"
  }
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
加密用法
```
provider "alibabacloudstack" {
  endpoints {
    sls_endpoint = "var.sls_openapi_endpoint"
  }
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


## 参数参考

支持以下参数：

* `project` - (必填，变更时重建) 日志库所属的项目名称。
* `name` - (必填，变更时重建) 日志库名称，在同一项目中必须唯一。
* `retention_period` - (可选) 数据保留时间(以天为单位)。有效值范围：[1-3650]。默认值为 `30`。当值为 `3650` 时，日志库数据将永久存储。
* `shard_count` - (必填) 该日志库中的分片数量。
* `auto_split` - (可选) 是否自动拆分分片。默认为 `false`。
* `max_split_shard_count` - (可选) 自动拆分的最大分片数，范围为 1 到 64。当 `auto_split` 为 true 时必须指定此参数。
* `append_meta` - (可选) 是否自动附加日志元数据。元数据包括日志接收时间和客户端 IP 地址。默认为 `true`。
* `enable_web_tracking` - (可选) 是否启用 Web 跟踪功能。默认为 `false`。
* `encryption`(可选) 启用加密。默认为 `false`
* `encrypt_type` (可选) 支持的加密类型，仅支持 `default(sm4_gcm)` 和 `aes_gcm`
* `cmk_key_id` (可选) 用户主密钥 ID。
* `arn`   (可选) 角色 ARN。

## 属性参考

导出以下属性：

* `id` - 日志项目的 ID。格式为 `<project>:<name>`。
* `project` - 项目名称。
* `name` - 日志库名称。
* `retention_period` - 数据保留时间。
* `shard_count` - 分片数量。
* `auto_split` - 是否自动拆分分片。
* `max_split_shard_count` - 自动拆分的最大分片数。
* `append_meta` - 是否自动附加日志元数据。
* `enable_web_tracking` - 是否启用 Web 跟踪。
* `shards` - 分片属性。
  * `id` - 分片的 ID。
  * `status` - 分片状态，只有两个状态：`readwrite` 和 `readonly`。
  * `begin_key` - 分片范围的起始值(MD5)，包含在分片范围内。
  * `end_key` - 分片范围的结束值(MD5)，不包含在分片范围内。

## 导入

日志库可以通过 ID 导入，例如：

```bash
$ terraform import alibabacloudstack_log_store.example tf-log:tf-log-store
```