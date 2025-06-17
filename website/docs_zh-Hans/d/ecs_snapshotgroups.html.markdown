---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshotgroups"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-snapshotgroups"
description: |-
  提供阿里云账号下拥有的ecs snapshotgroups列表。
---

# alibabacloudstack\_ecs\_snapshotgroups

此数据源提供根据指定过滤条件列出的阿里云账号下的ecs snapshotgroups资源列表。

## 示例用法
```
variable "name" {
  default = "tf-testAccEcsSnapshotGroupsDataSource-6474718"
}



data "alibabacloudstack_zones" default {
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
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}



resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}




data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
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


resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  data_disks {
      name                 = "disk1"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
      delete_with_instance = true
    }
  data_disks {
      name                 = "disk2"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
	  delete_with_instance = true
    }
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

data "alibabacloudstack_ecs_disks" "disks" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

resource "alibabacloudstack_ecs_snapshot_group" "default" {
  	description =                   "${var.name}"
	instance_id =                   "${alibabacloudstack_ecs_instance.default.id}"
	instant_access =                "true"
	instant_access_retention_days = "7"
	snapshot_group_name =           "${var.name}"
	disk_ids = [
		"${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.2.id}",
	]
}

 

data "alibabacloudstack_ecs_snapshot_groups" "default" {
  description_regex = "${alibabacloudstack_ecs_snapshot_group.default.description}"
}```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - <!--  propertie not found  -->
  * `name_regex` - (选填) - <!--  propertie not found  -->
  * `description_regex` - (选填) - <!--  propertie not found  -->
  * `instance_id` - (选填) - 实例ID
  * `snapshot_group_name` - (选填) - snapshot group name

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `snapshot_groups` - <!--  propertie not found  -->
    * `id` - <!--  propertie not found  -->
    * `create_time` - 代表创建时间的资源属性字段
    * `description` - 快照分组描述
    * `snapshots` - <!--  propertie not found  -->
    * `instance_id` - 实例ID
    * `snapshot_group_id` - 代表资源一级ID的资源属性字段
    * `snapshot_group_name` - snapshot group name
    * `status` - 代表资源状态的资源属性字段
