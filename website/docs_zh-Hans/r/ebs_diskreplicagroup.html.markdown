---
subcategory: "块存储 EBS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicagroup"
sidebar_current: "docs-Alibabacloudstack-resource-ebs-diskreplicagroup"
description: |-
  Provides a ebs Diskreplicagroup resource.
---

# alibabacloudstack\_ebs\_diskreplicagroup

Provides a ebs Diskreplicagroup resource.

## 示例用法
```
variable "name" {
  default = "tf-testaccebs-diskreplicagroup43963"
}
variable "region" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
  disk_replica_group_name = "${var.name}"
  description = "ebs_diskreplicagroup test"
  source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
  source_region_id = "${var.region}"
  destination_zone_id = "${data.alibabacloudstack_zones.default.zones[1].id}"
  destination_region_id = "${var.region}"
  site = "production"
  rpo = "300"
}
```

## 参数参考

支持以下参数：
  * `description` - (选填) - 一致性复制组的描述信息。
  * `destination_region_id` - (必填) - 灾备站点所属的地域ID。
  * `destination_zone_id` - (必填) - 灾备站点所属的可用区ID。
  * `disk_replica_group_name` - (选填) - 一致性复制组名称。
  * `last_recover_point` - (选填) - 一致性复制组的最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
  * `rpo` - (选填) - 一致性复制组的RPO值。该参数单位为秒。
  * `region_id` - (选填) - 一致性复制组所属的地域ID，与生产站点所属的地域相同。
  * `site` - (选填) - 复制对和一致性复制组的站点信息来源。可能值：- production：生产站点。- backup：灾备站点。
  * `source_region_id` - (必填) - 生产站点所属的地域ID。
  * `source_zone_id` - (必填) - 生产站点所属的可用区ID。
  * `status` - (选填) - 一致性复制组的状态。可能值：- invalid：失效。该状态表示一致性复制组中复制对存在异常。- creating：创建中。- created：已创建。- create_failed：创建失败。- manual_syncing：单次同步中。如果是第一次单次同步，则同步中也显示为该状态。- syncing：同步中。主盘和从盘之间非第一次进行异步复制数据时，将处于该状态。- normal：正常。当异步复制的当前周期内数据复制完成时，将处于该状态。- stopping：停止中。- stopped：已停止。- stop_failed：停止失败。- failovering：故障切换中。- failovered：故障切换完成。- failover_failed：故障切换失败。- reprotecting：反向复制操作中。- reprotect_failed：反向复制失败。- deleting：删除中。- delete_failed：删除失败。- deleted：已删除。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `last_recover_point` - 一致性复制组的最近一次异步复制操作完成的时间。该参数以时间戳的形式提供返回值。单位：秒。
  * `pair_ids` - 一致性复制组中包含的复制对ID列表。
  * `pair_number` - 一致性复制组中包含的复制对个数。
  * `rpo` - 一致性复制组的RPO值。该参数以秒（s）为单位。
  * `replica_group_id` - 一致性复制组ID。
  * `status` - 一致性复制组的状态。可能值：- invalid：失效。该状态表示一致性复制组中复制对存在异常。- creating：创建中。- created：已创建。- create_failed：创建失败。- manual_syncing：单次同步中。如果是第一次单次同步，则同步中也显示为该状态。- syncing：同步中。主盘和从盘之间非第一次进行异步复制数据时，将处于该状态。- normal：正常。当异步复制的当前周期内数据复制完成时，将处于该状态。- stopping：停止中。- stopped：已停止。- stop_failed：停止失败。- failovering：故障切换中。- failovered：故障切换完成。- failover_failed：故障切换失败。- reprotecting：反向复制操作中。- reprotect_failed：反向复制失败。- deleting：删除中。- delete_failed：删除失败。- deleted：已删除。
