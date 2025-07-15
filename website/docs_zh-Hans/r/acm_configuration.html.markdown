---
subcategory: "ACM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_acm_configuration"
sidebar_current: "docs-Alibabacloudstack-acm-configuration"
description: |-
  Provides a acm Configuration resource.
---

# alibabacloudstack\_acm\_configuration

Provides a acm Configuration resource.

## 示例用法
```
variable "name" {
    default = "tftest49196"
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
  type = "text"
  tags = "aaaaaa,bbbbbbb"
  namespace_id = "${alibabacloudstack_edas_namespace.default.id}"
  desc = "${var.name}"
  beta_ips = "192.168.1.1,192.168.1.2"
  data_id = "${var.name}"
  app_name = "${var.name}"
  content = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  group = "DEFAULT_GROUP"
}
```

## 参数参考

支持以下参数：
  * `app_name` - (选填) - 配置归属的应用名称
  * `content` - (必填) - 配置内容
  * `data_id` - (必填, 强制新建) - 配置ID
  * `desc` - (选填) - 配置描述
  * `group` - (必填, 强制新建) - 分组
  * `beta_ips` - (选填) - beta发布的ip，多个ip使用“,”连接。
  * `namespace_id` - (必填, 强制新建) - Edas命名空间ID
  * `tags` - (选填) - 配置的标签
  * `type` - (必填) - 配置内容的格式
  <!-- * `encrypt_algorithm` - (选填) - 加密方式 -->
  * `beta_content` - (选填) - 发布到bate环境的配置内容。
  * `beta_app_name` - (选填) - 发布到bate环境的配置归属的应用名称。
  * `beta_ips` - (选填) - bate环境的 IPs。
