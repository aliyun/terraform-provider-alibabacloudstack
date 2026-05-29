---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_pconn_association"
sidebar_current: "docs-Alibabacloudstack-resource-expressconnect-vbr-pconn-association"
description: |-
  Provides a expressconnect Vbrpconnassociation resource.
---

# alibabacloudstack\_expressconnect\_vbrpconnassociation

Provides a expressconnect Vbrpconnassociation resource.

## Example Usage

### Basic Association

```hcl
variable "name" {
  default = "tf-testaccexpressconnect-vbrpconn1073"
}

resource "alibabacloudstack_expressconnect_virtualborderrouter" "default" {
  physical_connection_id = alibabacloudstack_expressconnect_physicalconnection.default.id
  vlan_id                = 1073
  local_gateway_ip       = "10.0.0.1"
  peer_gateway_ip        = "10.0.0.2"
  peering_subnet_mask    = "255.255.255.252"
  name                   = var.name
}

resource "alibabacloudstack_expressconnect_vbr_pconn_association" "default" {
  physical_connection_id = alibabacloudstack_expressconnect_physicalconnection.default.id
  vbr_id                 = alibabacloudstack_expressconnect_virtualborderrouter.default.id
  vlan_id                = "1076"
  local_gateway_ip       = "10.100.0.1"
  peer_gateway_ip        = "10.100.0.2"
  peering_subnet_mask    = "255.255.255.0"
}
```

### Association with IPv6

```hcl
resource "alibabacloudstack_expressconnect_vbr_pconn_association" "ipv6" {
  physical_connection_id  = alibabacloudstack_expressconnect_physicalconnection.default.id
  vbr_id                  = alibabacloudstack_expressconnect_virtualborderrouter.default.id
  vlan_id                 = "1076"
  local_gateway_ip        = "10.100.0.1"
  peer_gateway_ip         = "10.100.0.2"
  peering_subnet_mask     = "255.255.255.0"
  enable_ipv6             = true
  local_ipv6_gateway_ip   = "2408:4004:cc:500::1"
  peer_ipv6_gateway_ip    = "2408:4004:cc:500::2"
  peering_ipv6_subnet_mask = "2408:4004:cc:500::/56"
}
```

## Argument Reference

The following arguments are supported:
  * `physical_connection_id` - (Required, ForceNew) - The ID of the leased line instance.
  * `vbr_id` - (Required, ForceNew) - The ID of the VBR instance.
  * `vlan_id` - (Required, ForceNew) - The VLAN ID of the VBR. Valid values: `0` to `2999`. Only the owner of the physical connection can specify this parameter. The VLAN IDs of two VBRs under the same physical connection cannot be the same.
  * `local_gateway_ip` - (Optional, ForceNew) - The Alibaba cloud IP address of the VBR instance.
  * `peer_gateway_ip` - (Optional, ForceNew) - The client-side IP address of the VBR instance. This attribute can only be specified or modified by the VBR owner. Required when creating a VBR instance for the physical connection owner.
  * `peering_subnet_mask` - (Optional, ForceNew) - The subnet mask for the Alibaba Cloud side and client side of the VBR instance. The two IP addresses must be in the same subnet.
  * `enable_ipv6` - (Optional, ForceNew) - Specifies whether to enable IPv6. Valid values: `true` (enabled), `false` (disabled, default).
  * `local_ipv6_gateway_ip` - (Optional, ForceNew) - The IPv6 address on the Alibaba Cloud side of the VBR instance.
  * `peer_ipv6_gateway_ip` - (Optional, ForceNew) - The IPv6 address on the client side of the VBR instance. This attribute can only be specified or modified by the VBR owner. Required when creating a VBR instance for the physical connection owner. Note: This parameter is required when `enable_ipv6` is set to `true` and must be used together with `local_ipv6_gateway_ip` and `peering_ipv6_subnet_mask`.
  * `peering_ipv6_subnet_mask` - (Optional, ForceNew) - The IPv6 subnet mask for the Alibaba Cloud side and client side of the VBR instance. The two IPv6 addresses must be in the same subnet. Note: This parameter is required when `enable_ipv6` is set to `true` and must be used together with `local_ipv6_gateway_ip` and `peer_ipv6_gateway_ip`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

  * `id` - The ID of the resource, formatted as `<physical_connection_id>:<vbr_id>`.
  * `circuit_code` - The circuit code provided by the operator for the physical connection. Only the owner of the physical connection can specify this parameter.
  * `status` - The status of the VBR-PCCN association. Valid values include `Creating`, `Associated`, `UnAssociating`.
  * `vlan_id` - The VLAN ID of the VBR.

## Import

VBR-PCCN Association can be imported using the `physical_connection_id` and `vbr_id` separated by a colon, e.g.

```
$ terraform import alibabacloudstack_expressconnect_vbr_pconn_association.example pc-abc12345:vbr-xyz67890
```
