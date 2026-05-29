---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_ha"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbr_ha"
description: |-
  Provides a Express Connect VBR HA resource.
---

# alibabacloudstack_expressconnect_vbr_ha

Provides a Express Connect VBR HA (VBR Failover Group) resource.

-> **Note:** This resource can also be referred to by the following alias:
-> - `alibabacloudstack_vbr_ha`

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf-testaccexpressconnect-vbrha1926"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default1" {
	physical_connection_id =     ""
	vlan_id =                    1926
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_1"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default2" {
	physical_connection_id =     ""
	vlan_id =                    1929
	local_gateway_ip =           "10.1.0.1"
	peer_gateway_ip =            "10.1.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_2"
}

resource "alibabacloudstack_expressconnect_vbr_ha" "default" {
  name        = var.name
  vbr_id      = alibabacloudstack_express_connect_virtual_border_router.default1.id
  peer_vbr_id = alibabacloudstack_express_connect_virtual_border_router.default2.id
  description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the VBR failover group.
* `vbr_id` - (Required, ForceNew) The ID of the VBR.
* `peer_vbr_id` - (Required, ForceNew) The ID of the other VBR in the VBR failover group.
* `description` - (Optional, ForceNew) The description of the VBR failover group. It must be 2 to 256 characters in length, must start with a letter or Chinese, but cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VBR failover group.

## Import

Express Connect VBR HA can be imported using the VBR HA ID (e.g. `vbrha-xxxxxxxxx`), e.g.

```
$ terraform import alibabacloudstack_expressconnect_vbr_ha.example vbrha-xxxxxxxxx
```