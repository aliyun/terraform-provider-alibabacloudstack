---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_instance"
sidebar_current: "docs-Alibabacloudstack-ecs-instance"
description: |- 
  Provides a ecs Instance resource.
---

# alibabacloudstack_ecs_instance
-> **NOTE:** Alias name has: `alibabacloudstack_instance`

Provides a ecs Instance resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" default {
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
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
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

variable "name" {
  default = "tf-testAccEcsInstanceConfigBasic2648"
}

resource "alibabacloudstack_ecs_instance" "default" {
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  instance_name        = "${var.name}"
  user_data            = "I_am_user_data"
  security_groups      = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
  vswitch_id          = "${alibabacloudstack_vpc_vswitch.default.id}"
  tags = {
    Bar = "Bar"
    foo = "foo"
  }
  image_id            = "${data.alibabacloudstack_images.default.images.0.id}"
  security_enhancement_strategy = "Active"
  instance_type       = "${local.default_instance_type_id}"
  availability_zone   = "${data.alibabacloudstack_zones.default.zones[0].id}"

  # IPv6 Configuration
  enable_ipv6         = true
  ipv6_cidr_block     = "fd00::/64"
  ipv6_address_count  = 3

  # Data Disks
  data_disks = [
    {
      category         = "cloud_efficiency"
      size             = 50
      delete_with_instance = true
    },
    {
      category         = "cloud_ssd"
      size             = 100
      snapshot_id      = "snap-12345678"
      delete_with_instance = false
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, ForceNew) The availability zone where the instance is launched. It must match the zone of the VSwitch if `vswitch_id` is specified.
* `zone_id` - (Optional, ForceNew) The ID of the available zone to which the instance belongs.
* `image_id` - (Required) The image ID used for the instance. Changing this value will force a new resource.
* `instance_type` - (Required) The type of instance to start. Changing this value will force a new resource.
* `instance_name` - (Optional) The name of the ECS instance. This can be a string of 2 to 128 characters and cannot begin with http:// or https://.
* `description` - (Optional) A description of the instance. This can be a string of 2 to 256 characters and cannot begin with http:// or https://.
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. Default is 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range: [0, 100]. Default is 0 Mbps.
* `host_name` - (Optional) Host name of the ECS instance. Changing this value will cause the instance to reboot.
* `password` - (Optional, Sensitive) Password for the instance. Must be 8 to 30 characters long and include uppercase/lowercase letters and numbers. Changing this value will cause the instance to reboot.
* `kms_encrypted_password` - (Optional) An KMS encrypted password used for the instance. If `password` is provided, this field will be ignored. Changing this value will cause the instance to reboot.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password`. Valid when `kms_encrypted_password` is set. Changing this value will cause the instance to reboot.
* `is_outdated` - (Optional) Whether to use outdated instance types. Default is `false`.
* `system_disk_category` - (Optional, ForceNew) The category of the system disk. Valid values: `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`. Default is `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of the system disk, measured in GiB. Value range: [20, 500]. Default is the maximum of {40, ImageSize}.
* `system_disk_name` - (Optional) Name of the system disk. Changing this value will cause the instance to reboot.
* `system_disk_description` - (Optional) Description of the system disk. Changing this value will cause the instance to reboot.
* `data_disks` - (Optional, ForceNew) A list of data disks created with the instance. Each data disk supports the following properties:
  * `category` - (Optional, ForceNew) The category of the data disk. Valid values: `cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`. Default is `cloud_efficiency`.
  * `size` - (Required, ForceNew) The size of the data disk in GiB.
  * `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk.
  * `delete_with_instance` - (Optional, ForceNew) Whether to delete the data disk when the instance is destroyed. Default is `true`.
  * `encrypted` - (Optional, Bool, ForceNew) Whether to encrypt the data disk. Default is `false`.
  * `kms_key_id` - (Optional) The KMS key ID corresponding to the data disk.
  * `name` - (Optional, ForceNew) The name of the data disk.
  * `description` - (Optional, ForceNew) The description of the data disk.
* `subnet_id` - (Removed since v1.210.0) The ID of the subnet. Conflicts with `vswitch_id`.
* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. This parameter must be set unless you can create classic network instances.
* `private_ip` - (Optional) The private IP address assigned to the instance. It is valid when `vswitch_id` is specified.
* `hpc_cluster_id` - (Optional, ForceNew) The ID of the Elastic High Performance Computing (E-HPC) cluster to which the instance belongs.
* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance. Changing this value will cause the instance to reboot.
* `role_name` - (Optional, ForceNew) The name of the RAM role associated with the instance.
* `key_name` - (Optional, ForceNew) The name of the key pair used for the instance.
* `storage_set_id` - (Optional, ForceNew) The ID of the storage set.
* `storage_set_partition_number` - (Optional, ForceNew) The number of partitions in the storage set.
* `security_enhancement_strategy` - (Optional, ForceNew) The security enhancement strategy. Valid values: `Active` (enable security enhancement strategy), `Deactive` (disable security enhancement strategy).
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable IPv6. Valid values: `false` (disable), `true` (enable).
* `ipv6_cidr_block` - (Optional) The IPv6 CIDR block of the VPC.
* `ipv6_address_count` - (Optional) The count of IPv6 addresses requested for allocation. If `enable_ipv6` is `true`, this value must be greater than 0.
* `ipv6_address_list` - (Optional, ForceNew) A list of IPv6 addresses to be assigned to the primary ENI. Supports up to 10 addresses.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `system_disk_tags` - (Optional) A mapping of tags to assign to the system disk.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `availability_zone` - The availability zone of the instance.
* `zone_id` - The ID of the available zone to which the instance belongs.
* `internet_max_bandwidth_in` - Maximum public access bandwidth.
* `host_name` - The host name of the instance.
* `system_disk_id` - The ID of the system disk.
* `subnet_id` - The ID of the subnet.
* `private_ip` - The private IP address of the instance.
* `hpc_cluster_id` - The ID of the HPC cluster to which the instance belongs.
* `status` - The status of the instance.
* `role_name` - The name of the RAM role associated with the instance.
* `key_name` - The name of the key pair used for the instance.
* `storage_set_id` - The ID of the storage set.
* `storage_set_partition_number` - The number of partitions in the storage set.