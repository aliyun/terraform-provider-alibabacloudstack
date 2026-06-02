---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vbr_attachment"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-vbr-attachment"
description: |-
  提供CEN转发路由器VBR附加资源。
---

# alibabacloudstack\_cen\_transit_router_vbr_attachment

提供CEN转发路由器VBR附加资源。

## Example Usage

### 基本用法

```hcl
variable "name" {
  default = "tf-testaccrouter-vbr-attachment"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc-xxxxxxxxx"
  virtual_border_router_name = var.name
  vlan_id                    = 1
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  cen_id            = alibabacloudstack_cen_instance.default.id
  transit_router_id = alibabacloudstack_cen_instance.default.transit_router_id
  vbr_id            = alibabacloudstack_express_connect_virtual_border_router.default.id
}
```

## Argument Reference

支持以下参数：

* `cen_id` - （必选）CEN实例的ID。
* `transit_router_id` - （必选）转发路由器的ID。
* `vbr_id` - （必选）虚拟边界路由器（VBR）的ID。
* `transit_router_attachment_name` - （可选）转发路由器附加项的名称。
* `transit_router_attachment_description` - （可选）转发路由器附加项的描述。
* `route_table_propagation_enabled` - （可选，变更时强制重建）是否启用路由表传播。修改此参数将强制重新创建资源。
* `route_table_association_enabled` - （可选，变更时强制重建）是否启用路由表关联。修改此参数将强制重新创建资源。
* `tags` - （可选）要分配给资源的标签映射。

## Attributes Reference

导出以下属性：

* `id` - 资源的ID，格式为 `{cen_id}:{transit_router_id}:{transit_router_attachment_id}:{vbr_id}`。
* `transit_router_attachment_id` - 转发路由器VBR附加项的ID。
* `auto_publish_route_enabled` - 是否启用了自动路由发布。
* `charge_type` - 转发路由器VBR附加项的计费类型。
* `creation_time` - 转发路由器VBR附加项的创建时间。
* `resource_type` - 资源类型。取值：`VBR`。
* `status` - 转发路由器VBR附加项的状态。
* `vbr_owner_id` - VBR所有者的阿里云账号ID。

## Import

CEN转发路由器VBR附加项可以使用 cen_id、transit_router_id、transit_router_attachment_id 和 vbr_id（以冒号分隔）进行导入，例如：

```
$ terraform import alibabacloudstack_cen_transit_router_vbr_attachment.example cen-xxxxxx:tr-xxxxxx:tr-attach-xxxxxx:vbr-xxxxxx
```