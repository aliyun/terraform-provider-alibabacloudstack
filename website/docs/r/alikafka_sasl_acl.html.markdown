---
subcategory: "Alikafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_sasl_acl"
sidebar_current: "docs-alibabacloudstack-resource-alikafka-sasl_acl"
description: |-
  Provides a Alibabacloudstack Alikafka Sasl Acl resource.
---

# alibabacloudstack_alikafka_sasl_acl

Provides an Alikafka sasl acl resource.

## Example Usage

Basic Usage

```
variable "name" {
	default = "testalikafkasslacl"
}

variable "password" {
}

resource "alibabacloudstack_alikafka_topic" "default" {
  instance_id = "cluster-private-paas-default"
  topic = "${var.name}"
  remark = "topic-remark"
}


resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = "cluster-private-paas-default"
  username = "${var.name}"
  password = var.password
  type     = "scram"
}

resource "alibabacloudstack_alikafka_sasl_acl" "default" {
    instance_id =               "cluster-private-paas-default"
    username =                 "${alibabacloudstack_alikafka_sasl_user.default.username}"
    acl_resource_type =         "Topic"
    acl_resource_name =         "${alibabacloudstack_alikafka_topic.default.topic}"
    acl_resource_pattern_type = "LITERAL"
    acl_operation_type =        "Write"
}

```



## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) Username for the sasl user. The length should between 1 to 64 characters. The user should be an existed sasl user.
* `acl_resource_type` - (Required, ForceNew) Resource type for this acl. The resource type can only be "Topic" and "Group".
* `acl_resource_name` - (Required, ForceNew) Resource name for this acl. The resource name should be a topic or consumer group name.
* `acl_resource_pattern_type` - (Required, ForceNew) Resource pattern type for this acl. The resource pattern support two types "LITERAL" and "PREFIXED". "LITERAL": A literal name defines the full name of a resource. The special wildcard character "*" can be used to represent a resource with any name. "PREFIXED": A prefixed name defines a prefix for a resource.
* `acl_operation_type` - (Required, ForceNew) Operation type for this acl. The operation type can only be "Write" and "Read".
* `host` - (Optional, ForceNew) Host for the acl.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>`.
* `host` - The host of the acl.