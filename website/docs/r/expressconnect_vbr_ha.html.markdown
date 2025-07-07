---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_ha"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbr_ha"
description: |-
  Provides a expressconnect Virtualborderrouter Ha Group resource.
---

# alibabacloudstack_vbr_ha

Provides a expressconnect Virtualborderrouter Ha Failover Group resource.

## Example Usage

Basic Usage
```
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
	local_gateway_ip =           "10.0.1.1"
	peer_gateway_ip =            "10.0.1.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_2"
}

resource "alibabacloudstack_expressconnect_vbr_ha" "default" {
  name = "tf-testaccexpressconnect-vbrha1926"
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default1.id}"
  peer_vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default2.id}"
  description = "tf-testaccexpressconnect-vbrha1926"
}
```

## Argument Reference

The following arguments are supported:
  * `name` - (Required, ForceNew) - The name of the Ha Failover Groups。 It must be 2 to 128 characters in length and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter.
  * `description` - (Optional, ForceNew) - It must be 2 to 256 characters in length and start with a letter. It cannot start with http:// or https://. 
  * `vbr_id` - (Required, ForceNew) - VBR ID。
  * `peer_vbr_id` - (Required, ForceNew) - Failover Group VBR ID。

## Attributes Reference

In addition to all the parameters mentioned above, there is a read-only property with no output.