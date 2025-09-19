---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachment"
sidebar_current: "docs-alibabacloudstack-resource-cen-transit-router-connect-attachment"
description: |-
  Provides a Alibaba Cloud CEN transit router connect attachment resource.
---

# alibabacloudstack\_cen\_transit_router_connect_attachment

Provides a CEN transit router connect attachment resource.

> **Note**: When creating the alibabacloudstack_cen_transit_router_connect_attachment resource, it requires that there is already an available alibabacloudstack_cen_transit_router_vbr_attachment resource in the cen_instance.

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccrouter_connect_attachment33915"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id = "YourPhysicalConnectionId"
	vlan_id =                    1
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" {
    transit_router_attachment_name = "${var.name}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
    depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_attachment_name` - (Optional) The name of the transit router attachment.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `{cen_id}:{transit_router_id}:{attachment_id}:{vpc_id}`.
* `transit_router_attachment_id` - The ID of the transit router attachment.
* `resource_id` - The resource ID.
* `creation_time` - The creation time of the attachment.
* `resource_type` - The type of the resource.
* `resource_owner_id` - The ID of the resource owner.
* `status` - The status of the attachment.

## Import

CEN transit router VPC attachment can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_cen_transit_router_vpc_attachment.default cen-abc12345678900001:tr-abc12345678900001:tr-attach-abc12345678900001:vpc-abc12345678900001
```
