---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vswitches"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-vswitches"
description: |- 
  Provides a list of vpc vswitches owned by an alibabacloudstack account.
---

# alibabacloudstack_vpc_vswitches
-> **NOTE:** Alias name has: `alibabacloudstack_vswitches`

This data source provides a list of vpc vswitches in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpc_vswitches" "default" {
  name_regex = "^my-vswitch-"
  vpc_id     = alibabacloudstack_vpc.example.id
  zone_id    = "cn-hangzhou-b"
}

output "vswitches" {
  value = data.alibabacloudstack_vpc_vswitches.default.vswitches[*].vswitch_name
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) The CIDR block of the vSwitch. You can use this parameter to filter vSwitches with specific CIDR blocks.
* `name_regex` - (Optional) A regex string used to filter results by the name of the vSwitch.
* `is_default` - (Optional, Type: bool) Specifies whether to query the default vSwitches in the specified region. Valid values:
  * **true** - Query only default vSwitches.
  * **false** - Exclude default vSwitches from the query.
  If you do not set this parameter, the system queries all vSwitches in the specified region by default.
* `vpc_id` - (Optional) The ID of the virtual private cloud (VPC) to which the vSwitches belong. At least one of `vpc_id` or `zone_id` must be specified.
* `zone_id` - (Optional) The ID of the availability zone where the vSwitches are located. You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/36064.html) operation to query the most recent zone list.
* `ids` - (Optional) A list of vSwitch IDs to filter the results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of vSwitch IDs.
* `names` - A list of vSwitch names.
* `vswitches` - A list of vSwitch details. Each element contains the following attributes:
  * `id` - The ID of the vSwitch.
  * `vpc_id` - The ID of the VPC to which the vSwitch belongs.
  * `zone_id` - The ID of the availability zone where the vSwitch is located.
  * `vswitch_name` - The name of the vSwitch.
  * `instance_ids` - A list of ECS instance IDs that are associated with the vSwitch.
  * `cidr_block` - The CIDR block of the vSwitch.
  * `description` - The description of the vSwitch.
  * `is_default` - Indicates whether the vSwitch is the default one in the region.
  * `creation_time` - The time when the vSwitch was created.
  * `available_ip_address_count` - The number of available IP addresses in the vSwitch.