---
subcategory: "Application Configuration Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_acm_configurations"
sidebar_current: "docs-Alibabacloudstack-datasource-acm-configurations"
description: |-
  Provides a list of acm configurations owned by an alibabacloudstack account.
---

# alibabacloudstack\_acm\_configurations

This data source provides a list of acm configurations in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf_testacmconfig_723"
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

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - The IDs of the Acm configurations.
  * `data_id` - (Required) - The data ID of the Acm configuration, which is unique in the same group.
  * `group` - (Required) - The group of the Acm configuration.
  * `app_name` - (Optional) - The name of the application to which the Acm configuration belongs.
  * `namespace_id` - (Optional) - The ID of the edas namespace.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `configurations` - A list of acm configurations.
    * `id` - The ID of the Acm configuration.
    * `app_name` - A name of the application to which the Acm configuration belongs.
    * `content` - Data content of the Acm configuration.
    * `data_id` - Data ID of the Acm configuration.
    * `desc` - The description of the Acm configuration.
    * `group` - The group of the Acm configuration.
    * `message_digest` - The data md5 of the Acm configuration.
    * `namespace_id` - The ID of the edas namespace.
    * `tags` - The tags of the Acm configuration.
    * `type` - The type of the Acm configuration.
    * `ud_version` - The data version of the Acm configuration.
