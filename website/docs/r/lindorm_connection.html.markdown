---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_lindorm_connection"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-connection"
description: |-
  Provides a AlibabacloudStack Lindorm Connection resource.
---

# alibabacloudstack\_lindorm\_connection

Provides a Lindorm connection resource.

## Example Usage

### Basic Usage

```terraform
variable "name" {
  default = "tf-testacc"
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
  vswitch_name   = "tf-testaccvpcvswitch"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}

resource "alibabacloudstack_lindorm_instance" "default" {
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

resource "alibabacloudstack_lindorm_connection" "example" {
  instance_id    = alibabacloudstack_lindorm_instance.default.id
  engine_type    = "lindorm"
  classic_network = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Lindorm instance.
* `engine_type` - (Required, ForceNew) The type of the engine.
* `classic_network` - (Optional) Whether to enable public network access. Default to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<instance_id>:<engine_type>`.
* `engine_username` - The username of the database engine.
* `engine_password` - The password of the database engine.
* `net_info_list` - A list of network connection information. Each element contains the following attributes:
  * `port` - Port number.
  * `connection_string` - Connection string.
  * `net_type` - Network type. Valid values: `0`(public network), `2`(classic network), `3`(VPC network).
  * `access_type` - Access type.

## Import

Lindorm connection can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_lindorm_connection.example li-12345678:lindorm
```