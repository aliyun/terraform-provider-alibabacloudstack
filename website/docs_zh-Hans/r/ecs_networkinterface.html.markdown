---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterface"
sidebar_current: "docs-Alibabacloudstack-ecs-networkinterface"
description: |- 
  编排云服务器（Ecs）弹性网卡。
---

# alibabacloudstack_ecs_networkinterface
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_network_interface`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）弹性网卡。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccecsnetwork_interface15831"
}

resource "alibabacloudstack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.0.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
  name       = "${var.name}"
  vpc_id     = alibabacloudstack_vpc.vpc.id
  cidr_block = "10.0.0.0/24"
  availability_zone = "cn-beijing-b"
}

resource "alibabacloudstack_security_group" "secgroup" {
  name        = "${var.name}"
  description = "Security Group"
  vpc_id      = alibabacloudstack_vpc.vpc.id
}

resource "alibabacloudstack_ecs_networkinterface" "default" {
  network_interface_name = "${var.name}-eni"
  vswitch_id             = alibabacloudstack_vswitch.vsw.id
  security_groups        = [alibabacloudstack_security_group.secgroup.id]
  primary_ip_address     = "10.0.0.10"
  private_ips_count      = 2
  description            = "Test ENI"
}
```

## 参数说明

支持以下参数：

* `network_interface_name` - (可选) 弹性网卡的名称。该名称可以包含2到128个字符，只能包含字母、数字或连字符(如“-”、“.”、“_”)，不能以连字符开头或结尾，不能以http://或https://开头。默认值为null。
* `vswitch_id` - (必填，变更时重建) 用于创建弹性网卡的交换机ID。
* `security_groups` - (必填) 要与弹性网卡关联的安全组ID列表。
* `primary_ip_address` - (可选，变更时重建) 弹性网卡的主要私有IP地址。如果不指定，阿里云将自动在交换机的CIDR块内分配一个。
* `private_ips` - (可选) 要分配给弹性网卡的次要私有IP地址列表。不要在同一弹性网卡资源块中同时使用`private_ips`和`private_ips_count`。
* `private_ips_count` - (可选) 要分配给弹性网卡的次要私有IP地址数量。不要在同一弹性网卡资源块中同时使用`private_ips`和`private_ips_count`。
* `description` - (可选) 弹性网卡的描述。该描述可以包含2到256个字符，不能以http://或https://开头。默认值为null。
* `tags` - (可选) 要分配给资源的标签映射。
* `mac_address` - (可选) 弹性网卡的MAC地址。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 弹性网卡的ID。
* `mac_address` - 弹性网卡的MAC地址。
* `network_interface_name` - 弹性网卡的名称。
* `primary_ip_address` - 弹性网卡的主要私有IP地址。
* `private_ips` - 分配给弹性网卡的所有私有IP地址列表，包括主要IP和任何次要IP。
* `private_ips_count` - 分配给弹性网卡的私有IP地址总数，包括主要IP和任何次要IP。