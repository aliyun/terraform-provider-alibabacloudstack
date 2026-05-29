---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_connect_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-connect-attachments"
description: |-
  Provides a list of CEN transit router connect attachments.
---

# alibabacloudstack\_cen\_transit\_router\_connect\_attachments

This data source provides a list of CEN transit router connect attachments.

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccrouter_connect_attachment33915"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "${var.name}"
  cen_instance_name = "${var.name}"
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

data "alibabacloudstack_cen_transit_router_connect_attachments" "example" {
  cen_id             = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
  transit_router_id  = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_id}"
}

output "first_attachment_id" {
  value = data.alibabacloudstack_cen_transit_router_connect_attachments.example.transitrouterattachments.0.transit_router_attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN instance.
* `transit_router_id` - (Required) The ID of the transit router.
* `ids` - (Optional) A list of transit router connect attachment IDs.
* `name_regex` - (Optional) A regex string to filter results by attachment name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of transit router connect attachment IDs.
* `transitrouterattachments` - A list of transit router connect attachments. Each element contains the following attributes:
  * `transit_router_id` - The ID of the transit router.
  * `transit_router_attachment_name` - The name of the transit router attachment.
  * `transit_router_attachment_id` - The ID of the transit router attachment.
  * `resource_id` - The ID of the resource.
  * `creation_time` - The creation time of the attachment.
  * `resource_type` - The type of the resource.
  * `resource_owner_id` - The ID of the resource owner.
  * `status` - The status of the attachment.
  * `protocol` - The protocol used by the attachment.
  * `transport_type` - The transport type of the attachment.
