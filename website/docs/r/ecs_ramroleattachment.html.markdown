---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_ramroleattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-ramroleattachment"
description: |- 
  Provides a ecs Ramroleattachment resource.
---

# alibabacloudstack_ecs_ramroleattachment
-> **NOTE:** Alias name has: `alibabacloudstack_ram_role_attachment`

Provides a ECS Ramroleattachment resource to bind a RAM role to one or more ECS instances.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "Test_ram_role_attachment"
}

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "192.168.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  image_id           = data.alibabacloudstack_images.default.images.0.id
  instance_type      = local.default_instance_type_id
  instance_name      = var.name
  security_groups    = [alibabacloudstack_security_group.default.id]
  availability_zone  = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_pperf"
  system_disk_size  = 100
  vswitch_id        = alibabacloudstack_vswitch.default.id
}

data "alibabacloudstack_ascm_ram_service_roles" "role" {
  product = "ecs"
}

resource "alibabacloudstack_ecs_ramroleattachment" "default" {
  role_name    = data.alibabacloudstack_ascm_ram_service_roles.role.roles.0.name
  instance_ids = [alibabacloudstack_instance.default.id]
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) The name of the RAM role to be attached. This name must be between 1 and 64 characters in length and can only contain alphanumeric characters or hyphens (`-`, `_`). It cannot start with a hyphen.
* `instance_ids` - (Required, ForceNew) A list of ECS instance IDs to which the RAM role will be attached.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_name` - The name of the RAM role that has been attached.
* `instance_ids` - The list of ECS instance IDs to which the RAM role has been attached.
