---
subcategory: "块存储 EBS"
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
  * `source_region_id` - (必填) 生产站点（主盘）所属的地域 ID。
  * `source_zone_id` - (必填) 生产站点（主盘）所属的可用区 ID。
  * `source_disk_id` - (必填) 生产站点的云盘 ID（主盘）。
  * `destination_region_id` - (必填) 灾备站点（从盘）所属的地域 ID。
  * `destination_zone_id` - (必填) 灾备站点（从盘）所属的可用区 ID。
  * `destination_disk_id` - (必填) 灾备站点的云盘 ID（从盘）。
  * `disk_replica_pair_name` - (选填) 异步复制关系的名称。长度为 2~128 个字符，必须以字母或中文开头，不能以 `http://` 或 `https://` 开头。可以包含中文、英文、数字、半角冒号（:）、下划线（_）、半角句号（.）或短划线（-）。
  * `description` - (选填) 异步复制关系的描述信息。长度为 2~256 个字符，不能以 `http://` 或 `https://` 开头。
  * `rpo` - (选填) 异步复制关系的 RPO（恢复点目标）。单位为秒。取值范围：300 到 86400。默认值：300。
  * `one_shot` - (选填) 是否立刻进行一次同步。取值：`true`（立刻开始数据同步）、`false`（在 RPO 时间周期之后才开始数据同步）。默认值：`false`。
  * `replica_group_id` - (选填) 复制组 ID。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 异步复制关系的 ID。
  * `replica_pair_id` - 异步复制关系的 ID。
  * `create_time` - 异步复制关系的创建时间，格式为 UTC 时间（例如：2006-01-02T15:04:05-07:00）。
  * `last_recover_point` - 异步复制关系最近一次异步复制操作完成的时间戳。单位：秒。
  * `status` - 异步复制关系的状态。取值：`creating`、`created`、`syncing`、`normal`、`stopped`、`failovered`、`deleting`、`failed` 等。

## 导入

EBS 磁盘复制对可以使用复制对 ID 导入，例如：

```
$ terraform import alibabacloudstack_ebs_diskreplicapair.example rp-xxxxxxxxx
```
