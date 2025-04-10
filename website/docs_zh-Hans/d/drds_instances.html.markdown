---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-instances"
description: |- 
  查询云原生分布式数据库（Drds）实例
---

# alibabacloudstack_drds_instances

根据指定过滤条件列出当前凭证权限可以访问的云原生分布式数据库（Drds）实例列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

variable "name" {
	default = "tf-testAccDRDSInstancesDataSource-8710705"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = "${alibabacloudstack_vpc.default.id}"
	cidr_block        = "172.16.0.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name              = "${var.name}"
}	

resource "alibabacloudstack_drds_instance" "default" {
	description = "${var.name}"
	zone_id = "${alibabacloudstack_vswitch.default.availability_zone}"
	instance_series = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	specification = "drds.sn2.4c16g.8C32G"
}

data "alibabacloudstack_drds_instances" "default" {
  name_regex      = "${alibabacloudstack_drds_instance.default.description}"
  description_regex = "example-description.*"
  ids             = ["${alibabacloudstack_drds_instance.default.id}"]
}

output "first_db_instance_id" {
  value = "${data.alibabacloudstack_drds_instances.default.instances.0.id}"
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (可选) 用于按实例名称筛选结果的正则表达式字符串。通过该参数，可以匹配符合特定命名规则的 DRDS 实例。
* `description_regex` - (可选) 用于按实例描述筛选结果的正则表达式字符串。通过该参数，可以匹配符合特定描述规则的 DRDS 实例。
* `ids` - (可选) DRDS 实例 ID 列表。通过该参数，可以限制查询结果为指定的 DRDS 实例。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - DRDS 实例 ID 列表。该列表包含所有符合条件的 DRDS 实例的 ID。
* `descriptions` - DRDS 实例描述列表。该列表包含所有符合条件的 DRDS 实例的描述。
* `instances` - DRDS 实例列表。每个实例包含以下属性：
  * `id` - DRDS 实例的唯一标识符。
  * `description` - DRDS 实例的描述信息。
  * `status` - DRDS 实例的当前状态，例如“Running”表示运行中，“Stopped”表示已停止。
  * `type` - DRDS 实例的类型或规格，表示实例的计算资源配置。
  * `create_time` - DRDS 实例的创建时间，格式为标准时间戳（ISO 8601 格式）。
  * `network_type` - DRDS 实例的网络类型，`Classic` 表示公共经典网络，`VPC` 表示私有网络。
  * `zone_id` - DRDS 实例所在的可用区 ID。
  * `version` - DRDS 实例的功能版本号，表示实例支持的功能集。