---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云: alibabacloudstack_cen_transit_router_connect_peer"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-connect-peer"
description: |-
  提供阿里云CEN转发路由器连接对等点资源。
---

# alibabacloudstack_cen_transit_router_connect_peer

提供CEN转发路由器连接对等点资源。

有关CEN转发路由器连接对等点的更多信息及使用方法，请参见[什么是转发路由器连接对等点](https://www.alibabacloud.com/help/doc-detail/65872.html)。

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

```

## 参数参考

以下参数被支持：

* `cen_id` - (必填, 变更时重建) 云企业网实例的ID。
* `connect_attachment_id` - (必填, 变更时重建) 转发路由器连接附件的ID。
* `local_ip` - (可选, 变更时重建) 转发路由器连接对等点的本地IP地址, 需在实例所在Transit Router VPC网段内。
* `name` - (可选, 变更时重建) 转发路由器连接对等点的名称。
* `peer_ip` - (必填, 变更时重建) 转发路由器连接对等点的对端IP地址。

## 属性参考

以下属性被导出：

* `id` - 资源的ID。格式为 `<cen_id>:<connect_attachment_id>:<peer_id>`。
* `peer_id` - 转发路由器连接对等点的ID。
* `region_id` - 转发路由器连接对等点的区域ID。

## 导入说明

CEN转发路由器连接对等点可以使用ID导入，例如：

```bash
$ terraform import alibabacloudstack_cen_transit_router_connect_peer.example cen-2cudwl7et3716ozh****:tr-attach-2cudwl7et3716ozh****:tr-cp-2cudwl7et3716ozh****
```