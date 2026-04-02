---
subcategory: "Bare Metal Compute Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_security_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-security-groups"
description: |- 
  Provides a list of BMCP Security Groups owned by an AlibabacloudStack account.
---

# alibabacloudstack_bmcp_security_groups

This data source provides a list of BMCP (Bare Metal Compute Platform) Security Groups in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_bmcp_security_groups" "example" {
  vpc_id = "vpc-7kkaxv72063xkz7exgv3x"
  name   = "my-security-group"
}

output "security_group_ids" {
  value = data.alibabacloudstack_bmcp_security_groups.example.ids
}

output "security_group_names" {
  value = data.alibabacloudstack_bmcp_security_groups.example.names
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional) The ID of the VPC to filter results.
* `name` - (Optional) The name of the security group to filter results.
* `name_regex` - (Optional) A regex string to filter results by security group name.
* `ids` - (Optional) A list of security group IDs to filter results. If not specified, all security groups will be considered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of matched security groups.
* `ids` - A list of IDs of matched security groups.
* `security_groups` - A list of matched security groups. Each element contains the following attributes:
  * `sg_id` - The ID of the security group.
  * `name` - The name of the security group.
  * `description` - The description of the security group.
  * `vpc_id` - The ID of the VPC where the security group is located.
  * `resource_group` - The resource group ID of the security group.
  * `resource_group_name` - The resource group name of the security group.
  * `department` - The department ID of the security group.
  * `department_name` - The department name of the security group.
  * `region_id` - The region ID of the security group.
  * `ascm_create_user` - The user who created the security group.
  * `create_time` - The creation time of the security group.
  * `update_time` - The last update time of the security group.
