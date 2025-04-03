---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackages"
sidebar_current: "docs-Alibabacloudstack-datasource-cbwp-commonbandwidthpackages"
description: |- 
  Provides a list of cbwp commonbandwidthpackages owned by an alibabacloudstack account.
---

# alibabacloudstack_cbwp_commonbandwidthpackages
-> **NOTE:** Alias name has: `alibabacloudstack_common_bandwidth_packages`

This data source provides a list of cbwp commonbandwidthpackages in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
resource "alibabacloudstack_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "tf-testAccCommonBandwidthPackage"
  description = "tf-testAcc-CommonBandwidthPackage"
}

data "alibabacloudstack_cbwp_commonbandwidthpackages" "foo" {
  name_regex = "^tf-testAcc.*"
  ids        = [alibabacloudstack_common_bandwidth_package.foo.id]
  resource_group_id = "your_resource_group_id"
}

output "common_bandwidth_packages" {
  value = data.alibabacloudstack_cbwp_commonbandwidthpackages.foo.packages
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Common Bandwidth Package IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by name.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which you want to move the resource. You can use resource groups to facilitate resource grouping and permission management for an Alibaba Cloud. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of Common Bandwidth Package names.
* `ids` - A list of Common Bandwidth Package IDs.
* `packages` - A list of Common Bandwidth Packages. Each element contains the following attributes:
  * `id` - The ID of the Common Bandwidth Package.
  * `bandwidth` - The maximum bandwidth of the Internet Shared Bandwidth instance. Unit: Mbit/s. Valid values: **1** to **1000**. Default value: **1**.
  * `status` - The status of the Common Bandwidth Package. Default value: **Available**.
  * `name` - The name of the Common Bandwidth Package.
  * `description` - The description of the Common Bandwidth Package instance. The description must be 2 to 256 characters in length and start with a letter. The description cannot start with `http://` or `https://`.
  * `business_status` - The business status of the Common Bandwidth Package instance. Values:
    - **Normal**: Normal.
    - **Financialized**: Arrears.
    - **Unactivated**: Not activated.
  * `isp` - The line type. Valid values:
    - **BGP**: All regions support BGP (Multi-ISP).
    - **BGP_PRO**: BGP (Multi-ISP) Pro lines are available in the China (Hong Kong), Singapore, Japan (Tokyo), Philippines (Manila), Malaysia (Kuala Lumpur), Indonesia (Jakarta), and Thailand (Bangkok) regions.
    - If you are allowed to use single-ISP bandwidth, you can also use one of the following values:
      - **ChinaTelecom**
      - **ChinaUnicom**
      - **ChinaMobile**
      - **ChinaTelecom_L2**
      - **ChinaUnicom_L2**
      - **ChinaMobile_L2**
    - If your services are deployed in China East 1 Finance, this parameter is required and you must set the value to **BGP_FinanceCloud**.
  * `creation_time` - The time when the Common Bandwidth Package was created.
  * `public_ip_addresses` - The public IP addresses that belong to the Common Bandwidth Package. Each element contains the following attributes:
    * `ip_address` - The address of the EIP.
    * `allocation_id` - The ID of the EIP instance.