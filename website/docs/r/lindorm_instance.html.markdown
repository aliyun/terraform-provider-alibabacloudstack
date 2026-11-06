---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instance"
description: |-
  Provides a AlibabacloudStack Lindorm Instance resource.
---

# alibabacloudstack\_lindorm\_instance

Provides a Lindorm instance resource.

## Example Usage

### Basic Usage

```terraform
variable "name" {
    default = "tf-testacc97984"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
}

resource "alibabacloudstack_lindorm_instance" "example" {
  zone_id         = "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_alias  = "${var.name}"
  cpu_brand       = "Intel"
  disk_category   = "HDD"
  engine_type     = "lindorm"
  instance_type   = "lindorm.c.xlarge"
  vpc_id          = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id      = "${alibabacloudstack_vpc_vswitch.default.id}"
  lindorm_num     = 2
  local_disk_num  = 2
  local_disk_size = 400
}
```

### Multi-Zone Usage

```terraform
variable "name" {
  default = "tf-testacc-mz"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "10.0.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "10.0.1.0/24"
  zone_id      = "<primary_zone_id>"
}

resource "alibabacloudstack_vpc_vswitch" "standby" {
  vswitch_name = "${var.name}_vsw_standby"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "10.0.2.0/24"
  zone_id      = "<standby_zone_id>"
}

resource "alibabacloudstack_lindorm_instance" "example" {
  zone_id            = "<primary_zone_id>"
  instance_alias     = "${var.name}"
  cpu_brand          = "Intel"
  disk_category      = "LV"
  engine_type        = "lindorm"
  instance_type      = "lindorm.g1.2c8g"
  core_spec          = "lindorm.g1.2c2g"
  vpc_id             = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id         = "${alibabacloudstack_vpc_vswitch.default.id}"
  standby_zone_id    = "<standby_zone_id>"
  standby_vswitch_id = "${alibabacloudstack_vpc_vswitch.standby.id}"
  lindorm_num        = 3
  local_disk_size    = 450
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) The ID of the zone where you want to deploy the instance.
* `instance_alias` - (Required) The alias of the instance.
* `cpu_brand` - (Required, ForceNew) The brand of CPU that you want to use.
* `disk_category` - (Optional, ForceNew) The category of the disk.
* `engine_type` - (Required, ForceNew) The type of the engine that you want to use.
* `instance_type` - (Required) The specification of the instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC where you want to deploy the instance.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch that is associated with the specified VPC.
* `lindorm_num` - (Required) The number of Lindorm nodes. Minimum value: 2, which is enforced by the server.
* `local_disk_size` - (Required, ForceNew) The size of the local disk. Unit: GiB. Valid values: 100 to 5000.
* `local_disk_num` - (Optional) The number of local disks per node. Valid values: 1 to 10. Default value: 1. Takes effect only when `disk_category` is `HDD` or `SSD`; when `disk_category` is `LV`, the number of disks per node is fixed to 1 by the server.
* `standby_zone_id` - (Optional, ForceNew) The standby zone ID. When specified, the instance is deployed in multi-zone mode and `topology_type` is automatically set to `multizone`. Must be specified together with `standby_vswitch_id`.
* `standby_vswitch_id` - (Optional, ForceNew) The vSwitch ID in the standby zone. Must be specified together with `standby_zone_id`.
* `core_spec` - (Optional, ForceNew) The core specification of the instance. Retrieve the available specifications from the `alibabacloudstack_lindorm_instance_types` data source with `engine_type` set to `store`.
* `deletion_protection` - (Optional) Indicates whether deletion protection is enabled. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
* `instance_id` - The ID of the instance.
* `instance_status` - The status of the instance.
* `create_time` - The creation time of the instance.
* `instance_storage` - The storage capacity of the instance.
* `disk_usage` - The usage of the disk.
* `enable_fs` - Indicates whether the file system is enabled.
* `switch_l_proxy_flag` - Indicates whether the L_Proxy switch is enabled.
* `switch_ssl_encryption_flag` - Indicates whether SSL encryption is enabled.
* `ali_uid` - The UID of the Alibaba Cloud account.
* `arch_version` - The architecture version of the instance.
* `topology_type` - The topology type of the instance. Valid values: `1azone`, `multizone`. It is automatically determined by whether `standby_zone_id` is set.

## Timeouts

* `create` - (Defaults to 60 mins) Used when create the resource.
* `update` - (Defaults to 20 mins) Used when update the resource.
* `delete` - (Defaults to 40 mins) Used when delete the resource.

## Import

Lindorm instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_lindorm_instance.example <id>
```
