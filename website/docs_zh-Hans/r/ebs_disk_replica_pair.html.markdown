---
subcategory: "Elastic Block Storage (EBS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicapair"
sidebar_current: "docs-Alibabacloudstack-ebs-diskreplicapair"
description: |-
  Provides a ebs Diskreplicapair resource.
---

# alibabacloudstack\_ebs\_diskreplicapair

Provides a ebs Diskreplicapair resource.

## 示例用法
```
variable "name" {
  default = "tf-testaccebs-diskreplicapair59477"
}

variable "region" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_ecs_disk" "disk1" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"

	lifecycle {
		ignore_changes = ["tags"]
	}
}

resource "alibabacloudstack_ecs_disk" "disk2" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
	lifecycle {
		ignore_changes = ["tags"]
	}
}

resource "alibabacloudstack_ebs_diskreplicapair" "default" {
  description = "ebs_diskreplicapair test"
  source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
  source_region_id = "${var.region}"
  source_disk_id = "${alibabacloudstack_ecs_disk.disk1.id}"
  destination_zone_id = "${data.alibabacloudstack_zones.default.zones[1].id}"
  destination_disk_id = "${alibabacloudstack_ecs_disk.disk2.id}"
  disk_replica_pair_name = "${var.name}"
  destination_region_id = "${var.region}"
  rpo = "300"
}
```

## 参数参考

支持以下参数：
  * `description` - (选填) - 异步复制关系的描述信息。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。
  * `destination_disk_id` - (必填) - 目标云盘（从盘）的云盘ID。
  * `destination_region_id` - (必填) - 灾备站点所属的地域ID。
  * `destination_zone_id` - (必填) - 灾备站点所属的可用区ID。
  * `source_disk_id` - (必填) - 生产站点的磁盘ID。
  * `disk_replica_pair_name` - (选填) - 异步复制关系的名称。长度为2~128个字符，必须以大小字母或中文开头，不能以`http://`或`https://`开头。可以包含中文、英文、数字、半角冒号（:）、下划线（_）、半角句号（.）或者短划线（-）。
  * `last_recover_point` - (选填) - 异步复制关系最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
  * `one_shot` - (选填) - 是否立刻进行一次同步。取值范围：- true：立刻开始数据同步。- false：在RPO时间周期之后才开始数据同步。默认值：false。
  * `rpo` - (选填) - 异步复制关系的Rpo。单位为秒。取值范围：300到86400。默认值：300。
  * `source_region_id` - (必填) - 生产站点的地域ID。
  * `replica_group_id` - (选填) - 复制组ID。
  * `source_zone_id` - (必填) - 生产站点所属的可用区ID。
  * `status` - (选填) - 代表资源状态的资源属性字段

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `last_recover_point` - 异步复制关系最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
  * `replica_pair_id` - 代表资源一级ID的资源属性字段
  * `create_time` - 代表创建时间的资源属性字段
  * `status` - 代表资源状态的资源属性字段
