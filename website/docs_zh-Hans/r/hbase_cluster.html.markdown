---
subcategory: "HBase"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hbase_cluster"
sidebar_current: "docs-Alibabacloudstack-hbase-cluster"
description: |- 
  编排云数据库hbase集群
---

# alibabacloudstack_hbase_cluster
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_hbase_instance`

使用Provider配置的凭证在指定的资源集下编排云数据库hbase集群。

## 示例用法

```hcl
variable "name" {
	default = "tf-testAccVpc1175381"
}

variable "password" {}

data "alibabacloudstack_zones" "default" {}

data "alibabacloudstack_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alibabacloudstack_vpc" "default" {
	name       = var.name
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "172.16.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.ids.0
	name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
	count   = 2
	vpc_id  = alibabacloudstack_vpc.default.id
	name    = var.name
}

resource "alibabacloudstack_hbase_instance" "default" {
	name                  = var.name
	ip_white             = "192.168.1.2"
	vswitch_id           = alibabacloudstack_vswitch.default.id
	account              = "adminu"
	zone_id              = data.alibabacloudstack_zones.default.zones.0.id
	security_groups      = [alibabacloudstack_security_group.default[0].id, alibabacloudstack_security_group.default[1].id]
	core_disk_size       = 480
	engine_version       = "2.0"
	cold_storage_size    = 900
	deletion_protection  = false
	core_instance_type   = "hbase.sn1.large"
	master_instance_type = "hbase.sn1.large"
	immediate_delete_flag = true
	maintain_start_time  = "14:00Z"
	password             = var.password
	maintain_end_time    = "16:00Z"
	core_disk_type       = "cloud_efficiency"
	tags = {
		Created = "TF-update"
		For     = "acceptance test 123"
	}
}
```

此示例展示了如何通过各种配置(如网络设置、安全组和存储选项)创建一个 HBase 实例。

## 参数参考

支持以下参数：

* `name` - (必填) HBase 集群的名称。长度必须在 2-128 个字符之间，可以包含中文字符、英文字母、数字、点 (`.`)、下划线 (`_`) 或短横线 (`-`)。
* `zone_id` - (可选，强制新值) HBase 实例将启动的可用区 ID。如果指定了 `vswitch_id`，则此字段可以为空或与 VSwitch 的可用区一致。
* `engine` - (可选，强制新值) 集群的引擎类型。有效值为 `hbase`、`hbaseue` 或 `bds`。
* `engine_version` - (必填，强制新值) HBase 的主要版本。有效值：
  - 对于 `hbase`：`1.1` 或 `2.0`
  - 对于 `hbaseue`：`2.0`
  - 对于 `bds`：`1.0`
* `master_instance_type` - (必填，强制新值) 主节点的规格。请参阅 [实例规格](https://help.aliyun.com/document_detail/53532.html) 或使用 `describeInstanceType` API。
* `core_instance_type` - (必填，强制新值) 核心节点的规格。请参阅 [实例规格](https://help.aliyun.com/document_detail/53532.html) 或使用 `describeInstanceType` API。
* `core_instance_quantity` - (可选) 核心节点的数量。默认值为 `2`，范围是 `[1-200]`。
* `core_disk_type` - (可选，强制新值) 核心节点的磁盘类型。有效值：
  - `cloud_ssd`
  - `cloud_essd_pl1`
  - `cloud_efficiency`
  - `local_hdd_pro`
  - `local_ssd_pro`
  当 `engine=bds` 时，无需设置磁盘类型(或将其留为空字符串)。
* `core_disk_size` - (可选) 单个核心节点的存储大小(以 GB 为单位)。当 `engine=hbase/hbaseue` 时有效。对于 `bds` 不需要设置。取值范围：
  - 自定义存储空间：`[20, 64000]`
  - 集群：`[400, 64000]`，按 40GB 增加。
  - 单独：`[20-500]`，按 1GB 增加。
* `pay_type` - (可选) 付款类型。有效值为 `PrePaid` 或 `PostPaid`。默认为 `PostPaid`。
* `duration` - (可选，强制新值) 订阅时长(以月为单位)。有效值：`1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36`。仅当 `pay_type=PrePaid` 时有效。值 `12, 24, 36` 分别表示 1 年、2 年和 3 年。
* `auto_renew` - (可选，强制新值) 是否启用自动续费。有效值为 `true` 或 `false`。默认为 `false`。仅当 `pay_type=PrePaid` 时有效。
* `vswitch_id` - (可选，强制新值) VSwitch 的 ID。如果指定，则网络类型为 `vpc`。如果不指定，则网络类型为 `classic`。国际站点不支持经典网络。
* `cold_storage_size` - (可选，强制新值) 冷存储大小(以 GB 为单位)。有效值：`0` 或 `[800, 1000000]`，按 10GB 增加。`0` 表示禁用冷存储。
* `maintain_start_time` - (可选) 维护时间段的开始时间(UTC 格式 `HH:mmZ`)。例如：`02:00Z`。
* `maintain_end_time` - (可选) 维护时间段的结束时间(UTC 格式 `HH:mmZ`)。例如：`04:00Z`。
* `deletion_protection` - (可选) 是否启用删除保护。有效值为 `true` 或 `false`。默认为 `false`。
* `immediate_delete_flag` - (可选) 是否启用即时删除。有效值为 `true` 或 `false`。默认为 `false`。
* `tags` - (可选) 分配给资源的标签映射。
* `account` - (可选) 集群 Web UI 的帐户。长度必须在 `0-128` 个字符之间。
* `password` - (可选) 集群 Web UI 帐户的密码。长度必须在 `0-128` 个字符之间。
* `ip_white` - (可选) 集群的 IP 白名单。
* `security_groups` - (可选) 与集群关联的安全组 ID 列表。
* `master_instance_quantity` - (可选，计算) 集群中的主节点数量。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - HBase 集群的 ID。
* `zone_id` - 集群所在的可用区 ID。
* `master_instance_quantity` - 集群中的主节点数量。
* `ui_proxy_conn_addrs` - Web UI 代理连接信息列表。
  * `net_type` - 连接地址的访问类型。返回值：
    - `2`: 内网访问。
    - `0`: 公网访问。
  * `conn_addr` - 连接地址。
  * `conn_addr_port` - 连接端口。
* `zk_conn_addrs` - Zookeeper 连接信息列表。
  * `net_type` - 连接地址的访问类型。返回值：
    - `2`: 内网访问。
    - `0`: 公网访问。
  * `conn_addr` - 连接地址。
  * `conn_addr_port` - 连接端口。
* `slb_conn_addrs` - SLB 连接信息列表。
  * `net_type` - 连接地址的访问类型。返回值：
    - `2`: 内网访问。
    - `0`: 公网访问。
  * `conn_addr` - 连接地址。
  * `conn_addr_port` - 连接端口。
* `maintain_start_time` - (计算) 维护时间段的开始时间(UTC 格式 `HH:mmZ`)。
* `maintain_end_time` - (计算) 维护时间段的结束时间(UTC 格式 `HH:mmZ`)。
* `ip_white` - (计算) 集群的 IP 白名单。
* `security_groups` - (计算) 与集群关联的安全组 ID 列表。