---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_diskattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-diskattachment"
description: |- 
  编排云绑定服务器（Ecs）磁盘和实例
---

# alibabacloudstack_ecs_diskattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_disk_attachment`

使用Provider配置的凭证在指定的资源集下编排云绑定服务器（Ecs）磁盘和实例。

## 示例用法

```hcl
variable "name" {
	default = "tf-testAccEcsDiskAttachmentConfig"
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
	cidr_block = "172.16.0.0/24"
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
	cidr_ip = "172.16.0.0/24"
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
	zone_id    = data.alibabacloudstack_zones.default.zones.0.id
	is_outdated          = false
	lifecycle {
		ignore_changes = [
			instance_type
		]
	}
}

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"

	tags = {
		Name = "TerraformTest-disk"
	}
}

resource "alibabacloudstack_ecs_diskattachment" "default" {
	disk_id = "${alibabacloudstack_ecs_disk.default.id}"
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
	device_name = "/dev/xvdb"
}
```

## 参数参考

支持以下参数：
  * `instance_id` - (必填, 变更时重建) - 目标 ECS 实例的 ID。
  * `disk_id` - (必填, 变更时重建) - 待挂载的云盘 ID。云盘(`DiskId`)和实例(`InstanceId`)必须在同一个可用区中。支持挂载数据盘和系统盘，相关约束条件请参见上文接口说明章节。
  * `device_name` - (选填, 变更时重建) - 暴露给实例的设备名称。它将由系统根据默认顺序从 `/dev/xvdb` 到 `/dev/xvdz` 自动分配。如果指定，则必须与实例上的可用设备名称之一匹配。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `instance_id` - 实例的 ID。
  * `disk_id` - 磁盘的 ID。
  * `device_name` - 暴露给实例的设备名称。这是实际分配给实例的设备名称，可能与配置中的 `device_name` 参数一致或由系统自动分配。