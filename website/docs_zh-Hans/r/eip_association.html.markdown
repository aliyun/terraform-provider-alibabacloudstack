---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_association"
sidebar_current: "docs-Alibabacloudstack-eip-association"
description: |- 
  编排绑定弹性公网地址和云服务器（Ecs）实例
---

# alibabacloudstack_eip_association

使用Provider配置的凭证在指定的资源集下编排绑定弹性公网地址和云服务器（Ecs）实例

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccEipAssociation10702"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id           = alibabacloudstack_vpc.default.id
  cidr_block       = "10.1.1.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name             = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name        = var.name
  description = "New security group"
  vpc_id      = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  vswitch_id         = alibabacloudstack_vswitch.default.id
  image_id          = data.alibabacloudstack_images.default.images.0.id
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_ssd"
  instance_type     = local.default_instance_type_id

  security_groups = [alibabacloudstack_security_group.default.id]
  instance_name   = var.name
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_eip" "default" {
  name = var.name
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = alibabacloudstack_eip.default.id
  instance_id   = alibabacloudstack_instance.default.id
  force         = false
  instance_type = "EcsInstance"
}
```

## 参数说明

支持以下参数：

* `allocation_id` - (必填，变更时重建) 弹性公网IP实例的ID。
* `instance_id` - (必填，变更时重建) 要绑定弹性公网IP的实例ID。可以输入NAT网关、传统型负载均衡CLB实例、云服务器ECS实例、辅助弹性网卡实例、高可用虚拟IP实例或IP地址的ID。
* `force` - (选填，变更时重建) 当弹性公网IP绑定了NAT网关，且NAT网关添加了DNAT或SNAT条目时，是否强制解绑弹性公网IP。取值范围：
  * **false** (默认值)：不强制解绑弹性公网IP。
  * **true**：强制解绑弹性公网IP。
* `instance_type` - (选填，变更时重建) 要绑定弹性公网IP的实例类型。取值范围：
  * **Nat**：NAT网关。
  * **SlbInstance**：传统型负载均衡CLB。
  * **EcsInstance** (默认值)：云服务器ECS。
  * **NetworkInterface**：辅助弹性网卡。
  * **HaVip**：高可用虚拟IP。
  * **IpAddress**：IP地址。
  > 默认要绑定弹性公网IP的实例类型为**EcsInstance**，如果您需要绑定弹性公网IP的实例类型不为**EcsInstance**，则该值必填。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `allocation_id` - 弹性公网IP实例的ID。
* `instance_id` - 与弹性公网IP关联的实例ID。
* `instance_type` - 与弹性公网IP关联的实例类型。取值范围：
  * **Nat**：NAT网关。
  * **SlbInstance**：传统型负载均衡CLB。
  * **EcsInstance** (默认值)：云服务器ECS。
  * **NetworkInterface**：辅助弹性网卡。
  * **HaVip**：高可用虚拟IP。
  * **IpAddress**：IP地址。
  > 默认要绑定弹性公网IP的实例类型为**EcsInstance**，如果您需要绑定弹性公网IP的实例类型不为**EcsInstance**，则该值必填。