---
subcategory: "ApsaraDB for Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_instance"
description: |-
  Provides a KVStore instance resource.
---

# alibabacloudstack_kvstore_instance
-> **NOTE:** Alias name has: `alibabacloudstack_redis_tairinstance`

Provides a KVStore (ApsaraDB for Redis / Memcache) instance resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "kvstoreinstanceexample"
}

data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details              = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  zone_id        = data.alibabacloudstack_zones.kv_zone.zones[0].id
  engine         = "Redis"
  engine_version = "5.0"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.kv_zone.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_kvstore_instance" "default" {
  tair_instance_name  = var.name
  instance_class      = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version      = "5.0"
  zone_id             = data.alibabacloudstack_zones.kv_zone.zones.0.id
  instance_type       = "Redis"
  vswitch_id          = alibabacloudstack_vswitch.default.id
  password            = var.password
  node_type           = "MASTER_SLAVE"
  maintain_start_time = "01:00Z"
  maintain_end_time   = "02:00Z"
  vpc_auth_mode       = "Open"
}
```

Multi-zone Usage

```hcl
variable "name" {
  default = "kvstoreinstancemultizoneexample"
}

data "alibabacloudstack_kvstore_zones" "default" {
  multi = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  engine         = "Redis"
  engine_version = "5.0"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[0]
}

resource "alibabacloudstack_kvstore_instance" "default" {
  tair_instance_name = var.name
  instance_class     = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version     = "5.0"
  instance_type      = "Redis"
  vswitch_id         = alibabacloudstack_vpc_vswitch.default.id
  zone_id            = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[0]
  secondary_zone_id  = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[1]
  node_type          = "double"
  password           = var.password
}
```

## Argument Reference

The following arguments are supported:

* `instance_class` - (Required) The instance class of the KVStore instance.
* `tair_instance_name` - (Optional) The name of the instance. It must be between 2 and 128 characters in length. Conflicts with `instance_name`.
* `instance_name` - (Optional, Deprecated) Field 'instance_name' is deprecated and will be removed in a future release. Please use new field `tair_instance_name` instead. Conflicts with `tair_instance_name`.
* `password` - (Optional, Sensitive) The password used to connect to the instance. If not specified, the instance can be accessed only after a password is set.
* `kms_encrypted_password` - (Optional) An KMS encrypted password used for the instance. If `password` is specified, this field will be ignored.
* `kms_encryption_context` - (Optional) An encryption context used to decrypt `kms_encrypted_password`. This is valid only when `kms_encrypted_password` is set.
* `engine_version` - (Optional, ForceNew) The database engine version. Default value is `5.0`.
* `zone_id` - (Optional, ForceNew) The ID of the primary availability zone where the instance resides. Conflicts with `availability_zone`.
* `secondary_zone_id` - (Optional, ForceNew) The secondary availability zone ID. If specified, a multi-zone instance is created; otherwise, a single-zone instance is created.
* `availability_zone` - (Optional, ForceNew, Deprecated) Field 'availability_zone' is deprecated and will be removed in a future release. Please use new field `zone_id` instead. Conflicts with `zone_id`.
* `instance_charge_type` - (Optional, Deprecated) Field 'instance_charge_type' is deprecated and will be removed in a future release, and not for any use now. Valid values: `PrePaid`, `PostPaid`.
* `instance_type` - (Optional, ForceNew) The engine type of the instance. Valid values: `Redis`, `Memcache`. Default value is `Redis`.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. If specified, the instance is deployed in a VPC network.
* `backup_id` - (Optional) The ID of the backup set if the instance is created based on another instance's backup.
* `security_ips` - (Optional) A set of IP addresses that are allowed to access the instance.
* `vpc_auth_mode` - (Optional) The VPC authentication mode. Valid values: `Open` (enables password authentication), `Close` (disables password authentication and enables password-free access). It takes effect only when the network type is VPC.
* `parameters` - (Optional) A set of parameters used to configure the instance. Each element supports the following:
  * `name` - (Required) The name of the parameter.
  * `value` - (Required) The value of the parameter.
* `maintain_start_time` - (Optional) The start time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `maintain_end_time` - (Optional) The end time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `cpu_type` - (Optional) The CPU type of the instance, e.g. `intel`.
* `node_type` - (Optional) The node type of the instance. Valid values: `MASTER_SLAVE`, `STAND_ALONE`, `double`, `readone`.
* `architecture_type` - (Optional, Deprecated) Field 'architecture_type' is not working properly, set through 'instance_class'. Valid values: `cluster`, `rwsplit`, `standard`.
* `series` - (Optional, Deprecated) Field 'series' is not working properly, set through 'instance_class'. Valid values: `community`, `enterprise`.
* `tde_status` - (Optional) Whether to enable Transparent Data Encryption (TDE). Valid values: `true`, `false`. TDE cannot be disabled once it is enabled.
* `enable_ssl` - (Optional) Whether to enable SSL for the instance. Valid values: `true`, `false`.
* `encryption_key` - (Optional, Deprecated) The encryption key used for TDE. TDE does not support the simultaneous use of `encryption_key` and `role_arn`. Conflicts with `role_arn`.
* `role_arn` - (Optional, Deprecated) The ARN of the RAM role used for TDE. TDE does not support the simultaneous use of `encryption_key` and `role_arn`. Conflicts with `encryption_key`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the KVStore instance.
* `connection_domain` - The internal endpoint of the instance.
* `private_ip` - The private IP address of the instance.

## Timeouts

* `create` - (Defaults to 45 mins) Used when create the resource.
* `update` - (Defaults to 30 mins) Used when update the resource.
* `delete` - (Defaults to 45 mins) Used when delete the resource.

## Import

alibabacloudstack_kvstore_instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_kvstore_instance.example <id>
```
