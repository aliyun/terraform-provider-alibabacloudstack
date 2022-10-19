---
subcategory: "Alikafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_sasl_acl"
sidebar_current: "docs-alibabacloudstack-resource-alikafka-sasl_acl"
description: |-
  Provides a Alibabacloudstack Alikafka Sasl Acl resource.
---

# alibabacloudstack\_alikafka\_sasl\_acl

Provides an ALIKAFKA sasl acl resource.

## Example Usage

Basic Usage

```
variable "username" {
  default = "testusername"
}

variable "password" {
  default = "testpassword"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_alikafka_instance" "default" {
  name        = "tf-testacc-alikafkainstance"
  topic_quota = "50"
  disk_type   = "1"
  disk_size   = "500"
  deploy_type = "5"
  io_max      = "20"
  vswitch_id  = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_alikafka_topic" "default" {
  instance_id = alibabacloudstack_alikafka_instance.default.id
  topic       = "test-topic"
  remark      = "topic-remark"
}

resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = alibabacloudstack_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}

resource "alibabacloudstack_alikafka_sasl_acl" "default" {
  instance_id               = alibabacloudstack_alikafka_instance.default.id
  username                  = alibabacloudstack_alikafka_sasl_user.default.username
  acl_resource_type         = "Topic"
  acl_resource_name         = alibabacloudstack_alikafka_topic.default.topic
  acl_resource_pattern_type = "LITERAL"
  acl_operation_type        = "Write"
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

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>`.
* `host` - The host of the acl.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_alikafka_sasl_acl.acl alikafka_post-cn-123455abc:username:Topic:test-topic:LITERAL:Write
```
