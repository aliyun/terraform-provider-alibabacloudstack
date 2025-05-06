---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_virtualborderrouter"
sidebar_current: "docs-Alibabacloudstack-expressconnect-virtualborderrouter"
description: |- 
  Provides a expressconnect Virtualborderrouter resource.
---

# alibabacloudstack_expressconnect_virtualborderrouter
-> **NOTE:** Alias name has: `alibabacloudstack_express_connect_virtual_border_router`

Provides a expressconnect Virtualborderrouter resource.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_express_connect_physical_connections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

resource "alibabacloudstack_expressconnect_virtualborderrouter" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alibabacloudstack_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = "example_value"
  vlan_id                    = 1
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
```

## Argument Reference

The following arguments are supported:

* `associated_physical_connections` - (Optional) The associated physical connection information.
* `bandwidth` - (Optional) The bandwidth of the VBR instance.
* `circuit_code` - (Optional) The circuit code provided by the operator for the physical connection.
* `description` - (Optional) The description information of the VBR. The length is from 2 to 256 characters, must start with a letter or Chinese character, and cannot start with `http://` or `https://`.
* `detect_multiplier` - (Optional) Multiple of detection time. That is, the maximum number of connection packet losses allowed by the receiver to send messages, which is used to detect whether the link is normal. Valid values: **3 to 10**.
* `enable_ipv6` - (Optional) Whether IPv6 is enabled. Valid values:
  - `true`: Enabled.
  - `false`: Disabled.
* `local_gateway_ip` - (Required) The IPv4 address on the Alibaba Cloud side of the VBR instance.
* `local_ipv6_gateway_ip` - (Optional) The IPv6 address on the Alibaba Cloud side of the VBR instance.
* `min_rx_interval` - (Optional) Configure the receiving interval of BFD packets. Values: **200 to 1000**, in ms.
* `min_tx_interval` - (Optional) Configure the sending interval of BFD packets. Value: **200~1000**, unit: ms.
* `peer_gateway_ip` - (Required) The IPv4 address of the client side of the VBR instance.
* `peer_ipv6_gateway_ip` - (Optional) The IPv6 address of the client side of the VBR instance.
* `peering_ipv6_subnet_mask` - (Optional) The subnet masks of the Alibaba Cloud-side IPv6 and the customer-side IPv6 of the VBR instance.
* `peering_subnet_mask` - (Required) The subnet masks of the Alibaba Cloud-side IPv4 and the customer-side IPv4 of the VBR instance.
* `physical_connection_id` - (Required, ForceNew) The ID of the physical connection to which the VBR belongs.
* `status` - (Optional) The status of the resource. Valid values:
  - `active`
  - `deleting`
  - `recovering`
  - `terminated`
  - `terminating`
  - `unconfirmed`
* `vbr_owner_id` - (Optional) The ID of the VBR instance owner. The default is the ID of the Alibaba Cloud account.
* `virtual_border_router_name` - (Optional) The name of the VBR instance. The length is from 2 to 128 characters, must start with a letter or Chinese character, can contain numbers, underscores (`_`), and dashes (`-`), but cannot start with `http://` or `https://`.
* `vlan_id` - (Required) The VLAN ID of the VBR instance. Valid range: **0 to 2999**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in Terraform of the Virtual Border Router.
* `route_table_id` - The ID of the route table of the VBR.
* `detect_multiplier` - Multiple of detection time. That is, the maximum number of connection packet losses allowed by the receiver to send messages, which is used to detect whether the link is normal. Valid values: **3 to 10**.
* `enable_ipv6` - Whether IPv6 is enabled. Valid values:
  - `true`: Enabled.
  - `false`: Disabled.
* `min_rx_interval` - Configure the receiving interval of BFD packets. Values: **200 to 1000**, in ms.
* `min_tx_interval` - Configure the sending interval of BFD packets. Value: **200~1000**, unit: ms.
* `status` - The status of the resource.
* `route_table_id` - The ID of the route table of the VBR.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 2 mins) Used when updating the Virtual Border Router.

## Import

Express Connect Virtual Border Router can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_expressconnect_virtualborderrouter.example <id>
```