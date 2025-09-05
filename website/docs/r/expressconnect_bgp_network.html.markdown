---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgpnetwork"
sidebar_current: "docs-Alibabacloudstack-expressconnect-bgpnetwork"
description: |-
  Provides a expressconnect Bgpnetwork resource.
---

# alibabacloudstack\_expressconnect\_bgpnetwork

Provides a expressconnect Bgpnetwork resource.

## Example Usage
```
variable "name" {
  default = "tf-testaccexpressconnect-bgp-group1106"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
    physical_connection_id = 
    vlan_id = 1106
    local_gateway_ip = "10.0.0.1"
    peer_gateway_ip = "10.0.0.2"
    peering_subnet_mask = "255.255.255.252"
    virtual_border_router_name = "${var.name}"
    description = "TestAccAlibabacloudStackExpressconnectBgpNetwork_basic0"
}



resource "alibabacloudstack_expressconnect_bgp_network" "default" {
  dst_cidr_block = "1.1.1.1"
  router_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `dst_cidr_block` - (Require) - expressconnect Bgpnetwork cidr block。
  * `router_id` - (Require, ForceNew) - vbr router id。

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `status` - expressconnect Bgpnetwork resource status。
