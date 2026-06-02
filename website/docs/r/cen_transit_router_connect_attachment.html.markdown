---
subcategory: "Cloud Enterprise Network"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachment"
description: |-
  Provides a CEN transit router connect attachment resource.
---

# alibabacloudstack_cen_transit_router_connect_attachment

Provides a CEN transit router connect attachment resource.

> **NOTE**: When creating an alibabacloudstack_cen_transit_router_connect_attachment resource, it requires that an available alibabacloudstack_cen_transit_router_vbr_attachment resource already exists in the cen_instance.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance. You can obtain the ID after creating a CEN instance.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router instance. You can obtain the ID after creating a transit router.
* `transit_router_attachment_name` - (Optional, ForceNew) The name of the network instance connection. The name must be 1 to 128 characters in length, and can contain Chinese, English, digits, underscores (_), and hyphens (-). It cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, with the format `{CenId:TransitRouterId:TransitRouterAttachmentId}`.
* `creation_time` - The creation time of the network instance connection. The time is formatted in ISO 8601 standard using UTC time. Format: YYYY-MM-DDThh:mmZ.
* `protocol` - The protocol type. Possible values include: BGP, Static, etc.
* `resource_id` - The ID of the network instance associated with the network instance connection.
* `resource_owner_id` - The ID of the account to which the network instance belongs.
* `resource_type` - The type of the network instance associated with the network instance connection.
* `status` - The status of the network instance connection.
* `transit_router_attachment_id` - The ID of the network instance connection.
* `transport_type` - The transport type.

## Import

CEN transit router connect attachment can be imported using the ID, e.g.

```shell
$ terraform import alibabacloudstack_cen_transit_router_connect_attachment.default cen-abc12345678900001:tr-abc12345678900001:tr-attach-abc12345678900001
```