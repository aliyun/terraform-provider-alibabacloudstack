---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_physicalconnection"
sidebar_current: "docs-Alibabacloudstack-expressconnect-physicalconnection"
description: |- 
  Provides a expressconnect Physicalconnection resource.
---

# alibabacloudstack_expressconnect_physicalconnection
-> **NOTE:** Alias name has: `alibabacloudstack_express_connect_physical_connection`

Provides a expressconnect Physicalconnection resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
    default = "tf-testaccexpress_connectphysical_connection34886"
}

resource "alibabacloudstack_expressconnect_physicalconnection" "default" {
  device_name              = var.name
  access_point_id          = "ap-cn-hangzhou-jg-B"
  line_operator            = "CO"
  peer_location            = "XX街道"
  physical_connection_name = var.name
  type                     = "VPC"
  description              = "abcabc"
  port_type                = "1000Base-LX"
  bandwidth                = 10
}
```

## Argument Reference

The following arguments are supported:

* `access_point_id` - (Required, ForceNew) The ID of the physical connection point.
* `bandwidth` - (Optional) The bandwidth of the physical connection. Unit: Gbps. Default value is `10`.
* `circuit_code` - (Optional) The circuit code provided by the operator for the physical connection.
* `description` - (Optional) The description of the physical connection.
* `line_operator` - (Required) Operators that provide access to physical lines. Valid values:
  * CT: China Telecom
  * CU: China Unicom
  * CM: China Mobile
  * CO: Other Chinese
  * Equinix: Equinix
  * Other: Other Overseas.
* `peer_location` - (Required) The geographic location of the local data center.
* `device_name` - (Required) The name of the physical device.
* `physical_connection_name` - (Optional) The name of the physical connection.
* `port_type` - (Optional) The port type of the physical connection. Valid values:
  * 100Base-T: Fast Ethernet electrical port.
  * 1000Base-T: Gigabit Ethernet electrical port.
  * 1000Base-LX: Gigabit single-mode optical port (10km).
  * 10GBase-T: 10 Gigabit Ethernet electrical port.
  * 10GBase-LR: 10 Gigabit single-mode optical port (10km).
  * 40GBase-LR: 40 Gigabit single-mode optical port.
  * 100GBase-LR: 100 Gigabit single-mode optical port.
  
  **NOTE:** From v1.185.0+, the `40GBase-LR` and `100GBase-LR` values are valid. Set these values based on the water levels of background ports. For details about the water levels, contact the business manager.
* `redundant_physical_connection_id` - (Optional) The ID of the redundant physical connection.
* `status` - (Optional) The status of the physical connection. Valid values:
  * Initial: In application.
  * Approved: Approved.
  * Allocating: Resources are being allocated.
  * Allocated: Access under construction.
  * Confirmed: Waiting for user confirmation.
  * Enabled: Activated.
  * Rejected: The application was rejected.
  * Canceled: Canceled.
  * Allocation Failed: Resource allocation failed.
  * Terminated: Terminated.
* `type` - (Optional, ForceNew) The type of the physical connection. Default value is `VPC`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the physical connection.
* `bandwidth` - The bandwidth of the physical connection. Unit: Gbps.
* `status` - The status of the physical connection. Valid values:
  * Initial: In application.
  * Approved: Approved.
  * Allocating: Resources are being allocated.
  * Allocated: Access under construction.
  * Confirmed: Waiting for user confirmation.
  * Enabled: Activated.
  * Rejected: The application was rejected.
  * Canceled: Canceled.
  * Allocation Failed: Resource allocation failed.
  * Terminated: Terminated.
* `type` - The type of the physical connection. Default value is `VPC`.
```
