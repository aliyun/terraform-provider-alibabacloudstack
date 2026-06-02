---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云: alibabacloudstack_cen_transit_router_connect_peers"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-connect-peers"
description: |-
  提供阿里云CEN转发路由器连接对等点列表，供 alibabacloudstack_cen_transit_router_connect_peer 资源使用。
---

# alibabacloudstack_cen_transit_router_connect_peers

本数据源根据指定的过滤条件提供阿里云账户中的CEN转发路由器连接对等点列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccTransitRouterConnectPeer48958"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
  transit_router_cidrs {
    cidr = "172.16.0.0/16"
  }
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  physical_connection_id = "pc-xxxxxxxxxxxxxxxxxxx"
  vlan_id =                    99
  local_gateway_ip =           "10.0.0.1"
  peer_gateway_ip =            "10.0.0.2"
  peering_subnet_mask =        "255.255.255.252"
  virtual_border_router_name = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" {
  cen_id = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
  transit_router_attachment_name = "${var.name}"
  depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

resource "alibabacloudstack_cen_transit_router_connect_peer" "default" {
  name                  = "${var.name}"
  local_ip              = "172.16.0.1"
  peer_ip               = "172.16.0.11"
  cen_id                = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_id}"
}

data "alibabacloudstack_cen_transit_router_connect_peers" "example" {
  cen_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.connect_attachment_id}"
}

output "first_connect_peer_id" {
  value = data.alibabacloudstack_cen_transit_router_connect_peers.example.peers.0.id
}
```

## 参数参考

以下参数被支持：

* `cen_id` - (必填, 变更时强制重建) 云企业网实例的ID。
* `connect_attachment_id` - (必填, 变更时强制重建) 转发路由器连接附件的ID。
* `ids` - (可选) 用于过滤结果的转发路由器连接对等点ID列表。
* `name_regex` - (可选) 用于按转发路由器连接对等点名称过滤结果的正则表达式字符串。
* `peer_name` - (可选, 变更时强制重建) 转发路由器连接对等点的名称。

## 属性参考

以下属性被导出：

* `peers` - 转发路由器连接对等点列表。每个元素包含以下属性：
  * `id` - 转发路由器连接对等点的ID。格式为 `{cen_id}:{connect_attachment_id}:{peer_id}`。
  * `cen_id` - 云企业网实例的ID。
  * `connect_attachment_id` - 转发路由器连接附件的ID。
  * `peer_id` - 转发路由器连接对等点的ID。
  * `name` - 转发路由器连接对等点的名称。
  * `local_ip` - 转发路由器连接对等点的本地IP地址。
  * `peer_ip` - 转发路由器连接对等点的对端IP地址。
  * `region_id` - 转发路由器连接对等点的区域ID。
  * `status` - 转发路由器连接对等点的状态。
  * `creation_time` - 转发路由器连接对等点的创建时间。