---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_network"
sidebar_current: "docs-Alibabacloudstack-resource-expressconnect-bgp-network"
description: |-
  提供高速通道 BGP 网络资源。
---

# alibabacloudstack\_expressconnect\_bgp\_network

提供高速通道 BGP 网络资源。

## 示例用法

```terraform
variable "name" {
  default = "tf-testaccexpressconnect-bgp-group"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
    physical_connection_id = var.physical_connection_id
    vlan_id = 1106
    local_gateway_ip = "10.0.0.1"
    peer_gateway_ip = "10.0.0.2"
    peering_subnet_mask = "255.255.255.252"
    virtual_border_router_name = var.name
    description = "BGP 网络测试"
}

resource "alibabacloudstack_expressconnect_bgp_network" "default" {
  dst_cidr_block = "10.10.0.0/24"
  router_id = alibabacloudstack_express_connect_virtual_border_router.default.id
}
```

## 参数参考

支持以下参数：

* `dst_cidr_block` - (必填) 需要和本地数据中心（IDC）互连的 VPC 或交换机的网段。
* `router_id` - (必填, 变更时强制重建) 路由器接口关联的路由器（VBR）ID。修改此参数会强制重新创建资源。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - BGP 网络的 ID。格式为 `<dst_cidr_block>:<router_id>`。
* `status` - BGP 网络的状态。

## Import

高速通道 BGP 网络可以使用 ID 导入，例如：

```
$ terraform import alibabacloudstack_expressconnect_bgp_network.example <dst_cidr_block>:<router_id>
```
