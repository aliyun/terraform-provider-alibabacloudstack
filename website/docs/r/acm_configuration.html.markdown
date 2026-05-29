---
subcategory: "Application Configuration Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_acm_configuration"
sidebar_current: "docs-Alibabacloudstack-acm-configuration"
description: |-
  Provides a acm Configuration resource.
---

# alibabacloudstack\_acm\_configuration

Provides a acm Configuration resource.

## Example Usage
```
variable "name" {
    default = "tftest49196"
}

variable "logical_id" {
  default = ":tf_testacmconfig_723"
}

resource "alibabacloudstack_edas_namespace" "default" {
  	description = "${var.name}"
	namespace_name = "${var.name}"
	namespace_logical_id = "${var.logical_id}"
}

resource "alibabacloudstack_acm_configuration" "default" {
  type = "text"
  tags = "aaaaaa,bbbbbbb"
  namespace_id = "ef8c9ec9-8b54-4be5-931b-7d8e5a21ca45"
  desc = "${var.name}"
  beta_ips = "192.168.1.1,192.168.1.2"
  data_id = "${var.name}"
  app_name = "${var.name}"
  content = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  group = "DEFAULT_GROUP"
}
```

## Argument Reference

The following arguments are supported:
  * `app_name` - (Optional) - The name of the application to which the Acm configuration belongs.
  * `content` - (Required) - Data content of the Acm configuration.
  * `data_id` - (Required, ForceNew) - Data ID of the Acm configuration.
  * `desc` - (Optional) - The description of the Acm configuration.
  * `group` - (Required, ForceNew) - The group of the Acm configuration.
  * `namespace_id` - (Required, ForceNew) - The ID of the edas namespace.
  * `tags` - (Optional) - The tags of the Acm configuration.
  * `type` - (Required) - The type of the Acm configuration.
  <!-- * `encrypt_algorithm` - (Optional) - The encrypt algorithm of the Acm configuration. -->
  * `beta_content` - (Optional) - The beta content of the beta environment.
  * `beta_app_name` - (Optional) - The app name of the beta environment.
  * `beta_ips` - (Optional) - The IPs of the beta environment.

