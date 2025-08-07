---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_polardbxinstance"
sidebar_current: "docs-Alibabacloudstack-drds-polardbxinstance"
description: |-
  Provides a drds Polardbxinstance resource.
---

# alibabacloudstack\_drds\_polardbxinstance

Provides a drds Polardbxinstance resource.

## 示例用法
```
variable "name" {
		default = "tf_acc_drds_polardb_16807"
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

resource "alibabacloudstack_drds_polardbx_instance" "default" {
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  cn_node_count = "2"
  dn_node_class = "mysql.n4.medium.25"
  dn_node_count = "2"
  description = "testtf1111"
  storage = "50"
  cn_node_class = "polarx.x4.medium.2e"
  engine_version = "5.7"
  network_type = "vpc"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}
```

## 参数参考

支持以下参数：
  * `series` - (选填) - 系列。
  * `cpu_type` - (选填) - cpu类型。
  * `storage` - (必填) - 存储容量。
  * `spec_series` - (选填) - 共享系列。
  * `cidr_block` - (选填) - 网段。
  * `description` - (选填) - 描述。
  * `cn_node_class` - (必填) -  计算节点规格。
  * `cn_node_count` - (必填) - 计算节点个数。
  * `create_time` - (选填) - 代表创建时间的资 实例节点源属性字段
  * `db_node_class` - (选填) - 节点规格：- **polarx.x4.medium.2e**：2核8G- **polarx.x4.large.2e**：4核16G- **polarx.x8.large.2e**：4核32G- **polarx.x4.xlarge.2e**：8核32G- **polarx.x8.xlarge.2e**：8核64G- **polarx.x4.2xlarge.2e**：16核64G- **polarx.x8.2xlarge.2e**：16核128G- **polarx.x4.4xlarge.2e**：32核128G- **polarx.x8.4xlarge.2e**：32核256G- **polarx.st.8xlarge.2e**：60核470G- **polarx.st.12xlarge.2e**：90核720G
  * `db_node_count` - (选填) - 实例节点数量，最小为2。
  * `dn_node_class` - (必填) - 实例节点规格。
  * `dn_node_count` - (必填) - 存储节点个数。
  * `engine_version` - (选填) - 5.7
  * `is_read_db_instance` - (选填) - 是否是只读实例。
  * `network_type` - (选填) - 网络类型，仅支持VPC网络。
  * `polardbx_instance_id` - (选填) - polardbx 实例ID.
  * `primary_db_instance_id` - (选填) - 主实例ID。
  * `primary_zone` - (选填) - 主可用区。
  * `resource_type` - (选填) - 资源类型，目前仅支持PolarDB-X 2.0实例一种类型的资源。
  * `secondary_zone` - (选填) - 次可用区。
  * `tertiary_zone` - (选填) - 第三可用区。
  * `topology_type` - (选填) - 拓扑类型：- **3azones**：三可用区；- **1azone**：单可用区。
  * `vswitch_id` - (必填) - 交换机ID.
  * `vpc_id` - (必填) - VPC ID。
  * `zone_id` - (选填) - 可用区。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `series` - 系列。
  * `cpu_type` - cpu类型。
  * `spec_series` - 共享类型。
  * `description` - 描述。
  * `create_time` - 代表创建时间的资源属性字段
  * `db_node_class` - 节点规格：- **polarx.x4.medium.2e**：2核8G- **polarx.x4.large.2e**：4核16G- **polarx.x8.large.2e**：4核32G- **polarx.x4.xlarge.2e**：8核32G- **polarx.x8.xlarge.2e**：8核64G- **polarx.x4.2xlarge.2e**：16核64G- **polarx.x8.2xlarge.2e**：16核128G- **polarx.x4.4xlarge.2e**：32核128G- **polarx.x8.4xlarge.2e**：32核256G- **polarx.st.8xlarge.2e**：60核470G- **polarx.st.12xlarge.2e**：90核720G
  * `db_node_count` - 实例节点数量，最小为2。
  * `status` - 代表资源状态的资源属性字段
