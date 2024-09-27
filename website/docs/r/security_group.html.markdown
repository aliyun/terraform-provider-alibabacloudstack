---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_security_group"
sidebar_current: "docs-alibabacloudstack-resource-security-group"
description: |-
  Provides a Alibabacloudstack Security Group resource.
---

# alibabacloudstack\_security\_group

Provides a security group resource.

-> **NOTE:** `alibabacloudstack_security_group` is used to build and manage a security group

## Example Usage

Basic Usage

```
resource "alibabacloudstack_security_group" "group" {
  name        = "terraform-test-group"
  description = "New security group"
}
```
Basic usage for vpc

```
resource "alibabacloudstack_security_group" "group" {
  name   = "new-group"
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security group. Defaults to null.
* `description` - (Optional, Forces new resource) The security group description. Defaults to null.
* `vpc_id` - (Optional, ForceNew) The VPC ID.	
* `type` - (Optional, ForceNew) The type of the security group. Valid values: `normal`, `enterprise`. Default value is `normal`.
* `inner_access_policy` - (Optional) Whether to allow both machines to access each other on all ports in the same security group. Valid values: ["Accept", "Drop"]
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group

