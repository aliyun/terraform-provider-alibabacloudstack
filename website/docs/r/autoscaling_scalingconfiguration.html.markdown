---
subcategory: "Auto Scaling(ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_scaling_configuration"
sidebar_current: "docs-alibabacloudstack-resource-ess-scaling-configuration"
description: |- 
  Provides a ESS scaling configuration resource.
---

# alibabacloudstack_ess_scaling_configuration

Provides a ESS scaling configuration resource.

-> **NOTE:** Several instance types have outdated in some regions and availability zones, such as `ecs.t1.*`, `ecs.s2.*`, `ecs.n1.*` and so on. If you want to keep them, you should set `is_outdated` to true. For more about the upgraded instance type, refer to `alibabacloudstack_instance_types` datasource.

## Example Usage

```hcl
variable "name" {
  default = "essscalingconfiguration"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vswitch.default.id]
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id  = alibabacloudstack_ess_scaling_group.default.id
  image_id          = data.alibabacloudstack_images.default.images[0].id
  instance_type     = data.alibabacloudstack_instance_types.default.instance_types[0].id
  security_group_ids = [alibabacloudstack_security_group.default.id]
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
  force_delete      = true
  active            = true
  internet_max_bandwidth_in = 50
  system_disk_category = "cloud_efficiency"
  system_disk_size = 40
  user_data = base64encode("echo 'Hello World' > /tmp/hello.txt")
  key_pair_name = "my-key-pair"
  tags = {
    Environment = "Test"
    Owner      = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional) The status of the scaling group configuration within the scaling group. Possible values:
  * `Active`: Indicates that the configuration is active, and the current scaling group will use this configuration to automatically create ECS instances.
  * `Inactive`: Indicates that the configuration is inactive, and the current scaling group will not use this configuration to automatically create ECS instances.
  
* `active` - (Optional) Whether to activate the current scaling configuration in the specified scaling group. Defaults to `false`.
  
* `enable` - (Optional) Whether to enable the specified scaling group (make it active) to which the current scaling configuration belongs.
  
* `scaling_group_id` - (Required, ForceNew) The ID of the scaling group to which the scaling configuration belongs.
  
* `image_id` - (Required) The ID of the image file used when creating an ECS instance.
  
* `instance_type` - (Required) The specification of the ECS instance.
  
* `security_group_ids` - (Required) A list of security group IDs to which the ECS instance belongs.
  
* `deployment_set_id` - (Required) The ID of the deployment set to which the ECS instance belongs.
  
* `zone_id` - (Optional) The zone ID to which the ECS instances belong.
  
* `scaling_configuration_name` - (Optional) The name of the scaling configuration. It must be 2~64 characters long and can contain letters, numbers, underscores (`_`), hyphens (`-`), or periods (`.`). If not specified, the default value is the ID of the scaling configuration.
  
* `internet_max_bandwidth_in` - (Optional) The maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). The value range is [1, 200]. If not specified, the default value is 200 Mbps.
  
* `system_disk_category` - (Optional) The category of the system disk. Valid options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, and `cloud`. Default is `cloud_efficiency`.
  
* `system_disk_size` - (Optional) The size of the system disk, in GiB. The valid range depends on the disk category:
  * `cloud`: [20, 500]
  * `cloud_efficiency`, `cloud_ssd`, `cloud_essd`: [20, 500]
  * `ephemeral_ssd`: [20, 500]
  
* `data_disk` - (Optional) A list of data disks to attach to the ECS instance. See [Block datadisk](#block-datadisk) below for details.
  
* `substitute` - (Optional) The scaling configuration that will automatically become active and replace the current configuration when `active` is set to `false`. This parameter is invalid when `active` is `true`.
  
* `system_disk_auto_snapshot_policy_id` - (Optional) The ID of the auto snapshot policy for the system disk.
  
* `is_outdated` - (Optional) Whether to use outdated instance types. Default is `false`.
  
* `user_data` - (Optional) Custom data for the ECS instance. It must be Base64-encoded, and the raw data must not exceed 16KB.
  
* `ram_role_name` - (Optional) The name of the RAM role for the ECS instance. You can query available RAM roles using the [ListRoles](~~ 28713 ~~) API. To create a RAM role, see [CreateRole](~~ 28710 ~~).
  
* `key_pair_name` - (Optional) The name of the key pair used to log in to the ECS instance. For Windows instances, this parameter will be ignored. For Linux instances, logging in with a password will be disabled if this parameter is specified.
  
* `force_delete` - (Optional) Whether to forcibly delete the last scaling configuration along with its scaling group. Default is `false`.
  
* `tags` - (Optional) A map of tags assigned to the resource. These tags will be applied to the ECS instances created by the scaling group.
  
* `instance_name` - (Optional) The name of the ECS instance.
  
* `override` - (Optional) Whether to overwrite existing data. Default is `false`.
  
* `host_name` - (Optional) The hostname of the server. Restrictions vary based on the operating system:
  * Windows: Hostname length is 2~15 characters, can include uppercase letters, digits, and hyphens (`-`). Cannot start or end with a period (`.`) or hyphen (`-`).
  * Linux: Hostname length is 2~64 characters, can include multiple dots (`.`). Each segment between dots can include uppercase letters, digits, and hyphens (`-`).
* `role_name` - (Optional) The name of the RAM role for the ECS instance.
* `key_name` - (Optional) The name of the key pair used to log in to the ECS instance.

### Block datadisk

The `data_disk` block supports the following:

* `size` - (Optional) The size of the data disk, in GB. The valid range depends on the disk category:
  * `cloud`: [5, 2000]
  * `ephemeral`: [5, 1024]
  * `ephemeral_ssd`: [5, 800]
  * `cloud_efficiency`, `cloud_ssd`, `cloud_essd`: [20, 32768]
  
* `device` - (Optional) The mount point of the data disk. Valid values are `/dev/xvdb` to `/dev/xvdz`.
  
* `category` - (Optional) The category of the data disk. Valid options are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, and `cloud`.
  
* `snapshot_id` - (Optional) The ID of the snapshot used to create the data disk. If specified, the `size` parameter is ignored.
  
* `delete_with_instance` - (Optional) Whether to delete the data disk when releasing the ECS instance. Valid values are `true` or `false`. Default is `true`.
  
* `encrypted` - (Optional) Whether to encrypt the data disk. Valid values are `true` or `false`. Default is `false`.
  
* `kms_key_id` - (Optional) The ID of the CMK used to encrypt the data disk.
  
* `name` - (Optional) The name of the data disk. Must be 2~128 characters long and cannot start with `http://` or `https://`.
  
* `description` - (Optional) The description of the data disk. Must be 2~256 characters long and cannot start with `http://` or `https://`.
  
* `auto_snapshot_policy_id` - (Optional) The ID of the auto snapshot policy for the data disk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the scaling configuration.
  
* `status` - The status of the scaling group configuration within the scaling group. Possible values:
  * `Active`: Indicates that the configuration is active.
  * `Inactive`: Indicates that the configuration is inactive.
  
* `active` - Whether the current scaling configuration is active in the specified scaling group.
  
* `zone_id` - The zone ID to which the ECS instances belong.
  
* `scaling_configuration_name` - The name of the scaling configuration.
  
* `internet_max_bandwidth_in` - The maximum incoming bandwidth from the public network, measured in Mbps.
  
* `substitute` - The scaling configuration that will automatically become active and replace the current configuration when `active` is set to `false`.
  
* `ram_role_name` - The name of the RAM role for the ECS instance.
  
* `role_name` - The name of the RAM role for the ECS instance.
  
* `key_pair_name` - The name of the key pair used to log in to the ECS instance.
* `is_outdated` - Whether to use outdated instance types.
* `key_name` - The name of the key pair used to log in to the ECS instance.