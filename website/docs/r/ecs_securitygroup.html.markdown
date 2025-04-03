---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_securitygroup"
sidebar_current: "docs-Alibabacloudstack-ecs-securitygroup"
description: |- 
  Provides a ecs Securitygroup resource.
---

# alibabacloudstack_ecs_securitygroup
-> **NOTE:** Alias name has: `alibabacloudstack_security_group`

Provides a ecs Securitygroup resource.

## Example Usage

### Basic Usage

```hcl
resource "alibabacloudstack_security_group" "basic_group" {
  name        = "terraform-test-group"
  description = "New security group"
}
```

### Basic Usage for VPC

```hcl
resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "group" {
  name           = "new-group"
  vpc_id         = alibabacloudstack_vpc.vpc.id
  description    = "Security group for VPC"
  type          = "normal"
  inner_access_policy = "Accept"
}
```

### Advanced Usage with Tags

```hcl
resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_security_group" "group_with_tags" {
  name           = "tagged-group"
  vpc_id         = alibabacloudstack_vpc.vpc.id
  description    = "Security group with tags"
  type          = "enterprise"
  inner_access_policy = "Drop"

  tags = {
    Environment = "Production"
    Owner      = "DevOps"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security group. It must be 2 to 128 characters in length and can contain letters, numbers, underscores (`_`), periods (`.`), and hyphens (`-`). It must start with a letter or Chinese character and cannot start with `http://` or `https://`. If not specified, Terraform will automatically generate a unique name.
  
* `description` - (Optional) The description of the security group. It must be 2 to 256 characters in length and cannot start with `http://` or `https://`. By default, this parameter is left empty.

* `vpc_id` - (Optional, ForceNew) The ID of the VPC in which you want to create the security group. This parameter is required only if you want to create security groups of the VPC type. In regions that support the classic network, you can create security groups of the classic network type without specifying the `vpc_id`.

* `type` - (Optional, ForceNew) The type of the security group. Valid values:
  * `normal`: Standard security group (default).
  * `enterprise`: Enterprise-level security group.

* `inner_access_policy` - (Optional) The internal access policy of the security group. Valid values:
  * `Accept`: All instances in the security group can communicate with each other.
  * `Drop`: All instances in the security group are isolated from each other.
  The value of this parameter is not case-sensitive. Default value is `Accept`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Each tag consists of a key-value pair. Tag keys must be unique within the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the security group.

* `inner_access_policy` - The internal access policy of the security group. Valid values:
  * `Accept`: All instances in the security group can communicate with each other.
  * `Drop`: All instances in the security group are isolated from each other.
  The value of this parameter is not case-sensitive. This attribute reflects the actual configuration of the security group.