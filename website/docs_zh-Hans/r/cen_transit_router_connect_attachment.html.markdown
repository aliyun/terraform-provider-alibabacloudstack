---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachment"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-connect-attachment"
description: |-
  创建CEN转发路由器连接附件
---

# alibabacloudstack_cen_transit_router_connect_attachment

提供CEN转发路由器连接附件资源。

> **注意**: 当创建alibabacloudstack_cen_transit_router_connect_attachment资源时，要求cen_instance中已存在可用的alibabacloudstack_cen_transit_router_vbr_attachment资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testaccrouter_connect_attachment76718"
}

resource "alibabacloudstack_cen_instance" "default" {
  description       = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  physical_connection_id     = ""
  vlan_id                    = 1
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  virtual_border_router_name = var.name
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  vbr_id            = alibabacloudstack_express_connect_virtual_border_router.default.id
  cen_id            = alibabacloudstack_cen_instance.default.id
  transit_router_id = alibabacloudstack_cen_instance.default.transit_router_id
}



resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" {
  transit_router_id              = alibabacloudstack_cen_instance.default.transit_router_id
  transit_router_attachment_name = var.name
  depends_on = [
    "alibabacloudstack_cen_transit_router_vbr_attachment.default"
  ]
  cen_id = alibabacloudstack_cen_instance.default.id
}
```

## 参数说明

支持以下参数：

* `cen_id` - (必填, 变更时重建) 云企业网实例ID。您可以在创建CEN实例后获取该ID。
* `transit_router_id` - (必填, 变更时重建) 转发路由器实例ID。您可以在创建转发路由器后获取该ID。
* `transit_router_attachment_name` - (可选, 变更时重建) 网络实例连接的名称。名称长度为1-128个字符，可以包含中文、英文、数字、下划线（_）和短划线（-），但不能以`http://`或`https://`开头。

## 属性说明

导出以下属性：

* `id` - 资源ID，格式为`{CenId:TransitRouterId:TransitRouterAttachmentId}`。
* `creation_time` - 网络实例连接的创建时间。时间按照ISO 8601标准表示，并使用UTC时间。格式为：YYYY-MM-DDThh:mmZ。
* `protocol` - 协议类型。可能的值包括：BGP、Static等。
* `resource_id` - 网络实例连接所关联的网络实例ID。
* `resource_owner_id` - 网络实例所属的账号ID。
* `resource_type` - 网络实例连接所关联的网络实例类型。
* `status` - 网络实例连接的状态。
* `transit_router_attachment_id` - 网络实例连接ID。
* `transport_type` - 传输类型。

## 导入

CEN转发路由器VPC附件可以通过ID导入，例如：

```shell
$ terraform import alibabacloudstack_cen_transit_router_vpc_attachment.default cen-abc12345678900001:tr-abc12345678900001:tr-attach-abc12345678900001:vpc-abc12345678900001
```