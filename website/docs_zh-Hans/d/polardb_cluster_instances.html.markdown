---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instances"
sidebar_current: "docs-alibabacloudstack-datasource-polardb-cluster-instances"
description: |-
  提供 PolarDB 集群实例列表。
---

# alibabacloudstack_polardb_cluster_instances

本数据源提供阿里云平台环境中的 PolarDB 集群实例列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-polar-test"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name      = "${var.name}_vsw"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description 	=  "${var.name}"
  zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_type 				= "${var.db_type}"
  db_version 			= "${var.db_version}"
  storage_space 		= "20"
  vpc_id 				= "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id			= "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class 		= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type 			= "ESSDPL1"
}
data "alibabacloudstack_polardb_cluster_instances" "default" {
  description_regex = "test"
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 实例 ID 列表。
* `name_regex` - (可选, 已弃用) 用于通过实例描述过滤结果的正则表达式。**字段 'name_regex' 已弃用，并将在未来版本中移除。请使用新字段 'description_regex' 替代。**
* `description_regex` - (可选) 用于通过实例描述过滤结果的正则表达式。

## 属性参考

导出以下属性：

* `ids` - 实例 ID 列表。
* `db_cluster_instances` - PolarDB 集群实例列表。每个元素包含以下属性：
  * `deletion_lock` - 删除锁定状态。
  * `category` - 集群实例的类别。
  * `resource_group_id` - 资源组的 ID。
  * `storage_pay_type` - 存储付费类型。
  * `hot_standby_cluster` - 热备集群状态。
  * `db_cluster_id` - 集群实例的 ID。
  * `db_type` - 数据库类型。
  * `db_cluster_network_type` - 集群实例的网络类型。
  * `db_proxy_cluster_class` - 代理集群类别。
  * `db_version` - 数据库版本。
  * `department` - 部门信息。
  * `db_nodes` - 数据库节点信息。
    * `zone_id` - 节点的可用区 ID。
    * `db_node_role` - 节点的角色。
    * `db_node_id` - 节点的 ID。
    * `region_id` - 节点的区域 ID。
    * `db_node_class` - 节点的类别。
  * `engine` - 数据库引擎。
  * `tags` - 集群实例的标签。
    * `tag_key` - 标签键。
    * `tag_value` - 标签值。
  * `resource_group` - 资源组信息。
  * `architecture` - 集群实例的架构。
  * `zone_id` - 集群实例的可用区 ID。
  * `db_cluster_status` - 集群实例的状态。
  * `create_time` - 集群实例的创建时间。
  * `db_cluster_description` - 集群实例的描述。
  * `expired` - 集群实例是否已过期。
  * `pay_type` - 集群实例的付费类型。
  * `lock_mode` - 集群实例的锁定模式。
  * `vswitch_id` - 交换机的 ID。
  * `ascm_create_user` - ASCM 创建用户。
  * `db_node_class` - 集群实例的节点类别。
  * `storage_used` - 已使用的存储空间。
  * `db_node_number` - 数据库节点数量。
  * `vpc_id` - 专有网络的 ID。
  * `storage_space` - 集群实例的存储空间。
  * `serverless_type` - Serverless 类型。
  * `department_name` - 部门名称。
  * `region_id` - 集群实例的区域 ID。
  * `expire_time` - 集群实例的过期时间。
  * `vip` - 集群实例的 VIP。
  * `resource_group_name` - 资源组的名称。