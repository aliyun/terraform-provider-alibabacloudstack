---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_vbr_attachment"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-vbr-attachment"
description: |-
  Provides a CEN transit router VBR attachment resource.

-> **Note:** This resource can also be referred to by the following alias:
-> - `apsarastack_cen_transit_router_vbr_attachment`
---

# alibabacloudstack\_cen\_transit_router_vbr_attachment

Provides a CEN transit router VBR attachment resource.

## Example Usage

### Basic Usage

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

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN instance.
* `transit_router_id` - (Required) The ID of the transit router.
* `vbr_id` - (Required) The ID of the Virtual Border Router (VBR).
* `transit_router_attachment_name` - (Optional) The name of the transit router attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router attachment.
* `route_table_propagation_enabled` - (Optional, ForceNew) Specifies whether to enable route table propagation. Modifying this parameter will force the recreation of the resource.
* `route_table_association_enabled` - (Optional, ForceNew) Specifies whether to enable route table association. Modifying this parameter will force the recreation of the resource.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `{cen_id}:{transit_router_id}:{transit_router_attachment_id}:{vbr_id}`.
* `transit_router_attachment_id` - The ID of the transit router VBR attachment.
* `auto_publish_route_enabled` - Indicates whether automatic route publishing is enabled.
* `charge_type` - The billing method of the transit router VBR attachment.
* `creation_time` - The time when the transit router VBR attachment was created.
* `resource_type` - The resource type. Valid value: `VBR`.
* `status` - The status of the transit router VBR attachment.
* `vbr_owner_id` - The Alibaba Cloud account ID of the VBR owner.

## Import

CEN Transit Router VBR Attachment can be imported using the cen_id, transit_router_id, transit_router_attachment_id and vbr_id (separated by colons), e.g.

```
$ terraform import alibabacloudstack_cen_transit_router_vbr_attachment.example cen-xxxxxx:tr-xxxxxx:tr-attach-xxxxxx:vbr-xxxxxx
```