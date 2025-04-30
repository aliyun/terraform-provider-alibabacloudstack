---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_dbcluster"
sidebar_current: "docs-Alibabacloudstack-adb-dbcluster"
description: |-
  编排adb数据库集群
---

# alibabacloudstack_adb_dbcluster
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_adb_cluster` `alibabacloudstack_adb_db_cluster`

使用Provider配置的凭证在指定的资源集下编排adb数据库集群。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccadbCluster73485"
}

data "alibabacloudstack_ascm_resource_groups" "default" {
  name_regex = ""
}

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "172.16.0.0/16"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "ADB"
}

data "alibabacloudstack_vswitches" "default" {
  vpc_id   = alibabacloudstack_vpc.default.id
  zone_id  = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_vswitch" "default" {
  name                 = "tf_testAccAdb_vpc"
  vpc_id               = alibabacloudstack_vpc.default.id
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  cidr_block          = "172.16.0.0/24"
}

resource "alibabacloudstack_adb_db_cluster" "default" {
  vswitch_id             = alibabacloudstack_vswitch.default.id
  db_cluster_category    = "Basic"
  description           = var.name
  db_node_storage       = "200"
  mode                  = "reserver"
  cpu_type              = "intel"
  db_cluster_version    = "3.0"
  db_node_count         = "2"
  db_node_class         = "C8"
  cluster_type          = "analyticdb"
  payment_type          = "PayAsYouGo"
  maintain_time         = "23:00Z-00:00Z"
  security_ips          = ["10.168.1.12", "10.168.1.11"]
}
```

## 参数参考

支持以下参数：

* `auto_renew_period` - (选填) - 集群的自动续费周期，单位为月。当 `payment_type` 为 `Subscription` 时有效。有效值：`1`, `2`, `3`, `6`, `12`, `24`, `36`。默认为 `1`。
* `compute_resource` - (选填) - 弹性模式下使用的计算资源规格，用于数据计算。增加资源可以加速查询。更多信息，请参见 [ComputeResource](https://www.alibabacloud.com/help/en/doc-detail/144851.htm)。
* `db_cluster_category` - (必填) - 数据库集群类别。有效值：`Basic`, `Cluster`, `MixedStorage`。
* `db_cluster_class` - (选填) - 与属性 `db_node_class` 重复。
* `storage_resource` - (选填) - 弹性模式下的存储资源规格。这些资源用于数据读写操作。增加资源可以提高集群的读写性能。更多信息，请参见 [Specifications](https://www.alibabacloud.com/help/en/doc-detail/144851.htm)。
* `storage_type` - (选填) - 预留参数，不涉及。
* `db_cluster_version` - (选填, 变更时重建) - 数据库集群版本。选项值：`3.0`。默认为 `3.0`。
* `cluster_type` - (必填) - 资源的集群类型。有效值：`analyticdb`, `AnalyticdbOnPanguHybrid`。
* `cpu_type` - (必填) - 资源的 CPU 类型。有效值：`intel`。
* `db_node_class` - (选填) - 数据库节点类。更多信息，请参见 [DBClusterClass](https://help.aliyun.com/document_detail/190519.html)。
* `db_node_count` - (选填) - 数据库节点的数量。
* `executor_count` - (选填) - 弹性模式下用的计算资源规格，对应后台实际的节点数。
* `db_node_storage` - (选填) - 每个数据库节点的存储容量。
* `description` - (选填) - DBCluster 的描述。
* `maintain_time` - (选填) - 集群实例可维护时间段。格式：`hh:mmZ-hh:mmZ`。
* `mode` - (必填) - 模式。取值说明：`reserver`：预留模式。`flexible`：弹性模式。
* `modify_type` - (选填) - 修改类型。
* `payment_type` - (选填, 变更时重建) - 代表付费类型的资源属性字段。有效值为 `PayAsYouGo` 和 `Subscription`。默认为 `PayAsYouGo`。
* `pay_type` - (选填, 变更时重建) - 已废弃字段。使用 `payment_type` 替代。
* `period` - (选填) - 指定预付费集群为包年或包月类型。当 `payment_type` 为 `Subscription` 时有效。有效值：[1~9], 12, 24, 36。
* `renewal_status` - (选填) - 有效值为 `AutoRenewal`, `Normal`, `NotRenewal`。默认为 `NotRenewal`。
* `resource_group_id` - (选填) - 代表资源组的资源属性字段。
* `security_ips` - (选填) - 允许访问所有数据库集群的 IP 地址列表。列表最多包含 1,000 个 IP 地址，以逗号分隔。支持的格式包括 `0.0.0.0/0`, `10.23.12.24`(IP)和 `10.23.12.24/24`(CIDR 模式。`/24` 表示 IP 地址前缀长度。前缀长度范围为 `[1,32]`)。
* `vswitch_id` - (选填, 变更时重建) - 交换机 ID。
* `zone_id` - (选填, 变更时重建) - 资源所在的可用区 ID。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `storage_resource` - 弹性模式下的存储资源规格。这些资源用于数据读写操作。增加资源可以提高集群的读写性能。
* `storage_type` - 预留参数，不涉及。
* `db_node_class` - 数据库节点类。
* `db_node_count` - 数据库节点的数量。
* `executor_count` - 弹性模式下用的计算资源规格，对应后台实际的节点数。
* `db_node_storage` - 每个数据库节点的存储容量。
* `description` - DBCluster 的描述。
* `maintain_time` - 集群实例可维护时间段。
* `payment_type` - 代表付费类型的资源属性字段。
* `pay_type` - 取值说明：`Postpaid`：按量付费。`Prepaid`：预付费(包年包月)。
* `resource_group_id` - 代表资源组的资源属性字段。
* `security_ips` - IP 白名单分组下的 IP 列表，最多 1000 个，以英文逗号(,)隔开。
* `status` - 代表资源状态的资源属性字段。
* `instance_inner_connection` - 集群的内部连接端点。
* `instance_inner_port` - 集群的内部端口。
* `instance_vpc_id` - VPC ID。
* `connection_string` - 集群的连接字符串。
* `port` - 集群的端口。
* `zone_id` - 资源所在的可用区 ID。

## 导入

AnalyticDB for MySQL (ADB) DBCluster 可以通过 id 导入，例如

```bash
$ terraform import alibabacloudstack_adb_db_cluster.example <id>
```