---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vbr_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-vbr-attachments"
description: |-
  Provides a list of CEN transit router VBR attachments owned by an Alibaba Cloud account.
---

# alibabacloudstack\_cen\_transit_router_vbr_attachments

This data source provides a list of CEN transit router VBR attachments in an Alibaba Cloud account according to the specified filters.

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
The following arguments are supported:

* `ids` - (Optional) A list of VBR attachment IDs.
* `cen_id` - (Required) The ID of the CEN instance.
* `vbr_id` - (Required) The ID of the VBR instance.
* `tags` - (Optional) The tag of the resource.
* `tag_key` - (Optional) The tag key.
* `tag_value` - (Optional) The tag value.
* `transit_router_id` - (Required) The ID of the transit router.
* `name_regex` - (Optional) A regex string to filter results by VBR attachment name.
* `description_regex` - (Optional) A regex string to filter results by VBR attachment description.
## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `transitrouterattachments` - A list of VBR attachments.
* `id` - The ID of the VBR attachment.
* `creation_time` - The creation time of the VBR attachment.
* `resource_type` - The resource type of the VBR attachment.
* `status` - The status of the VBR attachment.
* `auto_publish_route_enabled` - Whether to enable automatic * route publishing.
* `transit_router_id` - The ID of the transit router.
* `vbr_id` - The ID of the VBR instance.
* `transit_router_attachment_name` - The name of the transit router attachment.
* `transit_router_attachment_description` - The description of the transit router attachment.
* `transit_router_attachment_id` - The ID of the transit router attachment.