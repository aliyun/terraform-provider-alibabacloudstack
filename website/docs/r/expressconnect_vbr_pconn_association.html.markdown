---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbrpconnassociation"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbrpconnassociation"
description: |-
  Provides a expressconnect Vbrpconnassociation resource.
---

# alibabacloudstack\_expressconnect\_vbrpconnassociation

Provides a expressconnect Vbrpconnassociation resource.

## Example Usage
```
variable "name" {
  default = "tf-testaccexpressconnect-vbrpconn1073"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id =     ""
	vlan_id =                    1073
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	enable_ipv6              = true
	local_ipv6_gateway_ip = "2408:4004:cc:400::1"
	peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
	peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
}



resource "alibabacloudstack_expressconnect_vbr_pconn_association" "default" {
  vlan_id = "1076"
  peering_subnet_mask = "255.255.255.0"
  physical_connection_id = ""
  local_gateway_ip = "10.100.0.1"
  local_ipv6_gateway_ip = "2408:4004:cc:500::1"
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  peering_ipv6_subnet_mask = "2408:4004:cc:500::/56"
  peer_gateway_ip = "10.100.0.2"
  enable_ipv6 = "true"
  peer_ipv6_gateway_ip = "2408:4004:cc:500::2"
}
```

## Argument Reference

The following arguments are supported:
  * `physical_connection_id` - (Required, ForceNew) - The ID of the leased line instance.
  * `vbr_id` - (Required, ForceNew) - The ID of the VBR instance.
  * `vlan_id` - (Required, ForceNew) - VLAN ID of the VBR. Valid values: **0 to 2999 * *.> only the owner of the physical connection can specify this parameter. The VLAN ID of two VBRs under the same physical connection cannot be the same.
  * `local_gateway_ip` - (Optional, ForceNew) - The Alibaba cloud IP address of the VBR instance.
  * `peer_gateway_ip` - (Optional, ForceNew) - The client IP address of the VBR instance.-This attribute only allows the VBR owner to specify or modify.-Required when creating a VBR instance for the physical connection owner.
  * `peering_subnet_mask` - (Optional, ForceNew) - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.The two IP addresses must be in the same subnet.
  * `enable_ipv6` - (Optional, ForceNew) - Whether IPv6 is enabled. Value:-**true**: on.-**false** (default): Off.
  * `local_ipv6_gateway_ip` - (Optional, ForceNew) - The IPv6 address on the Alibaba Cloud side of the VBR instance.
  * `peer_ipv6_gateway_ip` - (Optional, ForceNew) - The IPv6 address of the client side of the VBR instance.-This attribute only allows the VBR owner to specify or modify.-Required when creating a VBR instance for the physical connection owner.
  * `peering_ipv6_subnet_mask` - (Optional, ForceNew) - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.Two IPv6 addresses must be in the same subnet.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `circuit_code` - The circuit code provided by the operator for the physical connection.> Only the owner of the physical line can specify this parameter.
  * `status` - The status of the resource
