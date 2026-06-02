---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_lifecycle_policy"
sidebar_current: "docs-Alibabacloudstack-resource-nas-lifecycle-policy"
description: |-
  Provides a nas Lifecyclepolicy resource.
---

# alibabacloudstack\_nas\_lifecyclepolicy

Provides a nas Lifecyclepolicy resource.

## 示例用法
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

## 参数参考

支持以下参数：
  * `lifecycle_policy_name` - (必填, 变更时强制重建) - 代表资源一级ID的资源属性字段
  * `storage_type` - (选填, 变更时强制重建) - 数据转储后的存储类型。默认值：InfrequentAccess（低频介质存储）
  * `file_system_id` - (必填, 变更时强制重建) - 文件系统ID。
  * `path` - (必填, 变更时强制重建) - 生命周期管理策略关联目录的绝对路径。仅支持关联单个目录。必须以正斜线（/）开头，并且是挂载点中真实存在的路径。> 建议您配置Paths.N，可以同时关联多个目录。
  * `recursive` - (选填) - 是否递归应用子路径。
  * `lifecycle_rule_name` - (必填) - 生命周期管理策略关联的管理规则。取值：- DEFAULT_ATIME_14：距今14天未访问的文件- DEFAULT_ATIME_30：距今30天未访问的文件- DEFAULT_ATIME_60：距今60天未访问的文件 - DEFAULT_ATIME_90：距今90天未访问的文件
  * `oss_bucket` - (必填, 变更时强制重建) - OSS存储桶名称。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `create_time` - 代表创建时间的资源属性字段
