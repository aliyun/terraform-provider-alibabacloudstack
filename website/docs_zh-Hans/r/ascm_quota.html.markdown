---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_quota"
sidebar_current: "docs-alibabacloudstack-resource-ascm-quota"
description: |-
  编排Ascm配额
---

# alibabacloudstack_ascm_quota

使用Provider配置的凭证在指定的资源集下编排Ascm配额。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_quota" "default" {
  quota_type = "organization"
  quota_type_id = alibabacloudstack_ascm_organization.default.parent_id // 用于创建父级配额
  //quota_type_id = alibabacloudstack_ascm_organization.default.org_id
  product_name = "ECS"
  total_cpu = 1000
  total_mem = 1000
  total_gpu = 1000
  total_disk_cloud_ssd = 10000
  total_disk_cloud_efficiency = 20000
}

output "quota" {
  value = alibabacloudstack_ascm_quota.default.*
}
```

## 参数参考

以下参数被支持：
### 在为组织的任何服务创建配额之前，首先通过使用该组织的parent_id创建父级配额。

* `product_name` - (必填) 服务名称。有效值：ECS、OSS、VPC、RDS、SLB、ODPS、GPDB、DDS、R-KVSTORE 和 EIP。
* `region_name` - (可选) 应用配额的区域名称。
* `quota_type` - (必填) 配额类型。有效值：organization 和 resourceGroup。
* `quota_type_id` - (必填) 配额类型的ID。当QuotaType参数设置为organization时，指定组织ID；当QuotaType参数设置为resourceGroup时，指定资源组ID。
* `cluster_name` - (可选) 集群名称。此保留参数是可选的，可以留空。
* `total_cpu` - (可选) 此保留参数是可选的，可以留空。
* `total_mem` - (可选) 此保留参数是可选的，可以留空。
* `total_gpu` - (可选) 此保留参数是可选的，可以留空。
* `total_disk_cloud_ssd` - (可选) 此保留参数是可选的，可以留空。
* `total_disk_cloud_efficiency` - (可选) 此保留参数是可选的，可以留空。
* `total_vip_internal` - (可选) 此保留参数是可选的，可以留空。
* `total_vip_public` - (可选) 此保留参数是可选的，可以留空。
* `total_vpc` - (可选) 此保留参数是可选的，可以留空。
* `total_amount` - (可选) 此保留参数是可选的，可以留空。
* `total_eip` - (可选) 此保留参数是可选的，可以留空。
* `total_disk` - (可选) 此保留参数是可选的，可以留空。
* `total_cu` - (可选) 此保留参数是可选的，可以留空。
* `target_type` - (可选) 此保留参数是可选的，可以留空。它仅在某些产品中使用。需要target_type的产品及其值为：RDS ("MySql")，R-KVSTORE ("redis") 和 DDS ("mongodb")。
* `region_name` - (可选) 应用配额的区域名称。
* `cluster_name` - (可选) 配额关联的集群名称。

您可以调用此操作来创建配额。根据产品名称使用参数。
以下是产品的示例：

ECS

* `total_cpu`:100, `total_mem`:100, `total_gpu`:100, `total_disk_cloud_ssd`:100, `total_disk_cloud_efficiency`:100

OSS

* `total_amount`:100


VPC

* `total_vpc`:100

RDS

* `total_cpu`:100, `total_mem`:100, `total_disk`:100, `target_type`: "MySql"

SLB

* `total_vip_internal`:100, `total_vip_public`:100

MaxCompute (ODPS)

* `total_cu`:100, `total_disk`:100

EIP

* `total_eip`:100

AnalyticDB for PostgreSQL (GPDB)

* `total_cpu`: 100, `total_mem`: 100, `total_disk`:100

KVStore for Redis (R-KVSTORE)

* `total_mem`: 100, `target_type`: "redis"

ApsaraDB for MongoDB (DDS)

* `total_cpu`: 100, `total_mem`: 100, `total_disk`:100, `target_type`: "mongodb"

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `quota_id` - 配额ID。
* `id` - 服务的ProductName、QuotaType和QuotaTypeId。格式为 `ProductName:QuotaType:QuotaTypeId`。