---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_quotas"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-quotas"
description: |-
    查询配额

---

# alibabacloudstack_ascm_quotas

根据指定过滤条件列出当前凭证权限可以访问的配额列表。

## 示例用法

```
resource "alibabacloudstack_ascm_organization" "default" {
    name = "Dummy_Test_1"
}

resource "alibabacloudstack_ascm_quota" "default" {
    quota_type = "organization"
    quota_type_id = alibabacloudstack_ascm_organization.default.parent_id
    product_name = "RDS"
    total_cpu = 1500
    total_disk = 1500
    total_mem = 1500
    target_type = "MySql"
}

data "alibabacloudstack_ascm_quotas" "default" {
    quota_type = "organization"
    quota_type_id = alibabacloudstack_ascm_organization.default.parent_id
    product_name = "RDS"
    target_type = "MySql"
    output_file = "Rds_quota"
}
output "quota" {
    value = data.alibabacloudstack_ascm_quotas.default.*
}
```

## 参数参考

支持以下参数：

  * `ids` - (可选，强制更新) 配额ID列表。
  * `product_name` - (必填)  服务名称。有效值：ECS、OSS、VPC、RDS、SLB、ODPS、GPDB、DDS、R-KVSTORE 和 EIP。
  * `quota_type` - (必填)  配额类型。有效值：organization 和 resourceGroup。
  * `quota_type_id` - (必填)  配额类型的ID。当 QuotaType 参数设置为 organization 时，指定一个组织ID；当 QuotaType 参数设置为 resourceGroup 时，指定一个资源组ID。
  * `target_type` - (可选) 这个保留参数是可选的，可以留空。它仅用于某些产品。需要 target_type 的产品及其值包括：RDS ("MySql")、R-KVSTORE ("redis") 和 DDS ("mongodb")。
  
## 属性参考

除了上述参数外，还导出以下属性：

* `quotas` - 配额列表。每个元素包含以下属性：
  * `id` - 配额ID。
  * `quota_type` - 组织或资源组的名称。
  * `quota_type_id` - 组织或资源组的ID。
  * `total_vip_internal` - 总内部VIP。
  * `target_type` - 它仅用于某些产品。需要 target_type 的产品及其值包括：RDS ("MySql")、R-KVSTORE ("redis") 和 DDS ("mongodb")。
  * `region` - 产品所属区域的名称。
  * `total_vip_public` - 总公共VIP。
  * `total_vpc` - 总VPC。
  * `total_cpu` - 总CPU。
  * `total_cu` - 总CU。
  * `total_disk` - 总磁盘。
  * `total_mem` - 总内存。
  * `used_mem` - 已使用内存。
  * `total_gpu` - 总GPU。
  * `total_amount` - 总金额。
  * `total_disk_cloud_ssd` - 总SSD云盘。
  * `used_disk` - 已使用磁盘。
  * `allocate_disk` - 已分配磁盘。
  * `allocate_cpu` - 已分配CPU。
  * `total_eip` - 总EIP。
  * `total_disk_cloud_efficiency` - 总高效云盘。
  * `allocate_vip_public` - 已分配的公共VIP地址数量。
  * `allocate_vip_internal` - 已分配的内部VIP地址数量。
  * `used_vip_public` - 已使用的公共VIP地址数量。