---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_dbclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-adb-dbclusters"
description: |- 
  查询adb数据库集群集群
---

# alibabacloudstack_adb_dbclusters
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_adb_clusters` `alibabacloudstack_adb_db_clusters`

根据指定过滤条件列出当前凭证权限可以访问的adb产品的数据库集群列表。

## 示例用法

```terraform
# 创建VPC资源
resource "alibabacloudstack_vpc" "default" {
  name        = "${var.name}"
  cidr_block  = "172.16.0.0/16"
}

# 查询可用区
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "ADB"
}

# 查询交换机
data "alibabacloudstack_vswitches" "default" {
  vpc_id      = "${alibabacloudstack_vpc.default.id}"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

# 创建交换机
resource "alibabacloudstack_vswitch" "default" {
  name                = "tf_testAccAdb_vpc"
  vpc_id              = "${alibabacloudstack_vpc.default.id}"
  availability_zone   = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cidr_block          = "172.16.0.0/24"
}

# 定义变量
variable "creation" {	
  default = "ADB"
}

variable "name" {
  default = "tf-testAccADBConfig_18550"
}

# 创建ADB集群
resource "alibabacloudstack_adb_db_cluster" "default" {
  db_cluster_category = "Basic"
  db_cluster_class    = "C8"
  db_node_storage     = "200"
  db_cluster_version  = "3.0"
  db_node_count       = "2"
  vswitch_id          = "${alibabacloudstack_vswitch.default.id}"
  description         = "${var.name}"
  mode                = "reserver"
  cluster_type        = "analyticdb"
  cpu_type            = "intel"
  security_ips        = ["10.168.1.12", "10.168.1.11"]
}

# 查询ADB集群
data "alibabacloudstack_adb_db_clusters" "default" {	
  enable_details      = true
  description_regex   = "${alibabacloudstack_adb_db_cluster.default.description}"
}

# 输出第一个ADB集群的ID
output "first_adb_dbcluster_id" {
  value = data.alibabacloudstack_adb_db_clusters.default.clusters.0.db_cluster_id
}
```

## 参数参考

以下参数是支持的：

* `description` - (可选，变更时重建) DBCluster 的描述。这可以用来通过其描述过滤集群。
* `description_regex` - (可选，变更时重建) 用于按 DBCluster 描述过滤结果的正则表达式字符串。
* `enable_details` - (可选) 默认为 `false`。将其设置为 `true` 以输出更多关于资源属性的详细信息。
* `ids` - (可选，变更时重建) DBCluster ID 列表。这可以用来通过其唯一标识符过滤集群。
* `resource_group_id` - (可选，变更时重建) 资源组的 ID。这可以用来过滤属于特定资源组的集群。
* `status` - (可选，变更时重建) 资源的状态。有效值包括 `Creating`、`Running`、`Stopping`、`Stopped` 和 `Starting`。
* `tags` - (可选) 分配给集群的标签映射。这可以用来通过其标签过滤集群。

## 属性参考

除了上述参数外，还导出以下属性：

* `descriptions` - DBCluster 描述列表。
* `clusters` - ADB DbClusters 列表。每个元素包含以下属性：
  * `commodity_code` - 与 DBCluster 关联的服务名称。
  * `connection_string` - 集群的端点。
  * `create_time` - DBCluster 的创建时间。
  * `db_cluster_category` - DBCluster 的类别。
  * `db_cluster_id` - DBCluster 的唯一标识符。
  * `db_cluster_network_type` - DBCluster 的网络类型。
  * `network_type` - DBCluster 的网络类型。
  * `db_cluster_type` - DBCluster 的类型。
  * `db_cluster_version` - DBCluster 的版本。
  * `db_node_class` - DBCluster 中 DB 节点的类。
  * `db_node_count` - DBCluster 中的 DB 节点数。
  * `db_node_storage` - DBCluster 中每个 DB 节点的存储大小。
  * `description` - DBCluster 的描述。
  * `disk_type` - DBCluster 使用的磁盘类型。
  * `dts_job_id` - 数据传输服务 (DTS) 中的数据同步任务的 ID。此参数仅对分析实例有效。
  * `elastic_io_resource` - 分配给 DBCluster 的弹性 I/O 资源。
  * `engine` - DBCluster 使用的数据库引擎。
  * `executor_count` - DBCluster 中的执行节点数。这些节点在弹性模式下用于数据计算。
  * `expire_time` - DBCluster 的过期时间。
  * `expired` - 表示 DBCluster 是否已过期。
  * `id` - DBCluster 的 ID。
  * `lock_mode` - DBCluster 的锁定模式。可能值包括：`Unlock`(正常)、`ManualLock`(手动触发锁定)、`LockByExpiration`(实例过期自动锁定)、`LockByRestoration`(实例回滚前的自动锁定)、`LockByDiskQuota`(实例空间满自动锁定)。
  * `lock_reason` - DBCluster 被锁定的原因。
  * `maintain_time` - DBCluster 的维护窗口。
  * `payment_type` - DBCluster 的支付类型。
  * `charge_type` - DBCluster 的计费类型。
  * `port` - 用于访问 DBCluster 的端口。
  * `rds_instance_id` - 数据从中同步到 DBCluster 的 ApsaraDB RDS 实例的 ID。此参数仅对分析实例有效。
  * `resource_group_id` - DBCluster 所属的资源组的 ID。
  * `security_ips` - 允许访问 DBCluster 所有数据库的 IP 地址列表。
  * `status` - DBCluster 的状态。
  * `storage_resource` - 弹性模式下的存储资源规格。增加这些资源可以提高 DBCluster 的读写性能。
  * `tags` - 分配给 DBCluster 的标签。
  * `vpc_cloud_instance_id` - 与 DBCluster 关联的 VPC 云实例 ID。
  * `vpc_id` - 与 DBCluster 关联的 VPC ID。
  * `vswitch_id` - 与 DBCluster 关联的交换机 ID。
  * `zone_id` - DBCluster 的可用区 ID。
  * `region_id` - DBCluster 的区域 ID。