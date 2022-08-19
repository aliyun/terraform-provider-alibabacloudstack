---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ecs_launch_template"
sidebar_current: "docs-apsarastack-resource-ecs-launch-template"
description: |-
  Provides a Apsarastack ECS Launch Template resource.
---

# apsarastack\_ecs\_launch\_template

Provides a ECS Launch Template resource.

For information about ECS Launch Template and how to use it, see [What is Launch Template](https://help.aliyun.com/document_detail/74686.html?spm=5176.21213303.J_6704733920.7.6c2853c9FgD3aj&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%4074686._.ID_help%40%40%E6%96%87%E6%A1%A3%40%4074686-RL_CreateLaunchTemplate-LOC_main-OR_ser-V_2-P0_0).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "apsarastack_instance_types" "default" {
  availability_zone = data.apsarastack_zones.default.zones[0].id
  cpu_core_count    = 2
  sorted_by         = "Memory"
 }
 locals {
  instance_type_id = sort(data.apsarastack_instance_types.default.ids)[0]
}
 
data "apsarastack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
variable "name" {
  default = "tf-testaccLaunchTemplateBasic3141889771392864639"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${apsarastack_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "apsarastack_launch_template" "default" {
  network_type = "vpc"
  key_pair_name = "tf-testaccLaunchTemplateBasic3141889771392864639"
  host_name = "tf-testaccLaunchTemplateBasic3141889771392864639"
  image_id = "${data.apsarastack_images.default.images.0.id}"
  spot_price_limit = "5"
  instance_name = "tf-testaccLaunchTemplateBasic3141889771392864639"
  internet_max_bandwidth_in = "5"
  ram_role_name = "tf-testaccLaunchTemplateBasic3141889771392864639"
  data_disks {
    name = "disk1"
    description = "test1"
  }
  data_disks {
    description = "test2"
    name = "disk2"
  }
  
  zone_id = "beijing-a"
  internet_charge_type = "PayByBandwidth"
  network_interfaces {
    name = "eth0"
    description = "hello1"
    primary_ip = "10.0.0.2"
    security_group_id = "xxxx"
    vswitch_id = "xxxxxxx"
  }
  
  tags = {
           tag2 = "world"
           tag1 = "hello"
         }
  vswitch_id = "sw-ljkngaksdjfj0nnasdf"
  instance_charge_type = "PrePaid"
  spot_strategy = "SpotWithPriceLimit"
  vpc_id = "vpc-asdfnbg0as8dfk1nb2"
  security_enhancement_strategy = "Active"
  internet_max_bandwidth_out = "0"
  system_disk_category = "cloud_ssd"
  security_group_id = "${apsarastack_security_group.default.id}"
  io_optimized = "none"
  system_disk_size = "40"
  system_disk_description = "tf-testaccLaunchTemplateBasic3141889771392864639"
  resource_group_id = "rg-zkdfjahg9zxncv0"
  instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
  system_disk_name = "tf-testaccLaunchTemplateBasic3141889771392864639"
  userdata = "xxxxxxxxxxxxxx"
  description = "tf-testaccLaunchTemplateBasic3141889771392864639"
  name = "tf-testaccLaunchTemplateBasic3141889771392864639"
}
```

## Argument Reference

The following arguments are supported:

* `auto_release_time` - (Optional) Instance auto release time. The time is presented using the ISO8601 standard and in UTC time. The format is  YYYY-MM-DDTHH:MM:SSZ.
* `data_disks` - (Optional) The list of data disks created with instance.
* `deployment_set_id` - (Optional) The Deployment Set Id.
* `description` - (Optional) Description of instance launch template version 1. It can be [2, 256] characters in length. It cannot start with "http://" or "https://". The default value is null.
* `enable_vm_os_config` - (Optional) Whether to enable the instance operating system configuration.
* `host_name` - (Optional) Instance host name.It cannot start or end with a period (.) or a hyphen (-) and it cannot have two or more consecutive periods (.) or hyphens (-).For Windows: The host name can be [2, 15] characters in length. It can contain A-Z, a-z, numbers, periods (.), and hyphens (-). It cannot only contain numbers. For other operating systems: The host name can be [2, 64] characters in length. It can be segments separated by periods (.). It can contain A-Z, a-z, numbers, and hyphens (-).
* `image_id` - (Optional) The Image ID.
* `image_owner_alias` - (Optional) Mirror source. Valid values: `system`, `self`, `others`, `marketplace`, `""`. Default to: `""`.
* `instance_charge_type` - (Optional) Billing methods. Valid values: `PostPaid`, `PrePaid`.
* `instance_type` - (Optional) Instance type. For more information, call resource_alicloud_instances to obtain the latest instance type list.
* `internet_charge_type` - (Optional) Internet bandwidth billing method. Valid values: `PayByTraffic`, `PayByBandwidth`.
* `internet_max_bandwidth_in` - (Optional) The maximum inbound bandwidth from the Internet network, measured in Mbit/s. Value range: [1, 200].
* `internet_max_bandwidth_out` - (Optional) Maximum outbound bandwidth from the Internet, its unit of measurement is Mbit/s. Value range: [0, 100].
* `io_optimized` - (Optional) Whether it is an I/O-optimized instance or not. Valid values: `none`, `optimized`.
* `key_pair_name` - (Optional) The name of the key pair.
    - Ignore this parameter for Windows instances. It is null by default. Even if you enter this parameter, only the  Password content is used.
    - The password logon method for Linux instances is set to forbidden upon initialization.
* `launch_template_name` - (Optional, ForceNew) The name of Launch Template.
* `network_interfaces` - (Optional) The list of network interfaces created with instance.
* `network_type` - (Optional) Network type of the instance. Valid values: `classic`, `vpc`.
* `password_inherit` - (Optional) Whether to use the password preset by the mirror.
* `period` - (Optional) The subscription period of the instance. Unit: months. This parameter takes effect and is required only when InstanceChargeType is set to PrePaid. If the DedicatedHostId parameter is specified, the value of the Period parameter must be within the subscription period of the dedicated host.
    - When the PeriodUnit parameter is set to `Week`, the valid values of the Period parameter are `1`, `2`, `3`, and `4`.
    - When the PeriodUnit parameter is set to `Month`, the valid values of the Period parameter are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`, `48`, and `60`.
* `private_ip_address` - (Optional) The private IP address of the instance.
* `ram_role_name` - (Optional) The RAM role name of the instance. You can use the RAM API ListRoles to query instance RAM role names.
* `resource_group_id` - (Optional) The ID of the resource group to which to assign the instance, Elastic Block Storage (EBS) device, and ENI.
* `security_enhancement_strategy` - (Optional) Whether or not to activate the security enhancement feature and install network security software free of charge. Valid values: `Active`, `Deactive`.
* `security_group_id` - (Optional) The security group ID.
* `security_group_ids` - (Optional) The ID of security group N to which to assign the instance.
* `spot_duration` - (Optional, Computed) The protection period of the preemptible instance. Unit: hours. Valid values: `0`, `1`, `2`, `3`, `4`, `5`, and `6`. Default to: `1`.
* `spot_price_limit` -(Optional) Sets the maximum hourly instance price. Supports up to three decimal places.
* `spot_strategy` - (Optional) The spot strategy for a Pay-As-You-Go instance. This parameter is valid and required only when InstanceChargeType is set to PostPaid. Valid values: `NoSpot`, `SpotAsPriceGo`, `SpotWithPriceLimit`.
* `system_disk` - (Optional) The System Disk.
* `template_resource_group_id` - (Optional, ForceNew) The template resource group id.
* `user_data` - (Optional, Computed) The User Data.
* `version_description` - (Optional) The description of the launch template version. The description must be 2 to 256 characters in length and cannot start with http:// or https://.                                    
* `vswitch_id` - (Optional) When creating a VPC-Connected instance, you must specify its VSwitch ID.
* `zone_id` - (Optional) The zone ID of the instance.
* `tags` - (Optional) A mapping of tags to assign to instance, block storage, and elastic network.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `template_tags` - (Optional) A mapping of tags to assign to the launch template.
  

#### Block system_disk

The system_disk supports the following: 

* `category` - (Optional, Computed) The category of the system disk. System disk type. Valid values: `all`, `cloud`, `ephemeral_ssd`, `cloud_essd`, `cloud_efficiency`, `cloud_ssd`, `local_disk`.
* `delete_with_instance` - (Optional) Specifies whether to release the system disk when the instance is released. Default to `true`.
* `description` - (Optional, Computed) System disk description. It cannot begin with http:// or https://.
* `iops` - (Optional) The Iops.
* `name` - (Optional, Computed) System disk name. The name is a string of 2 to 128 characters. It must begin with an English or a Chinese character. It can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-).
* `performance_level` - (Optional) The performance level of the ESSD used as the system disk. Valid Values: `PL0`, `PL1`, `PL2`, and `PL3`. Default to: `PL0`.
* `size` - (Optional, Computed) Size of the system disk, measured in GB. Value range: [20, 500].

#### Block network_interfaces

The network_interfaces supports the following: 

* `description` - (Optional) The ENI description.
* `name` - (Optional) The ENI name.
* `primary_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID must be one in the same VPC.
* `vswitch_id` - (Optional) The VSwitch ID for ENI. The instance must be in the same zone of the same VPC network as the ENI, but they may belong to different VSwitches.

#### Block data_disks

The data_disks supports the following: 

* `category` - (Optional) The category of the disk.
* `delete_with_instance` - (Optional) Indicates whether the data disk is released with the instance.
* `description` - (Optional) The description of the data disk.
* `encrypted` - (Optional) Encrypted the data in this disk.
* `name` - (Optional) The name of the data disk.
* `performance_level` - (Optional) The performance level of the ESSD used as the data disk.
* `size` - (Optional) The size of the data disk.
* `snapshot_id` - (Optional) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Launch Template.

## Import

ECS Launch Template can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_launch_template.example <id>
```
