---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachments"
sidebar_current: "docs-alibabacloudstack-datasource-cen-transit-router-connect-attachments"
description: |-
  提供CEN转发路由器连接附件列表。
---

# alibabacloudstack\_cen\_transit\_router\_connect\_attachments

该数据源提供CEN转发路由器连接附件列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccrouter_connect_attachment33915"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "${var.name}"
  cen_instance_name = "${var.name}"
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

data "alibabacloudstack_cen_transit_router_connect_attachments" "example" {
  cen_id             = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
  transit_router_id  = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_id}"
}

output "first_attachment_id" {
  value = data.alibabacloudstack_cen_transit_router_connect_attachments.example.transitrouterattachments.0.transit_router_attachment_id
}
```

## 参数参考

以下参数被支持：

* `cen_id` - (必填) CEN实例的ID。
* `transit_router_id` - (必填) 转发路由器的ID。
* `ids` - (可选) 转发路由器连接附件ID列表。
* `name_regex` - (可选) 用于按附件名称过滤结果的正则表达式。

## 属性参考

以下属性被导出：

* `ids` - 转发路由器连接附件ID列表。
* `transitrouterattachments` - 转发路由器连接附件列表。每个元素包含以下属性：
  * `transit_router_id` - 转发路由器的ID。
  * `transit_router_attachment_name` - 转发路由器附件的名称。
  * `transit_router_attachment_id` - 转发路由器附件的ID。
  * `resource_id` - 资源的ID。
  * `creation_time` - 附件的创建时间。
  * `resource_type` - 资源的类型。
  * `resource_owner_id` - 资源所有者的ID。
  * `status` - 附件的状态。
  * `protocol` - 附件使用的协议。
  * `transport_type` - 附件的传输类型。