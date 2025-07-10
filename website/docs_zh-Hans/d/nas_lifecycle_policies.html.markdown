---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_lifecyclepolicies"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-lifecyclepolicies"
description: |-
  提供阿里云账号下拥有的nas lifecyclepolicies列表。
---

# alibabacloudstack\_nas\_lifecyclepolicies

此数据源提供根据指定过滤条件列出的阿里云账号下的nas lifecyclepolicies资源列表。

## 示例用法
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

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 生命周期管理策略得ID列表。
  * `file_system_id` - (选填) - 文件系统ID。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `lifecycle_policies` - 生命周期管理策略列表。
    * `id` - 生命周期管理策略ID。
    * `create_time` - 代表创建时间的资源属性字段
    * `file_system_id` - 文件系统ID。
    * `lifecycle_policy_name` - 代表资源一级ID的资源属性字段
    * `lifecycle_rule_name` - 生命周期管理策略关联的管理规则。取值：- DEFAULT_ATIME_14：距今14天未访问的文件- DEFAULT_ATIME_30：距今30天未访问的文件- DEFAULT_ATIME_60：距今60天未访问的文件 - DEFAULT_ATIME_90：距今90天未访问的文件
    * `path` - 生命周期管理策略关联目录的绝对路径。仅支持关联单个目录。必须以正斜线（/）开头，并且是挂载点中真实存在的路径。> 建议您配置Paths.N，可以同时关联多个目录。
    * `recursive` - 是否递归应用子路径。
    * `storage_type` - 数据转储后的存储类型。默认值：InfrequentAccess（低频介质存储）
    * `oss_bucket` - oss存储桶名称
