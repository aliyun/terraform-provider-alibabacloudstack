---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_security_group"
sidebar_current: "docs-apsarastack-resource-security-group"
description: |-
  Provides a Apsarastack Security Group resource.
---

# apsarastack\_security\_group

Provides a security group resource.

-> **NOTE:** `apsarastack_security_group` is used to build and manage a security group

## Example Usage

Basic Usage

```
resource "apsarastack_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
}
```
Basic usage for vpc

```
resource "apsarastack_security_group" "group" {
  name   = "new-group"
  vpc_id = "${apsarastack_vpc.vpc.id}"
}

resource "apsarastack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security group. Defaults to null.
* `description` - (Optional, Forces new resource) The security group description. Defaults to null.
* `vpc_id` - (Optional, ForceNew) The VPC ID.	

* `inner_access_policy` - (Optional) Whether to allow both machines to access each other on all ports in the same security group. Valid values: ["Accept", "Drop"]
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group

