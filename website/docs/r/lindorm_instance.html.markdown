---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instance"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-instance"
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
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  description     = "modify_description"
  vswitch_name   = "tf-testaccvpcvswitch97984"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}

resource "alibabacloudstack_lindorm_instance" "example" {
  zone_id               = "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_alias        = "${var.name}"
  cpu_brand             = "Intel"
  disk_category         = "HHD"
  engine_type           = "tsdb"
  instance_type         = "lindorm.c.xlarge"
  vpc_id                = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id            = "${alibabacloudstack_vpc_vswitch.default.id}"
  lindorm_num           = 2
  local_disk_num        = 2
  local_disk_size       = "400"
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
* `lindorm_num` - (Required) The number of Lindorm nodes.
* `local_disk_num` - (Optional) The number of local disks. Valid values: 1 to 10. Default value: 1.
* `local_disk_size` - (Required, ForceNew) The size of the local disk. Unit: GiB.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
* `instance_id` - The ID of the instance.
* `instance_status` - The status of the instance.
* `create_time` - The creation time of the instance.
* `instance_storage` - The storage capacity of the instance.
* `deletion_protection` - Indicates whether deletion protection is enabled.
* `disk_usage` - The usage of the disk.
* `enable_fs` - Indicates whether the file system is enabled.
* `switch_l_proxy_flag` - Indicates whether the L_Proxy switch is enabled.
* `switch_ssl_encryption_flag` - Indicates whether SSL encryption is enabled.
* `ali_uid` - The UID of the Alibaba Cloud account.

## Import

Lindorm instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_lindorm_instance.example li-12345678
```