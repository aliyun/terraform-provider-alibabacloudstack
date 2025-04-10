---
subcategory: "HBase"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hbase_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-hbase-clusters"
description: |- 
  查询云数据库hbase集群。
---

# alibabacloudstack_hbase_clusters
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_hbase_instances`

根据指定过滤条件列出当前凭证权限可以访问的云数据库hbase集群列表。

## 示例用法

以下示例展示了如何使用 `alibabacloudstack_hbase_clusters` 数据源来查询 HBase 集群列表，并结合其他资源完成更复杂的配置。

```terraform
variable "name" {
  default = "tf-testAccHBaseInstance_datasource_17791"
}

# 查询可用区
data "alibabacloudstack_hbase_zones" "default" {}

# 查询 VPC
data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

# 查询或创建 VSwitch
data "alibabacloudstack_vswitches" "default" {
  vpc_id  = data.alibabacloudstack_vpcs.default.ids[0]
  zone_id = data.alibabacloudstack_hbase_zones.default.ids[0]
}

resource "alibabacloudstack_vswitch" "vswitch" {
  count           = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id          = data.alibabacloudstack_vpcs.default.vpcs[0].id
  cidr_block      = cidrsubnet(data.alibabacloudstack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id        = data.alibabacloudstack_hbase_zones.default.ids[0]
  vswitch_name   = var.name
}

locals {
  vswitch_id = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? data.alibabacloudstack_vswitches.default.ids[0] : concat(alibabacloudstack_vswitch.vswitch.*.id, [""])[0]
}

# 创建 HBase 实例
resource "alibabacloudstack_hbase_instance" "default" {
  name                  = var.name
  engine_version        = "2.0"
  master_instance_type  = "hbase.sn1.large"
  core_instance_type    = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type        = "cloud_efficiency"
  pay_type              = "PostPaid"
  duration              = 1
  auto_renew            = false
  vswitch_id            = local.vswitch_id
  cold_storage_size     = 0
  deletion_protection   = false
  immediate_delete_flag = true
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

# 查询 HBase 集群列表
data "alibabacloudstack_hbase_clusters" "hbase" {
  name_regex        = "${alibabacloudstack_hbase_instance.default.name}"
  availability_zone = data.alibabacloudstack_hbase_zones.default.ids[0]
  ids               = [alibabacloudstack_hbase_instance.default.id]
  output_file      = "hbase_clusters_output.txt"
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (选填) 应用于集群名称的正则表达式字符串。这允许基于名称使用正则表达式过滤集群。
* `ids` - (选填) HBase 集群的 ID 列表。可以使用此参数通过唯一标识符过滤集群。
* `availability_zone` - (选填) HBase 集群所在的可用区。使用此参数可以过滤特定区域内的集群。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - HBase 集群的 ID 列表。
* `names` - HBase 集群的名称列表。
* `instances` - HBase 集群的列表。每个元素包含以下属性：
  * `id` - HBase 集群的 ID。
  * `name` - HBase 集群的名称。
  * `region_id` - 集群所在区域的 ID。
  * `zone_id` - 集群所在可用区的 ID。
  * `engine` - 集群的引擎类型（例如 HBase）。
  * `engine_version` - 集群使用的引擎版本。
  * `network_type` - 集群的网络类型，如 `Classic` 或 `VPC`。
  * `master_instance_type` - 主节点的实例类型（例如 `hbase.sn2.2xlarge`）。
  * `master_node_count` - 集群中的主节点数量。
  * `core_instance_type` - 核心节点的实例类型（例如 `hbase.sn2.4xlarge`）。
  * `core_node_count` - 集群中的核心节点数量。
  * `core_disk_type` - 核心节点的磁盘类型，如 `Cloud_SSD` 或 `Cloud_Efficiency`。
  * `core_disk_size` - 核心节点的磁盘大小（以 GB 为单位）。
  * `vpc_id` - 与集群关联的 VPC ID。
  * `vswitch_id` - 与集群关联的 VSwitch ID。
  * `pay_type` - 集群的计费方式。可能的值包括 `PostPaid`（按量付费）和 `PrePaid`（包年包月订阅）。
  * `created_time` - 集群的创建时间。
  * `expire_time` - 集群的过期时间（如果适用）。
  * `status` - 集群的当前状态。
  * `backup_status` - 集群的备份状态。
  * `deletion_protection` - 指示是否为集群启用了删除保护。
  * `tags` - 分配给集群的标签映射。