---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_securitygroups"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-securitygroups"
description: |- 
  Provides a list of ecs securitygroups owned by an Alibabacloudstack account.
---

# alibabacloudstack_ecs_securitygroups
-> **NOTE:** Alias name has: `alibabacloudstack_security_groups`

This data source provides a list of ECS Security Groups in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
# Filter security groups and print the results into a file
data "alibabacloudstack_ecs_securitygroups" "sec_groups_ds" {
  name_regex = "^web-"
}

# In conjunction with a VPC
data "alibabacloudstack_ecs_securitygroups" "primary_sec_groups_ds" {
  vpc_id = var.vpc_id
}

output "first_group_id" {
  value = "${data.alibabacloudstack_ecs_securitygroups.primary_sec_groups_ds.groups.0.id}"
}

# Using tags to filter security groups
data "alibabacloudstack_ecs_securitygroups" "taggedSecurityGroups" {
  tags = {
    Environment = "Production"
    Department  = "Finance"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter the resulting security groups by their names.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC in which you want to retrieve security groups. This parameter is required only if you want to retrieve security groups of the VPC type. In regions that support the classic network, you can retrieve security groups of the classic network type without specifying this parameter.
* `ids` - (Optional) A list of Security Group IDs to filter the results.
* `tags` - (Optional) A map of tags assigned to the security groups. It must be in the format:
  ```hcl
  tags = {
    tagKey1 = "tagValue1",
    tagKey2 = "tagValue2"
  }
  ```

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Security Group IDs.
* `names` - A list of Security Group names.
* `groups` - A list of Security Groups. Each element contains the following attributes:
  * `id` - The ID of the security group.
  * `name` - The name of the security group.
  * `description` - The description of the security group. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`. By default, this parameter is left empty.
  * `vpc_id` - The ID of the VPC that owns the security group. This parameter is required only if you want to create or retrieve security groups of the VPC type. In regions that support the classic network, you can create or retrieve security groups of the classic network type without specifying this parameter.
  * `creation_time` - The creation time of the security group.
  * `tags` - A map of tags assigned to the security group.