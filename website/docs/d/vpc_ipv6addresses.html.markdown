---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_addresses"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-addresses"
description: |- 
  Provides a list of vpc ipv6 addresses owned by an Alibabacloudstack account.
---

# alibabacloudstack_vpc_ipv6_addresses

This data source provides a list of vpc ipv6addresses in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_vpc_ipv6_addresses" "example" {
  associated_instance_id = "your-associated-instance-id"
  vswitch_id            = "your-vswitch-id"
  vpc_id               = "your-vpc-id"
  status               = "Available"

  output_file = "output.txt"
}

output "ipv6_address_1" {
  value = data.alibabacloudstack_vpc_ipv6_addresses.example.addresses[0].ipv6_address
}

output "ipv6_address_name_1" {
  value = data.alibabacloudstack_vpc_ipv6_addresses.example.addresses[0].ipv6_address_name
}
```

## Argument Reference

The following arguments are supported:

* `associated_instance_id` - (Optional, ForceNew) The ID of the instance that is assigned the IPv6 address.
* `ids` - (Optional, ForceNew) A list of IPv6 addresses IDs.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Pending`, and `Deleting`.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch to which the IPv6 address belongs.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC to which the IPv6 address belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `addresses` - A list of VPC IPv6 addresses. Each element contains the following attributes:
  * `associated_instance_id` - The ID of the instance that is assigned the IPv6 address.
  * `associated_instance_type` - The type of the instance that is assigned the IPv6 address.
  * `create_time` - The creation time of the resource.
  * `id` - The ID of the IPv6 Address.
  * `ipv6_address` - The IPv6 address.
  * `ipv6_address_id` - Resource primary key attribute field.
  * `ipv6_address_name` - The name of the IPv6 Address. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
  * `ipv6_gateway_id` - The ID of the IPv6 gateway to which the IPv6 address belongs.
  * `network_type` - The type of communication supported by the IPv6 address. Valid values: `Private` or `Public`. 
    - `Private`: Communication within the private network.
    - `Public`: Communication over the public network.
  * `status` - The status of the resource. Valid values: `Available`, `Pending`, and `Deleting`.
  * `vswitch_id` - The ID of the vSwitch to which the IPv6 address belongs.
  * `vpc_id` - The ID of the VPC to which the IPv6 address belongs.

* `names` - A list of IPv6 Address names.