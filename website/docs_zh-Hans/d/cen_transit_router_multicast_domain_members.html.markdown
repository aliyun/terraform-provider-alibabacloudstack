---
subcategory: "云企业网(CEN)"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_cen_transit_router_multicast_domain_member"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-multicast-domain-member"
description: |-
  提供阿里云专有云CEN转发路由器组播域成员资源。
---

# alibabacloudstack_cen_transit_router_multicast_domain_member

提供阿里云专有云CEN转发路由器组播域成员资源。

有关CEN转发路由器组播域成员的信息以及如何使用，请参阅[什么是转发路由器组播域成员](https://www.alibabacloud.com/help/doc-detail/)。

## 示例用法

基本用法

```hcl
variable "name" {
  default = "tf-testAccRouteTable"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc" "example" {
  vpc_name       = var.name
  cidr_block     = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "example" {
  vpc_id       = alibabacloudstack_vpc.example.id
  cidr_block   = "10.1.0.0/16"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name = var.name
}

resource "alibabacloudstack_cen_instance" "example" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_cen_transit_router" "example" {
  cen_id = alibabacloudstack_cen_instance.example.id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "example" {
  cen_id                                = alibabacloudstack_cen_instance.example.id
  transit_router_id                     = alibabacloudstack_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name  = var.name
}

resource "alibabacloudstack_network_interface" "example" {
  vswitch_id = alibabacloudstack_vswitch.example.id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_member" "example" {
  group_ip_address                      = "224.0.0.1"
  network_interface_id                  = alibabacloudstack_network_interface.example.id
  transit_router_multicast_domain_id    = alibabacloudstack_cen_transit_router_multicast_domain.example.id
  vswitch_id                            = alibabacloudstack_vswitch.example.id
  resource_type                         = "VPC"
}
```

## 参数参考

以下参数被支持:

* `group_ip_address` - (必填, 变更时新建) 组播IP地址。
* `transit_router_multicast_domain_id` - (必填, 变更时新建) 转发路由器组播域ID。
* `resource_type` - (必填, 变更时新建) 资源类型。有效值: `VPC`, `Connect`。
* `vswitch_id` - (可选, 变更时新建) 组播成员所属的交换机ID。当`resource_type`为`VPC`时必填。
* `network_interface_id` - (可选, 变更时新建) 网络接口ID。当`resource_type`为`VPC`时必填。
* `connect_peer_id` - (可选, 变更时新建) 连接对端ID。当`resource_type`为`Connect`时必填。
* `connect_attachment_id` - (可选, 变更时新建, Computed) 连接附件ID。当`resource_type`为`Connect`时必填。

## 属性参考

以下属性会被导出:

* `id` - 资源ID，格式为 `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<key>`。
* `status` - 组播成员的状态。
* `vswitch_id` - 组播成员所属的交换机ID。
* `network_interface_id` - 网络接口ID。
* `connect_peer_id` - 连接对端ID。
* `connect_attachment_id` - 连接附件ID。

## 导入说明

CEN转发路由器组播域成员可以通过ID导入，例如：

```bash
$ terraform import alibabacloudstack_cen_transit_router_multicast_domain_member.default 224.0.0.1:tr-mcast-domain-1234567890abcdef0:VPC:eni-1234567890abcdef0
```