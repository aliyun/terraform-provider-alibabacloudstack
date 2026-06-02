---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_cen_transit_router_multicast_domain_source"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-multicast-domain-source"
description: |-
  提供阿里云专有云CEN转发路由器组播域源资源。
---

# alibabacloudstack_cen_transit_router_multicast_domain_source

提供阿里云专有云CEN转发路由器组播域源资源。

有关CEN转发路由器组播域源的信息以及如何使用，请参阅[什么是转发路由器组播域源](https://www.alibabacloud.com/help/doc-detail/)。

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

resource "alibabacloudstack_cen_transit_router_multicast_domain_source" "example" {
  group_ip_address                      = "224.0.0.1"
  network_interface_id                  = alibabacloudstack_network_interface.example.id
  transit_router_multicast_domain_id    = alibabacloudstack_cen_transit_router_multicast_domain.example.id
  vswitch_id                            = alibabacloudstack_vswitch.example.id
  resource_type                         = "VPC"
}
```

## 参数参考

以下参数被支持:

* `group_ip_address` - (必填) 组播IP地址。
* `transit_router_multicast_domain_id` - (必填, 变更后重建) 转发路由器组播域ID。修改此参数会强制重新创建资源。
* `resource_type` - (必填, 变更后重建) 资源类型。取值范围：`VPC`、`Connect`。修改此参数会强制重新创建资源。
* `vswitch_id` - (可选, 变更后重建, Computed) 交换机ID。当 `resource_type` 为 `VPC` 时必填。修改此参数会强制重新创建资源。此属性由 API 返回，无法手动设置。
* `network_interface_id` - (可选, 变更后重建) 弹性网卡 ENI ID。当 `resource_type` 为 `VPC` 时必填。修改此参数会强制重新创建资源。
* `connect_peer_id` - (可选, 变更后重建) Connect Peer ID。当 `resource_type` 为 `Connect` 时必填。修改此参数会强制重新创建资源。
* `connect_attachment_id` - (可选, 变更后重建, Computed) Connect Attachment ID。当 `resource_type` 为 `Connect` 时必填。修改此参数会强制重新创建资源。此属性由 API 返回，无法手动设置。

## 属性参考

以下属性会被导出:

* `id` - 资源ID，格式为 `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<key>`，其中当 `resource_type` 为 `VPC` 时 `<key>` 为 `network_interface_id`，当 `resource_type` 为 `Connect` 时 `<key>` 为 `connect_peer_id`。
* `status` - 组播源的状态。

## 导入说明

CEN转发路由器组播域源可以通过ID导入，例如：

```bash
$ terraform import alibabacloudstack_cen_transit_router_multicast_domain_source.example 224.0.0.1:tr-mcast-domain-1234567890abcdef0:VPC:eni-1234567890abcdef0
```
