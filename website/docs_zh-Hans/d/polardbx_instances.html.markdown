---
subcategory: "POLARDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-instances"
description: |-
  提供阿里云账号下拥有的 polardbx instances列表。
---

# alibabacloudstack\_polardbx\_instances

此数据源提供根据指定过滤条件列出的阿里云账号下的polardbx instances资源列表。

## 示例用法
```
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
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

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
	topology_type = "1azone"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

data "alibabacloudstack_polardbx_instances" "default" {
  description_regex = "${alibabacloudstack_polardbx_instance.default.description}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 实例的id列表。
  * `name_regex` - (选填) - 用名称对结果进行过滤的正则表达式。
  * `description_regex` - (选填) - 用描述对结果进行过滤的正则表达式。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `polardbx_instances` - polardbx 实力列表，支持以下属性
    * `id` - 实例ID。
    * `series` - 系列。
    * `storage` - 存储容量。
    * `description` - 实例描述。
    * `cpu_type` - cpu类型。
    * `cn_node_class` - 计算节点规格。
    * `cn_node_count` - 计算节点个数。
    * `create_time` - 代表创建时间的资源属性字段
    * `db_node_class` - 实例节点规格.
    * `db_node_count` - 实例节点数量，最小为2。
    * `dn_node_class` - 存储节点规格。
    * `dn_node_count` - 存储节点个数。
    * `engine_version` - 引擎版本。 取值：`5.7`、`8.0`。
    * `is_read_db_instance` - 是否是只读节点。
    * `network_type` - 网络类型，仅支持VPC网络。
    * `polardbx_instance_id` - 实例id.
    * `primary_zone` - 主可用区。
    * `resource_type` - 资源类型，目前仅支持PolarDB-X 2.0实例一种类型的资源。
    * `secondary_zone` - 次可用区。
    * `status` - 代表资源状态的资源属性字段
    * `tertiary_zone` - 第三可用区。
    * `topology_type` - 拓扑类型：- **3azones**：三可用区；- **1azone**：单可用区。
    * `vswitch_id` - VSwitch 交换机 ID。
    * `vpc_id` - VPC ID。
    * `zone` - 实例可用区。
