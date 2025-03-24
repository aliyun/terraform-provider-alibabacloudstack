---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshot"
sidebar_current: "docs-Alibabacloudstack-ecs-snapshot"
description: |- 
  Provides a ecs Snapshot resource.
---

# alibabacloudstack_ecs_snapshot
-> **NOTE:** Alias name has: `alibabacloudstack_snapshot`

Provides a ecs Snapshot resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccecssnapshot10982"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone     = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated         = false

  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_ecs_snapshot" "default" {
  description  = "rdk_test_description"
  snapshot_name = "rdk_test_name"
  disk_id       = alibabacloudstack_ecs_instance.default.system_disk_id
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The ID of the disk for which the snapshot will be created.
* `snapshot_name` - (Optional, ForceNew) The name of the snapshot. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (\_), and hyphens (-). It cannot start with `auto` because snapshots whose names start with `auto` are recognized as automatic snapshots.
* `description` - (Optional, ForceNew) The description of the snapshot. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### Timeouts

* `create` - (Defaults to 2 mins) Used when creating the snapshot (until it reaches the initial `SnapshotCreatingAccomplished` status).
* `delete` - (Defaults to 2 mins) Used when terminating the snapshot.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the snapshot.
* `snapshot_name` - The name of the snapshot. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (\_), and hyphens (-). It cannot start with `auto` because snapshots whose names start with `auto` are recognized as automatic snapshots.
* `description` - The description of the snapshot. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.