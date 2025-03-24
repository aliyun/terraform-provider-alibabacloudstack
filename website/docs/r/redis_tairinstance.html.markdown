---
subcategory: "Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_tairinstance"
sidebar_current: "docs-Alibabacloudstack-redis-tairinstance"
description: |- 
  Provides a redis Tairinstance resource.
---

# alibabacloudstack_redis_tairinstance
-> **NOTE:** Alias name has: `alibabacloudstack_kvstore_instance`

Provides a redis Tairinstance resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstances92773"
}

variable "kv_edition" {
    default = "community"
}

variable "kv_engine" {
    default = "Redis"
}

variable "password" {
}

data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  zone_id      = data.alibabacloudstack_zones.kv_zone.zones[0].id
  edition_type = var.kv_edition
  engine       = var.kv_engine
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

resource "alibabacloudstack_redis_tairinstance" "default" {
  tair_instance_name = var.name
  instance_class     = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version     = "5.0"
  zone_id           = data.alibabacloudstack_zones.kv_zone.zones.0.id
  instance_type     = "tair_rdb"
  vswitch_id        = alibabacloudstack_vswitch.default.id
  password          = var.password
  node_type         = "MASTER_SLAVE"
  architecture_type = "standard"
  maintain_start_time = "01:00Z"
  maintain_end_time   = "02:00Z"
  vpc_auth_mode      = "Open"
}
```

## Argument Reference

The following arguments are supported:

* `tair_instance_name` - (Optional) The name of the resource. It must be between 2 and 256 characters in length and start with a letter or number. It can include underscores, letters, and numbers.
* `password` - (Optional) The password used to connect to the instance. The password must be 8 to 32 characters long and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters (! @ # $ % ^ & * ( ) _ + - =).
* `kms_encrypted_password` - (Optional) An KMS encrypted password used for the instance. If `password` is specified, this field will be ignored.
* `kms_encryption_context` - (Optional) An encryption context used to decrypt `kms_encrypted_password`. This is valid only when `kms_encrypted_password` is set.
* `instance_class` - (Required) The instance type of the resource. For more information, see [Instance Types](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/instance-types).
* `engine_version` - (Optional, ForceNew) The database version. Default value is `5.0`. Rules for transferring parameters of different Tair product types:
  - `tair_rdb`: Compatible with Redis 5.0 and Redis 6.0 protocols, transmitted as `5.0` or `6.0`.
  - `tair_scm`: Persistent memory compatible with Redis 6.0 protocol, passed as `1.0`.
  - `tair_essd`: Disk (ESSD/SSD) compatible with Redis 4.0 and Redis 6.0 protocols, transmitted as `1.0` and `2.0` respectively.
* `zone_id` - (Optional, ForceNew) The ID of the availability zone where the instance resides.
* `availability_zone` - (Optional, ForceNew) The availability zone of the instance.
* `instance_charge_type` - (Optional) The billing method of the instance. Valid values are `PrePaid` and `PostPaid`. Default is `PostPaid`.
* `instance_type` - (Optional, ForceNew) The storage medium of the instance. Valid values: `tair_rdb`, `tair_scm`, `tair_essd`.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch.
* `backup_id` - (Optional) The ID of the backup set if the instance is created based on another instance's backup.
* `vpc_auth_mode` - (Optional) The VPC authentication mode. Valid values are `Open` (enables password authentication) and `Close` (disables password authentication and enables password-free access).
* `maintain_start_time` - (Optional) The start time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `maintain_end_time` - (Optional) The end time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `cpu_type` - (Optional) The CPU type of the resource. Valid values: `intel`.
* `node_type` - (Optional) The node type. Valid values:
  - `MASTER_SLAVE`: High availability (dual copy)
  - `STAND_ALONE`: Single copy
  - `double`: Double copy
  - `single`: Single copy
* `architecture_type` - (Optional) The architecture type of the instance. Valid values: `cluster`, `standard`, `rwsplit`.
* `series` - (Optional, ForceNew) The series of the instance.
* `tde_status` - (Optional) The status of Transparent Data Encryption (TDE).
* `encryption_key` - (Optional) The encryption key used for encrypting the instance data.
* `role_arn` - (Optional) The ARN of the RAM role.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Tair instance.
* `connection_domain` - The internal endpoint of the instance.
* `private_ip` - The private IP address of the instance.
* `vpc_auth_mode` - The VPC authentication mode. Valid values: `Open` (enables password authentication), `Close` (disables password authentication and enables password-free access).
* `maintain_start_time` - The start time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `maintain_end_time` - The end time of the maintenance window. The format is `HH:mmZ` (UTC time).
* `node_type` - The node type. Valid values:
  - `MASTER_SLAVE`: High availability (dual copy)
  - `STAND_ALONE`: Single copy
  - `double`: Double copy
  - `single`: Single copy
* `architecture_type` - The architecture type of the instance. Valid values: `cluster`, `standard`, `rwsplit`.
* `series` - The series of the instance.