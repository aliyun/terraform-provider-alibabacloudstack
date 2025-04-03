---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_vservergroup"
sidebar_current: "docs-Alibabacloudstack-slb-vservergroup"
description: |- 
  Provides a slb Vservergroup resource.
---

# alibabacloudstack_slb_vservergroup
-> **NOTE:** Alias name has: `alibabacloudstack_slb_server_group`

Provides a slb Vservergroup resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccslbv_server_group93962"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "instance" {
  image_id                   = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = ["${alibabacloudstack_security_group.default.id}"]
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name               = "${var.name}"
  address_type       = "internet"
  specification      = "slb.s2.small"
  vswitch_id         = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_vservergroup" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name             = "${var.name}"

  servers {
    server_ids = ["${alibabacloudstack_instance.instance[0].id}", "${alibabacloudstack_instance.instance[1].id}"]
    port       = 100
    weight     = 10
    type       = "ecs"
  }

  servers {
    server_ids = ["${alibabacloudstack_instance.instance.*.id}"]
    port       = 80
    weight     = 100
    type       = "eni"
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the SLB instance.
* `name` - (Optional) Name of the virtual server group.
* `vserver_group_name` - (Optional) Name of the VServer Group.
* `servers` - (Optional) A list of ECS instances to be added. At most 20 ECS instances can be supported in one resource. It contains several sub-fields as follows:
  * `server_ids` - (Required) A list of backend server IDs (ECS instance IDs).
  * `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
  * `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.
  * `type` - (Optional) Type of the backend server. Valid values: `ecs`, `eni`. Default to `ecs`.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the VServer Group.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new VServer Group.
* `name` - The name of the VServer Group.
* `vserver_group_name` - The name of the VServer Group.
* `servers` - A list of ECS instances that have been added.