---
subcategory: "ECS"  
layout: "alibabacloudstack"  
page_title: "Alibabacloudstack: alibabacloudstack_ecs_image"  
sidebar_current: "docs-Alibabacloudstack-ecs-image"  
description: |-  
  Provides a ecs Image resource.  
---

# alibabacloudstack_ecs_image
-> **NOTE:** Alias name has: `alibabacloudstack_image`

Provides an ECS image resource. You can then use a custom image to create ECS instances (RunInstances) or change the system disk for an existing instance (ReplaceSystemDisk).

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccEcsImageShareConfigBasic4783"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
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
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id              = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type         = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id              = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false

  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_ecs_image" "default" {
  instance_id = "${alibabacloudstack_ecs_instance.default.id}"
  image_name  = "${var.name}"
  description = "Custom image created by Terraform"
  tags = {
    Environment = "Test"
    Owner      = "Terraform"
  }
}

resource "alibabacloudstack_image_share_permission" "default" {
  image_id   = "${alibabacloudstack_ecs_image.default.id}"
  account_id = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew) The ID of the instance used to create the custom image.
* `snapshot_id` - (Optional, ForceNew) The snapshot ID used to create the custom image. Conflicts with `instance_id` and `disk_device_mapping`.
* `image_name` - (Optional) The name of the image. It must be 2 to 128 characters in length, starting with a letter or Chinese character. It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of the image. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `tags` - (Optional) A mapping of tags to assign to the resource. Maximum of 20 tag-value pairs.
* `disk_device_mapping` - (Optional, ForceNew) Description of the system disk and snapshots under the image. Conflicts with `snapshot_id` and `instance_id`. Each `disk_device_mapping` supports the following:
  * `size` - (Optional, ForceNew) Specifies the size of the disk in the combined custom image, in GiB. Value range: 5 to 2000.
  * `snapshot_id` - (Optional, ForceNew) Specifies the snapshot used to create the combined custom image.
* `force` - (Optional) Indicates whether to force delete the custom image. Default is `false`.
  - `true`: Force deletes the custom image, regardless of whether the image is currently being used by other instances.
  - `false`: Verifies that the image is not currently in use by any other instances before deleting the image.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### Timeouts

* `create` - (Defaults to 10 mins) Used when creating the image (until it reaches the initial `Available` status).
* `delete` - (Defaults to 10 mins) Used when terminating the image.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the image.
* `image_name` - The name of the image.
* `description` - The description of the image.
* `disk_device_mapping` - The disk device mappings associated with the image.