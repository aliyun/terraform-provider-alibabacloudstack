---
layout: "alicloud-doc"
page_title: "Resource: alibabacloudstack_bcmp_security_group"
subcategory: "Bare Metal Computing Platform (BMCP)"
---

# alibabacloudstack_bcmp_security_group

Provides a BMCP Security Group resource.

-> **Note:** This resource can also be referred to by the following alias:
- `alibabacloudstack_bcmp_security_group`

## Example Usage

```terraform
variable "name" {
  default = "tf-example-bmcp-sg"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = alibabacloudstack_vpc.default.id
}
```

## Argument Reference

The following arguments are supported:

- `vpc_id` - (Required, ForceNew) The ID of the VPC to which the security group belongs. Modifying this parameter will force a new resource to be created.
- `name` - (Required) The name of the security group. The length must be between 2 and 128 characters.
- `description` - (Optional) The description of the security group. The length must be between 0 and 256 characters.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the security group.
- `sg_id` - The security group ID.
- `name` - The name of the security group.
- `description` - The description of the security group.
- `vpc_id` - The ID of the VPC.
- `resource_group` - The ID of the resource group.
- `resource_group_name` - The name of the resource group.
- `department` - The ID of the department.
- `department_name` - The name of the department.
- `region_id` - The ID of the region.
- `ascm_create_user` - The ASCM user who created the security group.
- `create_time` - The time when the security group was created.
- `update_time` - The time when the security group was last updated.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain actions:

- `create` - (Defaults to 10 minutes) Used when creating the BMCP security group.
- `delete` - (Defaults to 10 minutes) Used when deleting the BMCP security group.

## Import

BMCP Security Group can be imported using the security group ID (sgId), e.g.

```
$ terraform import alibabacloudstack_bcmp_security_group.example sg-12345678
```
