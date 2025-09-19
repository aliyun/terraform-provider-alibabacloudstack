---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachment"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-connect-attachment"
description: |-
  提供阿里云CEN转发路由器连接附件资源。
---

# alibabacloudstack\_cen\_transit_router_connect_attachment

提供CEN转发路由器连接附件资源。

> **注意**: 当创建alibabacloudstack_cen_transit_router_connect_attachment资源时，要求cen_instance中已存在可用的alibabacloudstack_cen_transit_router_vbr_attachment资源。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccrouter_connect_attachment33915"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  physical_connection_id = "YourPhysicalConnectionId"
  vlan_id =                    1
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
  transit_router_attachment_name = "${var.name}"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
  depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

```

## 参数参考

以下参数被支持：

* `cen_id` - (必填, ForceNew) CEN实例的ID。
* `transit_router_id` - (必填, ForceNew) 转发路由器的ID。
* `transit_router_attachment_name` - (可选) 转发路由器附件的名称。

## 属性参考

以下属性被导出：

* `id` - 资源的ID，格式为`{cen_id}:{transit_router_id}:{attachment_id}:{vpc_id}`。
* `transit_router_attachment_id` - 转发路由器附件的ID。
* `resource_id` - 资源ID。
* `creation_time` - 附件的创建时间。
* `resource_type` - 资源的类型。
* `resource_owner_id` - 资源所有者的ID。
* `status` - 附件的状态。

## 导入

CEN转发路由器VPC附件可以通过ID导入，例如：

```shell
$ terraform import alibabacloudstack_cen_transit_router_vpc_attachment.default cen-abc12345678900001:tr-abc12345678900001:tr-attach-abc12345678900001:vpc-abc12345678900001
```