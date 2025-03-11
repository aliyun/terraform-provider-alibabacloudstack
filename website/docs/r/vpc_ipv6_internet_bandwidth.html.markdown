---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_internet_bandwidth"
sidebar_current: "docs-alibabacloudstack-resource-vpc-ipv6-internet-bandwidth"
description: |-
  Provides a Alibabacloudstack VPC Ipv6 Internet Bandwidth resource.
---

# alibabacloudstack\_vpc\_ipv6\_internet\_bandwidth

Provides a VPC Ipv6 Internet Bandwidth resource.

For information about VPC Ipv6 Internet Bandwidth and how to use it, see [What is Ipv6 Internet Bandwidth](https://www.alibabacloud.com/help/doc-detail/102213.htm).



## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_instances" "example" {
  name_regex = "ecs_with_ipv6_address"
  status     = "Running"
}

data "alibabacloudstack_vpc_ipv6_addresses" "example" {
  associated_instance_id = data.alibabacloudstack_instances.example.instances.0.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_internet_bandwidth" "example" {
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.example.addresses.0.id
  ipv6_gateway_id      = data.alibabacloudstack_vpc_ipv6_addresses.example.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}

```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The amount of Internet bandwidth resources of the IPv6 address, Unit: `Mbit/s`. Valid values: `1` to `5000`. **NOTE:** If `internet_charge_type` is set to `PayByTraffic`, the amount of Internet bandwidth resources of the IPv6 address is limited by the specification of the IPv6 gateway. `Small` (default): specifies the Free edition and the Internet bandwidth is from `1` to `500` Mbit/s. `Medium`: specifies the Medium edition and the Internet bandwidth is from `1` to `1000` Mbit/s. `Large`: specifies the Large edition and the Internet bandwidth is from `1` to `2000` Mbit/s.
* `internet_charge_type` - (Optional, Computed, ForceNew) The metering method of the Internet bandwidth resources of the IPv6 gateway. Valid values: `PayByBandwidth`, `PayByTraffic`.
* `ipv6_address_id` - (Required, ForceNew) The ID of the IPv6 address.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipv6 Internet Bandwidth.
* `status` - The status of the resource.Valid values:`Normal`, `FinancialLocked` and `SecurityLocked`.

## Import

VPC Ipv6 Internet Bandwidth can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_vpc_ipv6_internet_bandwidth.example <id>
```