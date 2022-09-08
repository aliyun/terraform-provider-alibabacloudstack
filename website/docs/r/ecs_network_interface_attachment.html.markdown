---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_network_interface_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ecs-network-interface-attachment"
description: |-
  Provides a Alibabacloudstack ECS Network Interface Attachment resource.
---

# alibabacloudstack\_ecs\_network\_interface\_attachment 

Provides a ECS Network Interface Attachment resource.

For information about ECS Network Interface Attachment and how to use it, see [What is Network Interface Attachment](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/ecs/enterprise-developer-guide/AttachNetworkInterface.html?spm=a2c4g.14484438.10001.364).

-> **NOTE:** Available in v1.123.1+.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

locals {
  instance_type_id = sort(data.alibabacloudstack_instance_types.default.ids)[0]
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    security_groups = ["${alibabacloudstack_security_group.default.id}"]

    instance_type = "${local.instance_type_id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alibabacloudstack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
    instance_id = "${alibabacloudstack_instance.default.id}"
    network_interface_id = "${alibabacloudstack_network_interface.default.id}"
}


```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The instance id.
* `network_interface_id` - (Required, ForceNew) The network interface id.
* `trunk_network_instance_id` - (Optional) The trunk network instance id.
* `wait_for_network_configuration_ready` - (Optional) The wait for network configuration ready.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Network Interface Attachment. The value is formatted `<network_interface_id>:<instance_id>`.

## Import

ECS Network Interface Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_network_interface_attachment.example eni-abcd1234:i-abcd1234
```
