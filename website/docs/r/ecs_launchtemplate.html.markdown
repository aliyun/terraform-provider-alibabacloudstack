---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_launchtemplate"
sidebar_current: "docs-Alibabacloudstack-ecs-launchtemplate"
description: |- 
  Provides a ecs Launchtemplate resource.
---

# alibabacloudstack_ecs_launchtemplate
-> **NOTE:** Alias name has: `alibabacloudstack_launch_template`

Provides a ecs Launchtemplate resource.

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
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testaccLaunchTemplateBasic12183"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = alibabacloudstack_security_group.default.id
  cidr_ip            = "172.16.0.0/24"
}

resource "alibabacloudstack_launch_template" "default" {
  name                          = var.name
  description                   = "Test launch template"
  host_name                     = var.name
  image_id                      = data.alibabacloudstack_images.default.images.0.id
  instance_name                 = var.name
  instance_type                 = local.default_instance_type_id
  internet_max_bandwidth_in     = 5
  internet_max_bandwidth_out    = 0
  io_optimized                  = "none"
  key_pair_name                 = "test-key-pair"
  ram_role_name                 = "xxxxx"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = 5
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_id             = alibabacloudstack_security_group.default.id
  system_disk_category          = "cloud_ssd"
  system_disk_description       = "Test disk"
  system_disk_name              = "hello"
  system_disk_size              = 40
  resource_group_id             = "rg-zkdfjahg9zxncv0"
  userdata                      = "xxxxxxxxxxxxxx"
  vswitch_id                    = alibabacloudstack_vswitch.default.id
  vpc_id                        = alibabacloudstack_vpc.default.id
  zone_id                       = data.alibabacloudstack_zones.default.zones.0.id

  tags = {
    tag1 = "hello"
    tag2 = "world"
  }

  network_interfaces {
    name              = "eth0"
    description       = "NI"
    primary_ip        = "10.0.0.2"
    security_group_id = "xxxx"
    vswitch_id        = "xxxxxxx"
  }

  data_disks {
    name        = "disk1"
    size        = 20
    category    = "cloud_efficiency"
    description = "test1"
  }

  data_disks {
    name        = "disk2"
    size        = 30
    category    = "cloud_ssd"
    description = "test2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) The name of the launch template. It must start with an English letter (uppercase or lowercase) and can contain numbers, periods (.), colons (:), underscores (_), and hyphens (-). Length should be between 2 and 128 characters. Cannot start with "http://" or "https://".
* `description` - (Optional) Description of the launch template version 1. Length should be between 2 and 256 characters. Cannot start with "http://" or "https://". Default value is null.
* `host_name` - (Optional) Host name of the instance. For Windows instances, the length should be between 2 and 15 characters, cannot start or end with a period (.) or a hyphen (-), and cannot have two or more consecutive periods (.) or hyphens (-). For other operating systems, the length should be between 2 and 64 characters.
* `image_id` - (Optional) Image ID used to create the instance.
* `image_owner_alias` - (Optional) The source of the image. Valid values:
  * `system`: Public images provided by Alibaba Cloud.
  * `self`: Custom images created by you.
  * `others`: Shared images from another Alibaba Cloud account.
  * `marketplace`: Marketplace images.
* `instance_charge_type` - (Optional) Billing method for the instance. Valid values:
  * `PrePaid`: Monthly or annual subscription.
  * `PostPaid`: Pay-as-you-go.
* `instance_name` - (Optional) Name of the instance. Must begin with an English or Chinese character and can contain A-Z, a-z, Chinese characters, numbers, periods (.), colons (:), underscores (_), and hyphens (-). Length should be between 2 and 128 characters.
* `instance_type` - (Optional) Instance type. You can use the `data alibabacloudstack_instances` data source to obtain the latest list of instance types.
* `internet_charge_type` - (Optional) Network billing method. Valid values:
  * `PayByBandwidth`: Billed by fixed bandwidth.
  * `PayByTraffic`: Billed by usage traffic.
* `internet_max_bandwidth_in` - (Optional) Maximum inbound bandwidth from the Internet, measured in Mbit/s. Value range: [1, 200].
* `internet_max_bandwidth_out` - (Optional) Maximum outbound bandwidth from the Internet, measured in Mbit/s. Value range: [0, 100].
* `io_optimized` - (Optional) Whether the instance is I/O optimized. Valid values:
  * `none`
  * `optimized`
* `key_pair_name` - (Optional) Name of the key pair used for SSH login. Ignored for Windows instances.
* `network_type` - (Optional) Network type of the instance. Valid values: `Classic`, `VPC`.
* `ram_role_name` - (Optional) RAM role name assigned to the instance.
* `resource_group_id` - (Optional) ID of the resource group to which the instance belongs.
* `security_enhancement_strategy` - (Optional) Whether to activate the security enhancement feature. Valid values: `Active`, `Deactive`.
* `security_group_id` - (Optional) Security group ID.
* `spot_price_limit` - (Optional) Maximum hourly price for a spot instance. Supports up to three decimal places.
* `spot_strategy` - (Optional) Spot strategy for a pay-as-you-go instance. Valid values:
  * `NoSpot`: Normal pay-as-you-go instance.
  * `SpotWithPriceLimit`: Spot instance with a maximum price limit.
  * `SpotAsPriceGo`: System automatically calculates the price.
* `system_disk_category` - (Optional) Category of the system disk. Valid values:
  * `cloud`: Basic cloud disk.
  * `cloud_efficiency`: Ultra cloud disk.
  * `cloud_ssd`: SSD cloud disk.
  * `ephemeral_ssd`: Local SSD disk.
  * `cloud_essd`: ESSD cloud disk.
* `system_disk_description` - (Optional) Description of the system disk.
* `system_disk_name` - (Optional) Name of the system disk.
* `system_disk_size` - (Optional) Size of the system disk, measured in GB. Value range: [20, 500].
* `userdata` - (Optional) User-defined data of the instance, Base64-encoded. Raw data size cannot exceed 16 KB.
* `vswitch_id` - (Optional) VSwitch ID when creating a VPC-connected instance.
* `vpc_id` - (Optional) VPC ID.
* `zone_id` - (Optional) Zone ID of the instance.
* `network_interfaces` - (Optional) List of network interfaces created with the instance.
  * `name` - (Optional) Name of the ENI.
  * `description` - (Optional) Description of the ENI.
  * `primary_ip` - (Optional) Primary private IP address of the ENI.
  * `security_group_id` - (Optional) Security group ID for the ENI.
  * `vswitch_id` - (Optional) VSwitch ID for the ENI.
* `data_disks` - (Optional) List of data disks created with the instance.
  * `name` - (Optional) Name of the data disk.
  * `size` - (Required) Size of the data disk, measured in GB.
    - `cloud`: [5, 2000]
    - `cloud_efficiency`: [20, 32768]
    - `cloud_ssd`: [20, 32768]
    - `cloud_essd`: [20, 32768]
    - `ephemeral_ssd`: [5, 800]
  * `category` - (Optional) Category of the data disk. Default is `cloud_efficiency`.
  * `encrypted` - (Optional, Bool) Whether the data disk is encrypted. Default is `false`.
  * `snapshot_id` - (Optional) Snapshot ID used to initialize the data disk.
  * `delete_with_instance` - (Optional) Whether to delete the data disk when the instance is destroyed. Default is `true`.
  * `description` - (Optional) Description of the data disk.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: Up to 64 characters in length. Cannot start with "aliyun", "acs:", "http://", or "https://".
  - Value: Up to 128 characters in length. Can be a null string.
* `auto_release_time` - (Optional) Scheduled release time for the instance.
* `user_data` - (Optional) User-defined data of the instance, Base64-encoded. Raw data size cannot exceed 16 KB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the launch template.
* `launch_template_name` - The name of the launch template.
* `internet_max_bandwidth_in` - The maximum public inbound bandwidth, in Mbit/s.
* `internet_max_bandwidth_out` - The maximum public outbound bandwidth, in Mbit/s.
* `name` - (Computed) The name of the launch template.