---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_gateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-gateways"
description: |- 
  Provides a list of vpc ipv6 gateways owned by an Alibabacloudstack account.
---

# alibabacloudstack_vpc_ipv6_gateways

This data source provides a list of vpc ipv6 gateways in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_vpc_ipv6_gateways" "ids" {
  ids = ["example_id"]
}
output "vpc_ipv6_gateway_id_1" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.ids.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "nameRegex" {
  name_regex = "^my-Ipv6Gateway"
}
output "vpc_ipv6_gateway_id_2" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.nameRegex.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "vpcId" {
  ids    = ["example_id"]
  vpc_id = "example_value"
}
output "vpc_ipv6_gateway_id_3" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.vpcId.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "status" {
  ids    = ["example_id"]
  status = "Available"
}
output "vpc_ipv6_gateway_id_4" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.status.gateways.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of IPv6 Gateway IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by IPv6 Gateway name.
* `ipv6_gateway_name` - (Optional, ForceNew) The name of the IPv6 gateway. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Pending`, and `Deleting`.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) to which the IPv6 gateway belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `gateways` - A list of VPC IPv6 Gateways. Each element contains the following attributes:
  * `business_status` - The status of the IPv6 gateway. Valid values: `Normal`, `FinancialLocked`, and `SecurityLocked`. 
    - `Normal`: Working as expected.
    - `FinancialLocked`: Locked due to overdue payments.
    - `SecurityLocked`: Locked due to security reasons.
  * `create_time` - The creation time of the resource.
  * `description` - The description of the IPv6 gateway. The description must be 2 to 256 characters in length. It cannot start with `http://` or `https://`.
  * `expired_time` - The expiration time of the IPv6 gateway.
  * `instance_charge_type` - The charge type of the IPv6 gateway. Valid value: `PayAsYouGo`.
  * `id` - The unique identifier (ID) of the IPv6 Gateway.
  * `ipv6_gateway_id` - The primary key attribute field for the IPv6 Gateway.
  * `ipv6_gateway_name` - The name of the IPv6 gateway. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
  * `spec` - The specification of the IPv6 gateway. This parameter is no longer used as IPv6 gateways do not distinguish between specifications.
  * `status` - The status of the IPv6 gateway. Valid values: `Available`, `Pending`, and `Deleting`.
  * `vpc_id` - The ID of the virtual private cloud (VPC) to which the IPv6 gateway belongs.