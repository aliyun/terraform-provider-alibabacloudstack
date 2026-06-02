---
subcategory: "共享带宽包"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_common_bandwidth_packages"
sidebar_current: "docs-Alibabacloudstack-datasource-common-bandwidth-packages"
description: |-
  提供阿里云账号下拥有的共享带宽包列表。
---

# alibabacloudstack\_common\_bandwidth\_packages

此数据源提供根据指定过滤条件列出的阿里云账号下的共享带宽包资源列表。

-> **注意：** 版本 v1.59.0+ 可用。

## 示例用法

```hcl
data "alibabacloudstack_common_bandwidth_packages" "example" {
  name_regex = "^my-CBWP"
}

output "first_cbwp_id" {
  value = data.alibabacloudstack_common_bandwidth_packages.example.packages.0.id
}
```

## 参数说明

支持以下参数：

* `ids` - (可选, 变更后重建) 共享带宽包ID列表。
* `name_regex` - (可选, 变更后重建) 用于按共享带宽包名称过滤结果的正则表达式字符串。
* `resource_group_id` - (可选, 变更后重建) 资源组ID。
* `output_file` - (可选, 已废弃) 输出文件路径。此字段已废弃，计划在 3.19.0 版本中移除。要将内容写入文件，请使用 `local_file` provider。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 共享带宽包ID列表。
* `names` - 共享带宽包名称列表。
* `packages` - 共享带宽包列表。每个元素包含以下属性：
  * `id` - 共享带宽包ID。
  * `bandwidth` - 共享带宽包的带宽峰值。单位：Mbps。
  * `status` - 共享带宽包的状态。取值：`Available`（可用）、`Modifying`（修改中）。
  * `name` - 共享带宽包的名称。
  * `description` - 共享带宽包的描述信息。
  * `business_status` - 共享带宽包的业务状态。取值：`Normal`（正常）、`FinancialLocked`（财务锁定）。
  * `isp` - 共享带宽包的线路类型。取值：`BGP`、`BGP_PRO`、`ChinaTelecom`、`ChinaUnicom`、`ChinaMobile`。
  * `creation_time` - 共享带宽包的创建时间。
  * `public_ip_addresses` - 与共享带宽包关联的公网IP地址列表。每个元素包含：
    * `ip_address` - 公网IP地址。
    * `allocation_id` - EIP的分配ID。
