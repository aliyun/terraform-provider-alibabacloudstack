---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_internet_bandwidths"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-internet-bandwidths"
description: |- 
  Provides a list of vpc ipv6 internet bandwidths owned by an alibabacloudstack account.
---

# alibabacloudstack_vpc_ipv6_internet_bandwidths
-> **NOTE:** Alias name has: `alibabacloudstack_vpc_ipv6_internetbandwidths`

This data source provides a list of vpc ipv6 internet bandwidths in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "ids" {
  ids = ["example_id"]
}

output "vpc_ipv6_internet_bandwidth_id_1" {
  value = data.alibabacloudstack_vpc_ipv6_internet_bandwidths.ids.bandwidths.0.id
}

data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "ipv6InternetBandwidthId" {
  ipv6_internet_bandwidth_id = "example_value"
}

output "vpc_ipv6_internet_bandwidth_id_2" {
  value = data.alibabacloudstack_vpc_ipv6_internet_bandwidths.ipv6InternetBandwidthId.bandwidths.0.id
}

data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "ipv6AddressId" {
  ipv6_address_id = "example_value"
}

output "vpc_ipv6_internet_bandwidth_id_3" {
  value = data.alibabacloudstack_vpc_ipv6_internet_bandwidths.ipv6AddressId.bandwidths.0.id
}

data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "status" {
  status = "Normal"
}

output "vpc_ipv6_internet_bandwidth_id_4" {
  value = data.alibabacloudstack_vpc_ipv6_internet_bandwidths.status.bandwidths.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of IPv6 Internet Bandwidth IDs.
* `ipv6_internet_bandwidth_id` - (Optional, ForceNew) The ID of the IPv6 Internet Bandwidth.
* `ipv6_address_id` - (Optional, ForceNew) The ID of the IPv6 address instance.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Normal`, `FinancialLocked`, and `SecurityLocked`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of IPv6 Internet Bandwidth names.
* `bandwidths` - A list of VPC IPv6 Internet Bandwidths. Each element contains the following attributes:
  * `bandwidth` - The amount of Internet bandwidth resources of the IPv6 address, Unit: `Mbit/s`. Valid values: `1` to `5000`. **NOTE:** If `internet_charge_type` is set to `PayByTraffic`, the amount of Internet bandwidth resources of the IPv6 address is limited by the specification of the IPv6 gateway. `Small` (default): specifies the Free edition and the Internet bandwidth is from `1` to `500` Mbit/s. `Medium`: specifies the Medium edition and the Internet bandwidth is from `1` to `1000` Mbit/s. `Large`: specifies the Large edition and the Internet bandwidth is from `1` to `2000` Mbit/s.
  * `id` - The ID of the IPv6 Internet Bandwidth.
  * `internet_charge_type` - The metering method of the Internet bandwidth resources of the IPv6 gateway. Valid values: `PayByBandwidth`, `PayByTraffic`.
  * `ipv6_address_id` - The ID of the IPv6 address instance.
  * `ipv6_gateway_id` - The ID of the IPv6 gateway to which the IPv6 address belongs.
  * `ipv6_internet_bandwidth_id` - The ID of the IPv6 Internet Bandwidth.
  * `payment_type` - The payment type of the resource.
  * `status` - The status of the resource. Valid values: `Normal`, `FinancialLocked`, and `SecurityLocked`.