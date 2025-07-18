---
subcategory: "块存储 EBS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicagroups"
sidebar_current: "docs-Alibabacloudstack-datasource-ebs-diskreplicagroups"
description: |-
  提供阿里云账号下拥有的ebs diskreplicagroups列表。
---

# alibabacloudstack\_ebs\_diskreplicagroups

此数据源提供根据指定过滤条件列出的阿里云账号下的ebs diskreplicagroups资源列表。

## 示例用法
```
variable "name" {
	default = "tf-testAccEbsDiskReplicaGroupsDataSource-8901034"
}

variable "region_id" {
  default = ""
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}



resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
    disk_replica_group_name = "${var.name}"
    description = "${var.name}"
    destination_region_id = "${var.region_id}"
    destination_zone_id ="${data.alibabacloudstack_zones.default.zones[1].id}"
    site = "production"
    source_region_id = "${var.region_id}"
    source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
}
 

data "alibabacloudstack_ebs_diskreplicagroups" "default" {
  description_regex = "${alibabacloudstack_ebs_diskreplicagroup.default.description}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 用于通过一致性复制组ID过滤结果。
  * `name_regex` - (选填) - 用于通过一致性复制组名称过滤结果。
  * `description_regex` - (选填) - 用于通过一致性复制组描述信息过滤结果。
  * `site` - (选填) - 复制对和一致性复制组的站点信息来源。可能值：- production：生产站点。- backup：灾备站点。
  * `source_region_id` - (选填) - 生产站点所属的地域ID。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `disk_replica_groups` - 一致性复制组列表。每个元素包含以下属性：
    * `id` - 一致性复制组的ID。
    * `description` - 一致性复制组的描述信息。
    * `destination_region_id` - 灾备站点所属的地域ID。
    * `destination_zone_id` - 灾备站点所属的可用区ID。
    * `disk_replica_group_name` - 一致性复制组名称。
    * `last_recover_point` - 一致性复制组的最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
    * `one_shot` - 是否立刻进行一次同步。取值范围：- true：立刻开始数据同步。- false：在RPO时间周期之后才开始数据同步。默认值：false。
    * `pair_ids` - 一致性复制组中包含的复制对ID列表。
    * `pair_number` - 一致性复制组中包含的复制对个数。
    * `rpo` - 一致性复制组的Rpo值
    * `replica_group_id` - 一致性复制组ID。
    * `site` - 复制对和一致性复制组的站点信息来源。可能值：- production：生产站点。- backup：灾备站点。
    * `source_region_id` - 生产站点所属的地域ID。
    * `source_zone_id` - 生产站点所属的可用区ID。
    * `status` - 一致性复制组的状态。可能值：- invalid：失效。该状态表示一致性复制组中复制对存在异常。- creating：创建中。- created：已创建。- create_failed：创建失败。- manual_syncing：单次同步中。如果是第一次单次同步，则同步中也显示为该状态。- syncing：同步中。主盘和从盘之间非第一次进行一致性复制数据时，将处于该状态。- normal：正常。当一致性复制的当前周期内数据复制完成时，将处于该状态。- stopping：停止中。- stopped：已停止。- stop_failed：停止失败。- failovering：故障切换中。- failovered：故障切换完成。- failover_failed：故障切换失败。- reprotecting：反向复制操作中。- reprotect_failed：反向复制失败。- deleting：删除中。- delete_failed：删除失败。- deleted：已删除。
    * `tags` - 代表资源标签的资源属性字段
