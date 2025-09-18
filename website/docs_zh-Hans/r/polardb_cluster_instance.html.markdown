---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instance"
sidebar_current: "docs-alibabacloudstack-resource-polardb-cluster-instance"
description: |-
  提供 PolarDB 集群实例资源。
---

# alibabacloudstack_polardb_cluster_instance

提供 PolarDB 集群实例资源。PolarDB 集群是一种具有云原生架构的分布式关系型数据库服务。

## 示例用法

### 创建 PolarDB MySQL 集群实例

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
```

## 参数参考

支持以下参数：

* `db_type` - (必填, 变更后新建) 数据库引擎类型。有效值：[MySQL](file://d:\terraform\terraform-provider-apsarastack\alibabacloudstack\extension_rds.go#L5-L5)、`PolarDB_PG`、`PolarDB_PPAS`。
* `db_version` - (必填, 变更后新建) 数据库引擎版本。
* `db_node_class` - (必填) 集群实例的节点类别。
* `zone_id` - (必填, 变更后新建) 集群实例的可用区 ID。
* `vpc_id` - (必填, 变更后新建) 专有网络 ID。
* `vswitch_id` - (必填, 变更后新建) 交换机 ID。
* `db_cluster_description` - (可选) 集群实例的描述。
* `storage_type` - (可选, 变更后新建) 集群实例的存储类型。
* `storage_space` - (可选, 变更后新建) 集群实例的存储空间。
* `readonly_node_num` - (可选) 只读节点数量。有效值：0 到 15。
* `hot_standby_cluster` - (可选) 指定是否启用热备集群。有效值：`standby`、`off`。
* `cpu_type` - (可选, 变更后新建) 集群实例的 CPU 类型。有效值：`intel`、`hygon`。
* `sub_category` - (可选) 集群实例的子类别。有效值：`normal_general`、`normal_exclusive`。
* `deletion_lock` - (可选) 集群实例的删除锁定。有效值：0、1。
* `maintain_time` - (可选) 集群实例的维护时间。
* `security_ips_groups` - (可选) 集群实例的安全 IP 组。
* `security_groups` - (可选) 集群实例的安全组。
* `ssl_enabled` - (可选) 是否启用 SSL 加密。
* `tde_enabled` - (可选) 是否启用透明数据加密 (TDE)。
* `encryption_key` - (可选) KMS 密钥的 ID。
* `encrypt_algorithm` - (可选) 加密算法。
* `parameters` - (可选) 集群实例的参数。
* `tags` - (可选) 分配给资源的标签映射。

## 属性参考

导出以下属性：

* `id` - 集群实例的 ID。
* `db_cluster_id` - 集群实例的 ID。
* `category` - 集群实例的类别。
* `db_cluster_network_type` - 集群实例的网络类型。
* `storage_max` - 集群实例的最大存储空间。
* `engine` - 数据库引擎。
* `zone_ids` - 集群实例的可用区 ID。
* `db_nodes` - 数据库节点信息。
  * `db_node_status` - 节点的状态。
  * `zone_id` - 节点的可用区 ID。
  * `max_connections` - 节点的最大连接数。
  * `added_cpu_cores` - 节点添加的 CPU 核心数。
  * `db_node_role` - 节点的角色。
  * `imci_switch` - 节点的 IMCI 开关。
  * `db_node_id` - 节点的 ID。
  * `max_iops` - 节点的最大 IOPS。
  * `db_node_class` - 节点的类别。
  * `creation_time` - 节点的创建时间。
  * `scc_mode` - 节点的 SCC 模式。
  * `failover_priority` - 节点的故障转移优先级。
  * `hot_replica_mode` - 节点的热副本模式。
  * `server_weight` - 节点的服务器权重。
* `proxy_status` - 集群实例的代理状态。
* `architecture` - 集群实例的架构。
* `db_cluster_status` - 集群实例的状态。
* `vip` - 集群实例的 VIP。
* `lock_mode` - 集群实例的锁定模式。
* `creation_time` - 集群实例的创建时间。
* `sql_size` - 集群实例的 SQL 大小。
* `proxy_cpu_cores` - 集群实例的代理 CPU 核心数。
* `pay_type` - 集群实例的付费类型。

## 导入

可以使用 id 导入 PolarDB 集群实例，例如：

```bash
$ terraform import alibabacloudstack_polardb_cluster_instance.example pc-12345678
```