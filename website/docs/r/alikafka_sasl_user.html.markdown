---
subcategory: "Alikafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_sasl_user"
sidebar_current: "docs-alibabacloudstack-resource-alikafka-sasl_user"
description: |-
  Provides a Alibabacloudstack Alikafka Sasl User resource.
---

# alibabacloudstack\_alikafka\_sasl\_user

Provides an Alikafka sasl user resource.

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

resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = alibabacloudstack_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `username` - (Required, ForceNew) Username for the sasl user. The length should between 1 to 64 characters. The characters can only contain 'a'-'z', 'A'-'Z', '0'-'9', '_' and '-'.
* `password` - (Optional, Sensitive) Operation password. It may consist of letters, digits, or underlines, with a length of 1 to 64 characters. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. You have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a user with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `type` - (Optional, ForceNew, Available in 1.159.0+) The authentication mechanism. Valid values: `plain`, `scram`. Default value: `plain`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value is formate as `<instance_id>:<username>`.

## Import

Alikafka Sasl User can be imported using the id, e.g.

```
terraform import alibabacloudstack_alikafka_sasl_user.example <instance_id>:<username>
```

