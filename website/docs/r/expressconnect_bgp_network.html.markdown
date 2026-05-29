---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_network"
sidebar_current: "docs-Alibabacloudstack-resource-expressconnect-bgp-network"
description: |-
  Provides a Express Connect BGP Network resource.
---

# alibabacloudstack\_expressconnect\_bgp\_network

Provides a Express Connect BGP Network resource.

## Example Usage

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
    description = "BGP Network Test"
}

resource "alibabacloudstack_expressconnect_bgp_network" "default" {
  dst_cidr_block = "10.10.0.0/24"
  router_id = alibabacloudstack_express_connect_virtual_border_router.default.id
}
```

## Argument Reference

The following arguments are supported:

* `dst_cidr_block` - (Required) The CIDR block of the VPC or vSwitch that you want to connect to a data center.
* `router_id` - (Required, ForceNew) The ID of the router (VBR) associated with the router interface. Modifying this parameter will force a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the BGP Network. The format is `<dst_cidr_block>:<router_id>`.
* `status` - The status of the BGP Network.

## Import

Express Connect BGP Network can be imported using the ID, e.g.

```
$ terraform import alibabacloudstack_expressconnect_bgp_network.example <dst_cidr_block>:<router_id>
```
