---
subcategory: "块存储"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicapairs"
sidebar_current: "docs-Alibabacloudstack-datasource-ebs-diskreplicapairs"
description: |-
  提供阿里云账号下拥有的ebs diskreplicapairs列表。
---

# alibabacloudstack\_ebs\_diskreplicapairs

此数据源提供根据指定过滤条件列出的阿里云账号下的ebs diskreplicapairs资源列表。

## 示例用法
```
variable "name" {
	default = "tf-testAccEbsDiskReplicaPairsDataSource-2021428"
}

variable "region" {
  default = ""
}

resource "alibabacloudstack_ecs_disk" "disk1" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
}

resource "alibabacloudstack_ecs_disk" "disk2" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
}
	
resource "alibabacloudstack_ebs_diskreplicapair" "default" {
	disk_replica_pair_name = "${var.name}"
	description =            "${var.name}"
	source_zone_id =         "cn-wulan-env82-amtest83002-b"
	source_region_id =       "${var.region}"
	source_disk_id =         "d-9rt00vs0qsrd2502tx5u"
	destination_zone_id =    "cn-wulan-env82-amtest82001-a"
	destination_region_id =  "${var.region}"
	destination_disk_id =    "d-9rt00vs0qsrd2502tx5p"
	rpo =                    300
}
 

data "alibabacloudstack_ebs_diskreplicapairs" "default" {
  description_regex = "${alibabacloudstack_ebs_diskreplicapair.default.description}"
}
```

## 参数参考
以下参数是支持的：
  * `name_regex` - (选填) - 用于通过异步复制关系名称过滤结果。
  * `description_regex` - (选填) - 用于通过异步复制关系描述信息过滤结果。
  * `ids` - (选填) - 用于通过异步复制关系ID过滤结果。
  * `source_region_id` - (选填) - 生产站点所属的地域ID。
  * `replica_group_id` - (选填) - 复制组ID。
  * `max_results` - (选填) - MaxItems本次请求所返回的最大记录条数。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `disk_replica_pairs` - 异步部异步复制关系列表。每个元素包含以下属性：
    * `id` - 异步择异步复制关系ID。
    * `bandwidth` - 云盘之间异步复制数据时的带宽。单位为Kbps。取值范围：- 10240 Kbps：等于10 Mbps。- 20480 Kbps：等于 20 Mbps。- 51200 Kbps：等于50 Mbps。- 102400 Kbps：等于100 Mbps。默认值：10240。当ChargeType取值为PayAsYouGo时，不能指定本参数值，系统取值为0，表示云盘异步复制时根据数据写入变化动态分配。
    * `description` - 异步复制关系的描述信息。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。
    * `destination_disk_id` - 目标云盘（从盘）的云盘ID。
    * `destination_region_id` - 灾备站点所属的地域ID。
    * `destination_zone_id` - 灾备站点所属的可用区ID。
    * `source_disk_id` - 生产站点的磁盘ID。
    * `disk_replica_pair_name` - 异步复制关系的名称。长度为2~128个字符，必须以大小字母或中文开头，不能以`http://`或`https://`开头。可以包含中文、英文、数字、半角冒号（:）、下划线（_）、半角句号（.）或者短划线（-）。
    * `last_recover_point` - 异步复制关系最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
    * `one_shot` - 是否立刻进行一次同步。取值范围：- true：立刻开始数据同步。- false：在RPO时间周期之后才开始数据同步。默认值：false。
    * `rpo` - 异步复制关系的Rpo。单位为秒。取值范围：300到86400。默认值：300。
    * `source_region_id` - 生产站点所属的地域ID。
    * `replica_group_id` - 复制组ID。
    * `replica_pair_id` - 代表资源一级ID的资源属性字段
    * `source_zone_id` - 生产站点所属的可用区ID。
    * `status` - 代表资源状态的资源属性字段
    * `tags` - 代表资源标签的资源属性字段
