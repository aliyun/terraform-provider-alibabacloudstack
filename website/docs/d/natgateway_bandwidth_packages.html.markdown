---
subcategory: "Network Address Translation Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_bandwidth_packages"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-bandwidth-packages"
description: |-
  Provides a list of NAT Gateway Bandwidth Packages owned by an Alibaba Cloud account.
---

# alibabacloudstack\_natgateway\_bandwidth\_packages

This data source provides a list of NAT Gateway Bandwidth Packages in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_natgateway_bandwidth_packages" "example" {
  name_regex = "^my-BWP"
}

output "first_bwp_id" {
  value = data.alibabacloudstack_natgateway_bandwidth_packages.example.bandwidth_packages.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of NAT Gateway Bandwidth Package IDs.
* `name_regex` - (Optional) A regex string to filter results by bandwidth package name.
* `description_regex` - (Optional) A regex string to filter results by bandwidth package description.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of NAT Gateway Bandwidth Package IDs.
* `bandwidth_packages` - A list of NAT Gateway Bandwidth Packages. Each element contains the following attributes:
  * `id` - The ID of the NAT Gateway Bandwidth Package.
  * `bandwidth` - The peak bandwidth of the bandwidth package. Unit: Mbps.
  * `bandwidth_package_id` - The ID of the bandwidth package.
  * `business_status` - The business status of the bandwidth package. Valid values: `Normal`, `FinancialLocked`, `Unactivated`.
  * `creation_time` - The creation time of the bandwidth package.
  * `name` - The name of the bandwidth package.
  * `description` - The description of the bandwidth package.
  * `instance_charge_type` - The billing method of the bandwidth package. Valid values: `PostPaid`, `PrePaid`.
  * `internet_charge_type` - The internet charge type. Valid values: `PayByBandwidth`, `PayBy95`.
  * `natgateway_id` - The ID of the NAT Gateway associated with the bandwidth package.
  * `ip_count` - The number of EIPs bound to the bandwidth package.
  * `isp` - The line type. Valid values: `BGP`, `BGP_PRO`.
  * `public_ip_addresses` - A list of public IP addresses. Each element contains:
    * `ip_address` - The public IP address.
    * `allocation_id` - The allocation ID of the EIP.
  * `status` - The status of the bandwidth package. Valid values: `Available`.
