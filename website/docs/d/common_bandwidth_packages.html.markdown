---
subcategory: "Common Bandwidth Package"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_common_bandwidth_packages"
description: |-
  Provides a list of Common Bandwidth Packages owned by an Alibaba Cloud account.
---

# alibabacloudstack\_common\_bandwidth\_packages

This data source provides a list of Common Bandwidth Packages in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_common_bandwidth_packages" "example" {
  name_regex = "^my-CBWP"
}

output "first_cbwp_id" {
  value = data.alibabacloudstack_common_bandwidth_packages.example.packages.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Common Bandwidth Package IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Common Bandwidth Package name.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `output_file` - (Optional, Deprecated) The output file path. This field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the `local_file` provider instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Common Bandwidth Package IDs.
* `names` - A list of Common Bandwidth Package names.
* `packages` - A list of Common Bandwidth Packages. Each element contains the following attributes:
  * `id` - The ID of the Common Bandwidth Package.
  * `bandwidth` - The peak bandwidth of the Common Bandwidth Package. Unit: Mbps.
  * `status` - The status of the Common Bandwidth Package. Valid values: `Available`, `Modifying`.
  * `name` - The name of the Common Bandwidth Package.
  * `description` - The description of the Common Bandwidth Package.
  * `business_status` - The business status of the Common Bandwidth Package. Valid values: `Normal`, `FinancialLocked`.
  * `isp` - The line type of the Common Bandwidth Package. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`.
  * `creation_time` - The creation time of the Common Bandwidth Package.
  * `public_ip_addresses` - A list of public IP addresses associated with the Common Bandwidth Package. Each element contains:
    * `ip_address` - The public IP address.
    * `allocation_id` - The allocation ID of the EIP.
