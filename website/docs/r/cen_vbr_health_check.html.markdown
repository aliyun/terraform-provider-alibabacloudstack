---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_vbr_health_check"
sidebar_current: "docs-alibabacloudstack-resource-cen-vbr-health-check"
description: |-
  Provides a AlibabacloudStack CEN VBR Health Check resource.
---

# alibabacloudstack\_cen\_vbr\_health\_check

Provides a CEN VBR Health Check resource to configure health check for CEN VBR.

For information about CEN VBR Health Check and how to use it, see [What is VBR Health Check](https://www.alibabacloud.com/help/doc-detail/65883.htm).

## Example Usage

```terraform
variable "name" {
  default = "tf-testAcc9527"
}


resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc123456"
  virtual_border_router_name = var.name
  vlan_id                    = 1991
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_vbr_health_check" "default" {
    cen_id                 = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
    vbr_instance_id        = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
    health_check_source_ip = "172.16.0.1"
    health_check_target_ip = "10.0.0.1"
    healthy_threshold      = 3
    health_check_interval  = 2
    health_check_only      = true
    depends_on             = [alibabacloudstack_cen_transit_router_vbr_attachment.default]
}

```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `health_check_source_ip` - (Required, ForceNew) The source IP address of the health check.
* `health_check_target_ip` - (Required, ForceNew) The destination IP address of the health check.
* `health_check_interval` - (Required, ForceNew) The interval of the health check. Valid values: 1 to 3.
* `healthy_threshold` - (Required, ForceNew) The number of consecutive health check successes before the link is considered healthy. Valid values: 1 to 8.
* `health_check_only` - (Required, ForceNew) Whether to perform health checks only. Valid values: `true` and `false`.
* `vbr_instance_id` - (Required, ForceNew) The ID of the VBR instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<cen_id>:<vbr_instance_id>`.
* `vbr_instance_owner_id` - The owner ID of the VBR instance.
* `link_status` - The status of the link.
* `delay` - The latency of the health check.
* `packet_loss` - The packet loss rate of the health check.

## Import

CEN VBR Health Check can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cen_vbr_health_check.default cen-abc123456:vbr-def789012
```