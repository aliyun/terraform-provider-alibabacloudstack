---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_vbr_health_checks"
sidebar_current: "docs-alibabacloudstack-datasource-cen-vbr-health-checks"
description: |-
  Provides a list of CEN VBR Health Checks to the user.
---

# alibabacloudstack\_cen\_vbr\_health\_checks

This data source provides a list of CEN VBR Health Checks in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
variable "name" {
  default = "tf-testAccCenVbrHealthChecks"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc-123456"
  virtual_border_router_name = var.name
  vlan_id                    = 1991
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
  vbr_id            = alibabacloudstack_express_connect_virtual_border_router.default.id
  cen_id            = alibabacloudstack_cen_instance.default.id
  transit_router_id = alibabacloudstack_cen_instance.default.transit_router_id
}

resource "alibabacloudstack_cen_vbr_health_check" "default" {
  cen_id                 = alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id
  vbr_instance_id        = alibabacloudstack_express_connect_virtual_border_router.default.id
  health_check_source_ip = "172.16.0.1"
  health_check_target_ip = "10.0.0.1"
  healthy_threshold      = 3
  health_check_interval  = 2
  health_check_only      = true
  depends_on             = [alibabacloudstack_cen_transit_router_vbr_attachment.default]
}

data "alibabacloudstack_cen_vbr_health_checks" "default" {
  cen_id          = alibabacloudstack_cen_vbr_health_check.default.cen_id
  vbr_instance_id = alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id
  ids             = ["${alibabacloudstack_cen_vbr_health_check.default.cen_id}:${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"]
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Optional, ForceNew) The ID of the CEN instance.
* `vbr_instance_id` - (Optional, ForceNew) The ID of the VBR instance.
* `ids` - (Optional, ForceNew) A list of VBR Health Check IDs in the format of `<cen_id>:<vbr_instance_id>`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of VBR Health Check IDs.
* `vbr_health_checks` - A list of VBR Health Checks. Each element contains the following attributes:
  * `cen_id` - The ID of the CEN instance.
  * `vbr_instance_id` - The ID of the VBR instance.
  * `vbr_instance_region_id` - The region ID of the VBR instance.
  * `health_check_source_ip` - The source IP address of the health check.
  * `health_check_target_ip` - The destination IP address of the health check.
  * `health_check_interval` - The interval of the health check.
  * `healthy_threshold` - The number of consecutive health check successes before the link is considered healthy.
  * `health_check_only` - Whether to perform health checks only.
```
