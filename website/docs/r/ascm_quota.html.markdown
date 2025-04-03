---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_quota"
sidebar_current: "docs-alibabacloudstack-resource-ascm-quota"
description: |-
  Provides a Ascm quota resource.
---

# alibabacloudstack_ascm_quota

Provides a Ascm quota resource.

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_quota" "default" {
  quota_type = "organization"
  quota_type_id = alibabacloudstack_ascm_organization.default.parent_id // For creating FatherQuota
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
## Argument Reference

The following arguments are supported:
### Before creating quota for any service of an organization, create Father Quota first by using parent_id of that organization.

* `product_name` - (Required) The name of the service. Valid values: ECS, OSS, VPC, RDS, SLB, ODPS, GPDB, DDS, R-KVSTORE, and EIP.
* `region_name`- (Optional) The name of region to apply quota.
* `quota_type` - (Required) The type of the quota. Valid values: organization and resourceGroup.
* `quota_type_id` - (Required) The ID of the quota type. Specify an organization ID when the QuotaType parameter is set to organization. Specify a resource set ID when the QuotaType parameter is set to resourceGroup.
* `cluster_name` - (Optional) The name of the cluster. This reserved parameter is optional and can be left empty.
* `total_cpu` - (Optional) This reserved parameter is optional and can be left empty.
* `total_mem` - (Optional) This reserved parameter is optional and can be left empty.
* `total_gpu` - (Optional) This reserved parameter is optional and can be left empty.
* `total_disk_cloud_ssd` - (Optional) This reserved parameter is optional and can be left empty.
* `total_disk_cloud_efficiency` - (Optional) This reserved parameter is optional and can be left empty.
* `total_vip_internal` - (Optional) This reserved parameter is optional and can be left empty.
* `total_vip_public` - (Optional) This reserved parameter is optional and can be left empty.
* `total_vpc` - (Optional) This reserved parameter is optional and can be left empty.
* `total_amount` - (Optional) This reserved parameter is optional and can be left empty.
* `total_eip` - (Optional) This reserved parameter is optional and can be left empty.
* `total_disk` - (Optional) This reserved parameter is optional and can be left empty.
* `total_cu` - (Optional) This reserved parameter is optional and can be left empty.
* `target_type` - (Optional) This reserved parameter is optional and can be left empty. It will be used only for some products. Products where target_type are required with their values - RDS ("MySql"), R-KVSTORE ("redis") and DDS ("mongodb").
* `region_name` - (Optional)  The name of the region to apply the quota.
* `cluster_name` - (Optional)  The name of the cluster associated with the quota.

You can call this operation to create a quota. Use parameters according to the product name.
 Sample for the product.

ECS

* `total_cpu`:100,`total_mem`:100,`total_gpu`:100,`total_disk_cloud_ssd`:100,`total_disk_cloud_efficiency`:100

OSS

* `total_amount`:100


VPC

* `total_vpc`:100

RDS

* `total_cpu`:100,`total_mem`:100,`total_disk`:100, `target_type`: "MySql"

SLB

* `total_vip_internal`:100,`total_vip_public`:100

MaxCompute (ODPS)

* `total_cu`:100,`total_disk`:100

EIP

* `total_eip`:100

AnalyticDB for PostgreSQL (GPDB)

* `total_cpu`: 100, `total_mem`: 100, `total_disk`:100

KVStore for Redis (R-KVSTORE)

* `total_mem`: 100, `target_type`: "redis"

ApsaraDB for MongoDB (DDS)

* `total_cpu`: 100, `total_mem`: 100, `total_disk`:100, `target_type`: "mongodb"

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `quota_id` - ID of the quota. 
* `id` - ProductName, QuotaType and QuotaTypeId of the Service. The value is in format `ProductName:QuotaType:QuotaTypeId`.