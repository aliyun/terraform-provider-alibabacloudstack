---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_pconn_associations"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-vbrpconnassociations"
description: |-
  Provides a list of expressconnect vbrpconnassociations owned by an alibabacloudstack account.
---

# alibabacloudstack\_expressconnect\_vbrpconnassociations

This data source provides a list of expressconnect vbrpconnassociations in an alibabacloudstack account according to the specified filters.

## Example Usage
```
data "alibabacloudstack_expressconnect_vbr_pconn_associations" "example" {
  vbr_id = "vbr-bp1d8yixxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:
  * `vbr_id` - (Required) - The ID of the VBR instance.
  * `physical_connection_id` - (Optional) - The ID of the leased line instance.
  * `vlan_id` - (Optional) - VLAN ID of the VBR. Valid values: **0 to 2999 * *.> only the owner of the physical connection can specify this parameter. The VLAN ID of two VBRs under the same physical connection cannot be the same.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `vbr_pconn_associations` - The list of VBRs associated with the VBR.
    * `circuit_code` - The circuit code provided by the operator for the physical connection.> Only the owner of the physical line can specify this parameter.
    * `enable_ipv6` - Whether IPv6 is enabled. Value:-**true**: on.-**false** (default): Off.
    * `local_gateway_ip` - The Alibaba cloud IP address of the VBR instance.
    * `local_ipv6_gateway_ip` - The IPv6 address on the Alibaba Cloud side of the VBR instance.
    * `peer_gateway_ip` - The client IP address of the VBR instance.-This attribute only allows the VBR owner to specify or modify.-Required when creating a VBR instance for the physical connection owner.
    * `peer_ipv6_gateway_ip` - The IPv6 address of the client side of the VBR instance.-This attribute only allows the VBR owner to specify or modify.-Required when creating a VBR instance for the physical connection owner.
    * `peering_ipv6_subnet_mask` - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.Two IPv6 addresses must be in the same subnet.
    * `peering_subnet_mask` - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.The two IP addresses must be in the same subnet.
    * `physical_connection_id` - The ID of the leased line instance.
    * `status` - The status of the resource
    * `vbr_id` - The ID of the VBR instance.
    * `vlan_id` - VLAN ID of the VBR. Valid values: **0 to 2999 * *.> only the owner of the physical connection can specify this parameter. The VLAN ID of two VBRs under the same physical connection cannot be the same.
