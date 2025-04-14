---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_physicalconnections"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-physicalconnections"
description: |- 
  查询高速通道物理连接
---

# alibabacloudstack_expressconnect_physicalconnections
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_express_connect_physical_connections`

根据指定过滤条件列出当前凭证权限可以访问的查询高速物理连接列表。

## 示例用法

### 基础用法：

```terraform
variable "name" {
  default = "tf-testAlibabacloudstackExpressconnectPhysicalConnections45872"
}

data "alibabacloudstack_expressconnect_physicalconnections" "ids" {
  ids = ["pc-2345678"]
}

output "express_connect_physical_connection_id_1" {
  value = data.alibabacloudstack_expressconnect_physicalconnections.ids.connections.0.id
}

data "alibabacloudstack_expressconnect_physicalconnections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

output "express_connect_physical_connection_id_2" {
  value = data.alibabacloudstack_expressconnect_physicalconnections.nameRegex.connections.0.id
}
```

## 参数说明

以下参数是支持的：

* `include_reservation_data` - (可选，变更时重建) 指定是否包括尚未生效的保留数据。有效值：`true` 或 `false`。
* `ids` - (可选，变更时重建) 物理专线 ID 列表。
* `name_regex` - (可选，变更时重建) 用于通过物理专线名称筛选结果的正则表达式字符串。
* `status` - (可选，变更时重建) 物理专线的状态。有效值：
  * `Initial`: 在申请中。
  * `Approved`: 已批准。
  * `Allocating`: 正在分配资源。
  * `Allocated`: 接入建设中。
  * `Confirmed`: 等待用户确认。
  * `Enabled`: 已激活。
  * `Rejected`: 申请被拒绝。
  * `Canceled`: 已取消。
  * `Allocation Failed`: 资源分配失败。
  * `Terminated`: 已终止。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 物理专线名称列表。
* `connections` - Express Connect 物理专线列表。每个元素包含以下属性：
  * `access_point_id` - 物理专线接入点的 ID。
  * `ad_location` - 物理专线接入设备所在的物理位置。
  * `bandwidth` - 物理专线的带宽。单位：Gbps。
  * `business_status` - 物理专线的商业状态。有效值：
    * `Normal`: 已经激活。
    * `FinancialLocked`: 欠费锁定。
    * `SecurityLocked`: 安全原因锁定。
  * `circuit_code` - 运营商为物理专线提供的电路编码。
  * `create_time` - 物理专线的创建时间。时间按照 ISO8601 标准表示，并使用 UTC 时间。格式为：YYYY-MM-DDThh:mm:ssZ。
  * `description` - 物理专线的描述信息。
  * `enabled_time` - 物理专线的开通时间。
  * `end_time` - 物理专线的过期时间。
  * `has_reservation_data` - 是否包含未生效的订单数据。有效值：`true` 或 `false`。
  * `id` - 物理专线的 ID。
  * `line_operator` - 提供物理线路接入的运营商。有效值：
    * `CT`: 中国电信。
    * `CU`: 中国联通。
    * `CM`: 中国移动。
    * `CO`: 中国其他。
    * `Equinix`: Equinix。
    * `Other`: 其他国外。
  * `loa_status` - LOA 的状态。有效值：
    * `Applying`: LOA 申请中。
    * `Accept`: LOA 申请通过。
    * `Available`: LOA 可用。
    * `Rejected`: LOA 申请被拒绝。
    * `Completing`: 专线正在施工。
    * `Complete`: 专线施工完成。
    * `Deleted`: LOA 已删除。
  * `payment_type` - 资源的付费类型。
  * `peer_location` - 本地数据中心的地理位置。
  * `physical_connection_id` - 专线实例的 ID。
  * `physical_connection_name` - 物理专线的名称。
  * `port_number` - 物理专线设备的端口号。
  * `port_type` - 物理专线端口类型。有效值：
    * `100Base-T`: 100 兆电口。
    * `1000Base-T`: 千兆电口。
    * `1000Base-LX`: 千兆单模光口(10km)。
    * `10GBase-T`: 万兆电口。
    * `10GBase-LR`: 万兆单模光口(10km)。
  * `redundant_physical_connection_id` - 冗余物理专线的 ID。
  * `reservation_active_time` - 未生效订单的生效时间。
  * `reservation_internet_charge_type` - 未生效订单的付费类型。值：`PayByBandwidth`，表示按带宽付费。
  * `spec` - 物理专线的规格。单位：G 表示 Gbps。
  * `status` - 物理专线的状态。参见 `status` 参数的有效值。
  * `type` - 物理专线的类型。默认值为 `VPC`。
  * `reservation_order_type` - 预留订单的订单类型。