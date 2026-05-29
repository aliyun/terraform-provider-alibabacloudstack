---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vbr_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-vbr-attachments"
description: |-
  提供阿里云账户拥有的CEN转发路由器VBR附件列表。
---

# alibabacloudstack\_cen\_transit_router_vbr_attachments

该数据源根据指定的过滤条件提供阿里云账户中的CEN转发路由器VBR附件列表。

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccRouterVbrAttachmentsDatasource"
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
  transit_router_attachment_name = var.name
  transit_router_attachment_description = var.name
}

data "alibabacloudstack_cen_transit_router_vbr_attachments" "default" {
  cen_id            = alibabacloudstack_cen_instance.default.id
  transit_router_id = alibabacloudstack_cen_instance.default.transit_router_id
  vbr_id            = alibabacloudstack_express_connect_virtual_border_router.default.id
  name_regex        = alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name
}
```
## Argument Reference
支持以下参数：

* `ids` - (可选) VBR附件的ID列表。
* `cen_id` - (必选) CEN实例ID。
* `vbr_id` - (必选) VBR实例ID。
* `tags` - (可选) 资源的标签。
* `tag_key` - (可选) 标签键。
* `tag_value` - (可选) 标签值。
* `transit_router_id` - (必选) 转发路由器ID。
* `name_regex` - (可选) VBR附件名称的正则表达式。
* `description_regex` - (可选) VBR附件描述的正则表达式。
Attributes Reference
除了上述参数外，还导出以下属性：

* `transitrouterattachments` - VBR附件列表。
* `id` - VBR附件的ID。
* `creation_time` - 创建时间。
* `resource_type` - 资源类型。
* `status` - VBR附件状态。
* `auto_publish_route_enabled` - 是否启用自动发布路由。
* `transit_router_id` - 转发路由器ID。
* `vbr_id` - VBR实例ID。
* `transit_router_attachment_name` - 转发路由器附件名称。
* `transit_router_attachment_description` - 转发路由器附件描述。
* `transit_router_attachment_id` - 转发路由器附件ID。