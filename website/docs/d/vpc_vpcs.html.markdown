---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vpcs"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-vpcs"
description: |- 
  Provides a list of vpc vpcs owned by an alibabacloudstack account.
---

# alibabacloudstack_vpc_vpcs
-> **NOTE:** Alias name has: `alibabacloudstack_vpcs`

This data source provides a list of VPCs in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpc_vpcs" "vpcs_ds" {
  cidr_block = "172.16.0.0/12"
  status     = "Available"
  name_regex = "^foo"
}

output "first_vpc_id" {
  value = "${data.alibabacloudstack_vpc_vpcs.vpcs_ds.vpcs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional, ForceNew) The CIDR block of the VPC. You can specify one of the following CIDR blocks or their subsets as the primary IPv4 CIDR block of the VPC:
  * Standard private CIDR blocks defined by RFC documents: `192.168.0.0/16`, `172.16.0.0/12`, and `10.0.0.0/8`. The subnet mask must be between 8 and 28 bits in length.
  * Custom CIDR blocks other than `100.64.0.0/10`, `224.0.0.0/4`, `127.0.0.0/8`, `169.254.0.0/16`, and their subnets.
* `status` - (Optional, ForceNew) The status of the VPC. Valid values:
  * `Pending`: The VPC is being configured.
  * `Available`: The VPC is available.
* `name_regex` - (Optional, ForceNew) A regex string to filter VPCs by name.
* `is_default` - (Optional, ForceNew) Specifies whether to create the default VPC in the specified region. Valid values:
  * `true`
  * `false` (default)
* `vswitch_id` - (Optional, ForceNew) Filters results by the specified VSwitch.
* `ids` - (Optional, ForceNew) A list of VPC IDs.
* `dhcp_options_set_id` - (Optional, ForceNew) The ID of the DHCP options set.
* `dry_run` - (Optional, ForceNew) Specifies whether to perform a dry run. Valid values:
  * `true`: Performs only a dry run. The system checks the required parameters, request syntax, and limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  * `false` (default): Performs a dry run and sends the request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which you want to move the resource.
* `vpc_name` - (Optional, ForceNew) The name of the VPC. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `vpc_owner_id` - (Optional, ForceNew) The owner ID of the VPC.
* `enable_details` - (Optional) Default to `true`. Set it to `true` to output the `route_table_id`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of VPC IDs.
* `names` - A list of VPC names.
* `vpcs` - A list of VPCs. Each element contains the following attributes:
  * `id` - The ID of the VPC.
  * `region_id` - The ID of the region where the VPC is located.
  * `resource_group_id` - The ID of the resource group to which the VPC belongs.
  * `status` - The status of the VPC.
  * `vpc_name` - The name of the VPC.
  * `vswitch_ids` - A list of VSwitch IDs in the specified VPC.
  * `cidr_block` - The CIDR block of the VPC.
  * `vrouter_id` - The ID of the VRouter.
  * `route_table_id` - The route table ID of the VRouter.
  * `description` - The description of the VPC.
  * `is_default` - Indicates whether the VPC is the default VPC in the region.
  * `creation_time` - The time when the VPC was created.
  * `tags` - A map of tags assigned to the VPC.
  * `ipv6_cidr_block` - The IPv6 CIDR block of the VPC.
  * `router_id` - The ID of the VRouter.
  * `secondary_cidr_blocks` - A list of secondary IPv4 CIDR blocks of the VPC.
  * `user_cidrs` - A list of user CIDRs.
  * `vpc_id` - The ID of the VPC.
  * `available_ip_address_count` - The count of available IP addresses in the VPC.