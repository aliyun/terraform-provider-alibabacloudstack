---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_virtualborderrouters"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-virtualborderrouters"
description: |- 
  Provides a list of expressconnect virtualborderrouters owned by an alibabacloudstack account.
---

# alibabacloudstack_expressconnect_virtualborderrouters
-> **NOTE:** Alias name has: `alibabacloudstack_express_connect_virtual_border_routers`

This data source provides a list of Express Connect Virtual Border Routers (VBRs) in an Alibaba Cloud Stack account according to the specified filters.

## Example Usage

Basic Usage:

```terraform
data "alibabacloudstack_expressconnect_virtualborderrouters" "ids" {
}

output "express_connect_virtual_border_router_id_1" {
  value = data.alibabacloudstack_expressconnect_virtualborderrouters.ids.routers.0.id
}

data "alibabacloudstack_expressconnect_virtualborderrouters" "nameRegex" {
  name_regex = "^my-VirtualBorderRouter"
}

output "express_connect_virtual_border_router_id_2" {
  value = data.alibabacloudstack_expressconnect_virtualborderrouters.nameRegex.routers.0.id
}

data "alibabacloudstack_expressconnect_virtualborderrouters" "filter" {
  filter {
    key    = "PhysicalConnectionId"
    values = ["pc-xxxx1"]
  }
  filter {
    key    = "VbrId"
    values = ["vbr-xxxx1", "vbr-xxxx2"]
  }
}

output "express_connect_virtual_border_router_id_3" {
  value = data.alibabacloudstack_expressconnect_virtualborderrouters.filter.routers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) Custom filter block as described below:
  * `key` - (Required) The key of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
* `ids` - (Optional, ForceNew) A list of Virtual Border Router IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Virtual Border Router name.
* `status` - (Optional, ForceNew) The instance state with. Valid values: `active`, `deleting`, `recovering`, `terminated`, `terminating`, `unconfirmed`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Virtual Border Router names.
* `routers` - A list of Express Connect Virtual Border Routers. Each element contains the following attributes:
  * `access_point_id` - The ID of the physical connection point.
  * `activation_time` - The time when the VBR is first activated.
  * `circuit_code` - The circuit code provided by the operator for the physical connection.
  * `cloud_box_instance_id` - The ID of the cloud box instance associated with the VBR.
  * `create_time` - The creation time of the VBR.
  * `description` - The description information of the VBR.
  * `detect_multiplier` - Multiple of detection time. That is, the maximum number of connection packet losses allowed by the receiver to send messages, which is used to detect whether the link is normal. Valid values: **3 to 10**.
  * `ecc_id` - The ID of the high-speed cloud service instance.
  * `enable_ipv6` - Whether IPv6 is enabled:
    - **true**: Enabled.
    - **false**: Disabled.
  * `id` - The ID of the Virtual Border Router.
  * `local_gateway_ip` - The IPv4 address on the Alibaba Cloud side of the VBR instance.
  * `local_ipv6_gateway_ip` - The IPv6 address on the Alibaba Cloud side of the VBR instance.
  * `min_rx_interval` - Configure the receiving interval of BFD packets. Values: **200 to 1000**, in ms.
  * `min_tx_interval` - Configure the sending interval of BFD packets. Value: **200~1000**, unit: ms.
  * `payment_vbr_expire_time` - The billing expiration time of the VBR.
  * `peer_gateway_ip` - The IPv4 address of the client side of the VBR instance.
  * `peer_ipv6_gateway_ip` - The IPv6 address of the client side of the VBR instance.
  * `peering_ipv6_subnet_mask` - The subnet masks of the Alibaba Cloud-side IPv6 and the customer-side IPv6 of the VBR instance.
  * `peering_subnet_mask` - The subnet masks of the Alibaba Cloud-side IPv4 and the customer-side IPv4 of the VBR instance.
  * `physical_connection_business_status` - The business status of the physical connection:
    - **Normal**: Normal.
    - **Financialized**: Arrears locked.
  * `physical_connection_id` - The ID of the physical connection to which the VBR belongs.
  * `physical_connection_owner_uid` - The ID of the account to which the physical connection belongs.
  * `physical_connection_status` - The status of the physical connection:
    - **Initial**: The application is in progress.
    - **Approved**: Approved.
    - **Allocating**: Resources are being allocated.
    - **Allocated**: Access is under construction.
    - **Confirmed**: Waiting for user confirmation.
    - **Enabled**: Activated.
    - **Rejected**: The application was rejected.
    - **Canceled**: Canceled.
    - **Allocation Failed**: Resource allocation failed.
    - **Terminated**: Terminated.
  * `recovery_time` - The last time VBR returned from the Terminated state to the Active state.
  * `route_table_id` - The ID of the route table of the VBR.
  * `status` - The status of the resource.
  * `termination_time` - The time when VBR was last terminated.
  * `type` - VBR type.
  * `virtual_border_router_id` - The ID of the VBR.
  * `virtual_border_router_name` - The name of the VBR instance.
  * `vlan_id` - The VLAN ID of the VBR instance.
  * `vlan_interface_id` - The ID of the VBR router interface.