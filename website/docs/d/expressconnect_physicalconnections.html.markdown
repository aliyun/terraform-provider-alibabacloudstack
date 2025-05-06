---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_physicalconnections"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-physicalconnections"
description: |- 
  Provides a list of expressconnect physicalconnections owned by an Alibabacloudstack account.
---

# alibabacloudstack_expressconnect_physicalconnections
-> **NOTE:** Alias name has: `alibabacloudstack_express_connect_physical_connections`

This data source provides a list of expressconnect physicalconnections in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_expressconnect_physicalconnections" "ids" {
  ids = ["pc-2345678"]
}

output "express_connect_physical_connection_id_1" {
  value = data.alibabacloudstack_expressconnect_physicalconnections.ids.connections.0.id
}

data "alibabacloudstack_expressconnect_physicalconnections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

output "express_connect_physical_connection_id_2" {
  value = data.alibabacloudstack_expressconnect_physicalconnections.nameRegex.connections.0.id
}
```

## Argument Reference

The following arguments are supported:

* `include_reservation_data` - (Optional, ForceNew) Specifies whether to include reservation data that is not yet in effect. Valid values: `true` or `false`.
* `ids` - (Optional, ForceNew) A list of Physical Connection IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Physical Connection name.
* `status` - (Optional, ForceNew) The status of the physical connection. Valid values:
  * `Initial`: In application.
  * `Approved`: Approved.
  * `Allocating`: Resources are being allocated.
  * `Allocated`: Access under construction.
  * `Confirmed`: Waiting for user confirmation.
  * `Enabled`: Activated.
  * `Rejected`: The application was rejected.
  * `Canceled`: Canceled.
  * `Allocation Failed`: Resource allocation failed.
  * `Terminated`: Terminated.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Physical Connection names.
* `connections` - A list of Express Connect Physical Connections. Each element contains the following attributes:
  * `access_point_id` - The ID of the physical connection point.
  * `ad_location` - The physical location where the physical connection access device is located.
  * `bandwidth` - The bandwidth of the physical connection. Unit: Gbps.
  * `business_status` - The commercial status of the physical connection. Valid values:
    * `Normal`: Already activated.
    * `FinancialLocked`: Arrears locked.
    * `SecurityLocked`: Locked for security reasons.
  * `circuit_code` - The circuit code provided by the operator for the physical connection.
  * `create_time` - When the physical connection was created. Time is expressed according to ISO8601 standard and UTC time is used. The format is: YYYY-MM-DDThh:mm:ssZ.
  * `description` - The description of the physical connection.
  * `enabled_time` - The opening time of the physical connection.
  * `end_time` - The expiration time of the physical connection.
  * `has_reservation_data` - Whether to include order data that is not in effect. Valid values: `true` or `false`.
  * `id` - The ID of the physical connection.
  * `line_operator` - Operators that provide access to physical lines. Valid values:
    * `CT`: China Telecom.
    * `CU`: China Unicom.
    * `CM`: China Mobile.
    * `CO`: China Other.
    * `Equinix`: Equinix.
    * `Other`: Other abroad.
  * `loa_status` - The state of LOA. Valid values:
    * `Applying`: LOA application.
    * `Accept`: LOA application passed.
    * `Available`: LOA is available.
    * `Rejected`: LOA application rejected.
    * `Completing`: The dedicated line is under construction.
    * `Complete`: The construction of the dedicated line is completed.
    * `Deleted`: LOA has been deleted.
  * `payment_type` - The payment type of the resource.
  * `peer_location` - The geographic location of the local data center.
  * `physical_connection_id` - The ID of the instance of the leased line.
  * `physical_connection_name` - The name of the physical connection.
  * `port_number` - The port number of the physical connection device.
  * `port_type` - Physical connection port type. Valid values:
    * `100Base-T`: 100 megabytes of electrical ports.
    * `1000Base-T`: Gigabit port.
    * `1000Base-LX`: Gigabit single-mode optical port (10km).
    * `10GBase-T`: 10 gigabyte electrical port.
    * `10GBase-LR`: 10 trillion single mode optical port (10km).
  * `redundant_physical_connection_id` - The ID of the redundant physical connection.
  * `reservation_active_time` - The effective time of the uneffective order.
  * `reservation_internet_charge_type` - The Payment type of the order that does not take effect. Value: `PayByBandwidth`, which means pay-by-bandwidth.
  * `spec` - Specifications of the physical connection. Unit: G means Gbps.
  * `status` - The status of the physical connection. See the valid values for the `status` argument.
  * `type` - The type of the physical connection. The default value is `VPC`.
  * `reservation_order_type` - The order type of the reservation.