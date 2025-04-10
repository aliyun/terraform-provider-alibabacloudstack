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

## 参数说明

支持以下参数：

  * `ids` - (可选，强制更新) 配额ID列表。用于筛选特定的配额。
  * `product_name` - (必填) 服务名称。有效值包括：ECS（弹性计算服务）、OSS（对象存储服务）、VPC（虚拟私有云）、RDS（关系型数据库服务）、SLB（负载均衡）、ODPS（大数据处理服务）、GPDB（分析型数据库）、DDS（文档型数据库服务）、R-KVSTORE（缓存服务）和 EIP（弹性公网IP）。
  * `quota_type` - (必填) 配额类型。有效值为：`organization`（组织级别配额）和 `resourceGroup`（资源组级别配额）。
  * `quota_type_id` - (必填) 配额类型的ID。当 `quota_type` 参数设置为 `organization` 时，需要指定一个组织ID；当 `quota_type` 参数设置为 `resourceGroup` 时，需要指定一个资源组ID。
  * `target_type` - (可选) 这个保留参数是可选的，可以留空。它仅用于某些产品。需要 `target_type` 的产品及其值包括：RDS (`MySql`)、R-KVSTORE (`redis`) 和 DDS (`mongodb`)。

## 属性说明

除了上述参数外，还导出以下属性：

* `quotas` - 配额列表。每个元素包含以下属性：
  * `id` - 配额的唯一标识符。
  * `quota_type` - 配额所属的类型名称，例如组织或资源组。
  * `quota_type_id` - 配额所属的类型ID，例如组织ID或资源组ID。
  * `total_vip_internal` - 内部VIP的总数量。
  * `target_type` - 配额目标类型。对于某些产品，此字段是必需的。需要 `target_type` 的产品及其值包括：RDS (`MySql`)、R-KVSTORE (`redis`) 和 DDS (`mongodb`)。
  * `region` - 配额所属区域的名称。
  * `total_vip_public` - 公共VIP的总数量。
  * `total_vpc` - 总VPC数量。
  * `total_cpu` - 总CPU核心数。
  * `total_cu` - 总计算单元（CU）数量。
  * `total_disk` - 总磁盘容量（单位通常为GB）。
  * `total_mem` - 总内存容量（单位通常为GB）。
  * `used_mem` - 已使用的内存容量。
  * `total_gpu` - 总GPU数量。
  * `total_amount` - 总金额或预算限制。
  * `total_disk_cloud_ssd` - 总SSD云盘容量。
  * `used_disk` - 已使用的磁盘容量。
  * `allocate_disk` - 已分配的磁盘容量。
  * `allocate_cpu` - 已分配的CPU核心数。
  * `total_eip` - 总EIP（弹性公网IP）数量。
  * `total_disk_cloud_efficiency` - 总高效云盘容量。
  * `allocate_vip_public` - 已分配的公共VIP地址数量。
  * `allocate_vip_internal` - 已分配的内部VIP地址数量。
  * `used_vip_public` - 已使用的公共VIP地址数量。

**注意**：只修改了中文文档中的`## 参数说明`和`## 属性说明`部分，其他内容保持不变。