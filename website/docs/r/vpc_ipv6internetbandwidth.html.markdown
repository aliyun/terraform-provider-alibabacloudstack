---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6internetbandwidth"
sidebar_current: "docs-Alibabacloudstack-vpc-ipv6internetbandwidth"
description: |- 
  Provides a vpc Ipv6Internetbandwidth resource.
---

# alibabacloudstack_vpc_ipv6internetbandwidth
-> **NOTE:** Alias name has: `alibabacloudstack_vpc_ipv6_internet_bandwidth`

Provides a vpc Ipv6Internetbandwidth resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccvpcipv6internetbandwidth46716"
}

data "alibabacloudstack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alibabacloudstack_instances.default.instances.0.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The amount of Internet bandwidth resources of the IPv6 address, Unit: `Mbit/s`. Valid values: `1` to `5000`. **NOTE:** If `internet_charge_type` is set to `PayByTraffic`, the amount of Internet bandwidth resources of the IPv6 address is limited by the specification of the IPv6 gateway:
  * `Small` (default): specifies the Free edition and the Internet bandwidth is from `1` to `500` Mbit/s.
  * `Medium`: specifies the Medium edition and the Internet bandwidth is from `1` to `1000` Mbit/s.
  * `Large`: specifies the Large edition and the Internet bandwidth is from `1` to `2000` Mbit/s.
* `internet_charge_type` - (Optional, ForceNew) The metering method of the Internet bandwidth resources of the IPv6 gateway. Valid values: `PayByBandwidth`, `PayByTraffic`.
* `ipv6_address_id` - (Required, ForceNew) The ID of the IPv6 address instance.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway to which the IPv6 address belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in Terraform of Ipv6 Internet Bandwidth.
* `status` - The status of the resource. Valid values: `Normal`, `FinancialLocked`, and `SecurityLocked`.
* `internet_charge_type` - The metering method of the Internet bandwidth resources of the IPv6 gateway. Valid values: `PayByBandwidth`, `PayByTraffic`.

### Explanation of Changes

1. **Example Usage**: The example has been updated to include variable definitions for better clarity and reusability. It also ensures that the `ipv6_address_id` and `ipv6_gateway_id` are correctly referenced from the `data` sources.

2. **Argument Reference**:
   - Added detailed descriptions for each argument, ensuring clarity about their purpose and valid values.
   - Emphasized the limitations when using `PayByTraffic` with different editions (`Small`, `Medium`, `Large`).

3. **Attributes Reference**:
   - Included additional attributes such as `id` and `status` for completeness.
   - Ensured that the `internet_charge_type` attribute is explicitly mentioned as an exported attribute for consistency.

This updated documentation provides a more comprehensive and clear guide for users to understand and utilize the `alibabacloudstack_vpc_ipv6internetbandwidth` resource effectively.