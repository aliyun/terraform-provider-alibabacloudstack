---
subcategory: "ACM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_acm_configurations"
sidebar_current: "docs-Alibabacloudstack-datasource-acm-configurations"
description: |-
  提供阿里云账号下拥有的acm configurations列表。
---

# alibabacloudstack\_acm\_configurations

此数据源提供根据指定过滤条件列出的阿里云账号下的acm configurations资源列表。

## 示例用法
```
variable "name" {
  default = "tf_testacmconfig_723"
}

variable "logical_id" {
  default = "{region_id}:tf_testacmconfig_723"
}

resource "alibabacloudstack_edas_namespace" "default" {
  description = "${var.name}"
	namespace_name = "${var.name}"
	namespace_logical_id = "${var.logical_id}"
}

resource "alibabacloudstack_acm_configuration" "default" {
	app_name = "${var.name}"
	content = "test"
	data_id = "${var.name}"
	group = "DEFAULT_GROUP"
	type = "text"
	namespace_id = "${alibabacloudstack_edas_namespace.default.id}"
}


data "alibabacloudstack_acm_configurations" "default" {
  namespace_id = "${alibabacloudstack_edas_namespace.default.id}"
  data_id = "${var.name}"
  group_id = "DEFAULT_GROUP"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - acm配置ID的列表。
  * `data_id` - (必填) - acm配置的数据ID
  * `group` - (必填) - 数据分组
  * `app_name` - (选填) - acm配置归属的应用名称
  * `namespace_id` - (选填) - edas命名空间ID

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `configurations` - acm配置的列表。
    * `id` - acm配置ID。
    * `app_name` - 配置归属的应用名称
    * `content` - 配置内容
    * `data_id` - 配置ID
    * `desc` - 配置描述
    * `group` - 分组
    * `message_digest` - 配置的消息摘要
    * `namespace_id` - edas命名空间ID
    * `tags` - 配置的标签
    * `type` - 配置内容的格式
    * `ud_version` - 数据版本。
