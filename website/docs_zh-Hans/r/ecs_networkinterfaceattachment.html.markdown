---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterfaceattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-networkinterfaceattachment"
description: |- 
  编排绑定云服务器（Ecs）弹性网卡和实例
---

# alibabacloudstack_ecs_networkinterfaceattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_network_interface_attachment`

使用Provider配置的凭证在指定的资源集下编排绑定云服务器（Ecs）弹性网卡和实例。

## 示例用法

以下示例展示了如何创建一个VPC、交换机、安全组、实例、网络接口以及将网络接口附加到实例的完整流程。

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

data "alibabacloudstack_instance_types" "eni2" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount       = 2
  sorted_by        = "Memory"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

# 创建VPC
resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "192.168.0.0/24"
}

# 创建交换机
resource "alibabacloudstack_vswitch" "default" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  availability_zone = reverse(data.alibabacloudstack_zones.default.zones)[0].id
  vpc_id            = alibabacloudstack_vpc.default.id
}

# 创建安全组
resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

# 创建ECS实例
resource "alibabacloudstack_instance" "default" {
  availability_zone         = reverse(data.alibabacloudstack_zones.default.zones)[0].id
  security_groups          = [alibabacloudstack_security_group.default.id]
  instance_type            = data.alibabacloudstack_instance_types.eni2.instance_types[0].id
  system_disk_category     = "cloud_efficiency"
  image_id                 = data.alibabacloudstack_images.default.images[0].id
  instance_name            = var.name
  vswitch_id               = alibabacloudstack_vswitch.default.id
  internet_max_bandwidth_out = 10
}

# 创建弹性网卡
resource "alibabacloudstack_network_interface" "default" {
  name           = var.name
  vswitch_id     = alibabacloudstack_vswitch.default.id
  security_groups = [alibabacloudstack_security_group.default.id]
}

# 将弹性网卡附加到实例
resource "alibabacloudstack_network_interface_attachment" "default" {
  instance_id           = alibabacloudstack_instance.default.id
  network_interface_id  = alibabacloudstack_network_interface.default.id
}
```

## 参数参考

支持以下参数：

* `instance_id` - (必填, 变更时重建) - 要附加弹性网卡(ENI)的ECS实例ID。更改此值将强制创建新的附件。
* `network_interface_id` - (必填, 变更时重建) - 要附加到指定实例的弹性网络接口(ENI)的ID。更改此值将强制创建新的附件。

> **说明**: 
> - `instance_id` 和 `network_interface_id` 是必填参数，且更改它们中的任何一个都会导致创建一个新的资源实例。
> - 弹性网卡(ENI)只能附加到与之位于同一可用区和VPC中的实例。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `id` - ENI附件资源的唯一标识符。它格式化为 `<network_interface_id>:<instance_id>`。

> **说明**: 
> - `id` 是该资源的唯一标识符，用于在Terraform中引用该资源。
> - 此属性可以用于其他依赖于此资源的模块或配置中。