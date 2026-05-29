---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_security_group"
sidebar_current: "docs-Alibabacloudstack-bmcp-security-group"
description: |-
  Provides a BMCP Security Group resource.
---

# alibabacloudstack_bmcp_security_group

**NOTE:** This resource can also be referred to by the following alias:
- `alibabacloudstack_bcmp_security_group`

Provides a BMCP (Bare Metal Compute Platform) Security Group resource.

## Example Usage

```hcl
resource "alibabacloudstack_bmcp_security_group" "default" {
  vpc_id      = "vpc-7kkaxv72063xkz7exgv3x"
  name        = "my-security-group"
  description = "My BMCP Security Group"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC where the security group is located. Modifying this parameter will force a new resource to be created.
* `name` - (Required) The name of the security group. It must be 2 to 128 characters in length.
* `description` - (Optional) The description of the security group. It must be 0 to 256 characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the security group.
* `sg_id` - The ID of the security group.
* `name` - The name of the security group.
* `description` - The description of the security group.
* `vpc_id` - The ID of the VPC where the security group is located.
* `resource_group` - The resource group ID of the security group.
* `resource_group_name` - The resource group name of the security group.
* `department` - The department ID of the security group.
* `department_name` - The department name of the security group.
* `region_id` - The region ID of the security group.
* `ascm_create_user` - The ASCM user who created the security group.
* `create_time` - The creation time of the security group.
* `update_time` - The last update time of the security group.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain actions:

- `create` - (Defaults to 10 minutes) Used when creating the BMCP security group.
- `delete` - (Defaults to 10 minutes) Used when deleting the BMCP security group.

## Import

BMCP Security Group can be imported using the security group ID (sgId), e.g.

```
$ terraform import alibabacloudstack_bmcp_security_group.example sg-12345678
```
