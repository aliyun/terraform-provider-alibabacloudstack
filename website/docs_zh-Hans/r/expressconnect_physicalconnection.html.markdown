---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_physicalconnection"
sidebar_current: "docs-Alibabacloudstack-expressconnect-physicalconnection"
description: |- 
  编排高速物理通道
---

# alibabacloudstack_expressconnect_physicalconnection
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_express_connect_physical_connection`

使用Provider配置的凭证在指定的资源集下编排高速物理通道。

## 示例用法

### 基础用法

```terraform
variable "name" {
    default = "tf-testaccexpress_connectphysical_connection34886"
}

resource "alibabacloudstack_expressconnect_physicalconnection" "default" {
  device_name              = var.name
  access_point_id          = "ap-cn-hangzhou-jg-B"
  line_operator            = "CO"
  peer_location            = "XX街道"
  physical_connection_name = var.name
  type                     = "VPC"
  description              = "abcabc"
  port_type                = "1000Base-LX"
  bandwidth                = 10
}
```

## 参数说明

支持以下参数：

* `access_point_id` - (必填，变更时重建) 物理专线接入点的ID。
* `bandwidth` - (选填，可回读) 物理专线的带宽。单位：Gbps。如果不指定，则使用 API 的默认值。
* `circuit_code` - (选填) 运营商为物理专线提供的电路编码。
* `description` - (选填) 物理专线的描述信息。
* `line_operator` - (必填) 提供接入物理线路的运营商。有效值：
  * CT：中国电信
  * CU：中国联通
  * CM：中国移动
  * CO：中国其他
  * Equinix：Equinix
  * Other：境外其他
* `peer_location` - (必填) 本地数据中心的地理位置。
* `device_name` - (选填) 物理设备的名称。
* `physical_connection_name` - (选填) 物理专线的名称。
* `port_type` - (必填) 物理专线接入端口类型。有效值：
  * 100Base-T：百兆电口
  * 1000Base-T：千兆电口
  * 1000Base-LX：千兆单模光口(10千米)
  * 10GBase-T：万兆电口
  * 10GBase-LR：万兆单模光口(10千米)
  * 40GBase-LR：40千兆单模光口
  * 100GBase-LR：100千兆单模光口
  
  **注意**：从 v1.185.0+ 开始，`40GBase-LR` 和 `100GBase-LR` 值是有效的。根据背景端口的水位设置这些值。有关水位的详细信息，请联系业务经理。
* `redundant_physical_connection_id` - (选填) 冗余物理专线的ID。
* `status` - (选填) 物理专线的状态。用户可配置的有效值：
  * Confirmed：等待用户确认
  * Enabled：已开通
  * Canceled：已取消
  * Terminated：已终止
  
  -> **注意**：其他状态值（Initial、Approved、Allocating、Allocated、Rejected、Allocation Failed）为只读状态，由 API 返回，无法手动设置。
* `type` - (选填，变更时重建，可回读) 物理专线的类型。默认值为 `VPC`。

-> **注意**：`port_type` 和 `redundant_physical_connection_id` 参数在资源创建后无法修改。如果需要变更这些参数，资源将被重新创建。

-> **注意**：标注为"可回读"的参数由 API 返回并可读取，但可能无法由用户手动设置。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 物理专线的ID。
* `bandwidth` - 物理专线的带宽。单位：Gbps。
* `status` - 物理专线的状态。有效值：
  * Initial：申请中
  * Approved：审批通过
  * Allocating：正在分配资源
  * Allocated：接入施工中
  * Confirmed：等待用户确认
  * Enabled：已开通
  * Rejected：申请被拒绝
  * Canceled：已取消
  * Allocation Failed：资源分配失败
  * Terminated：已终止
* `type` - 物理专线的类型。默认值为 `VPC`。

## Import

高速物理通道资源可以通过 PhysicalConnectionId 进行导入，例如：

```
$ terraform import alibabacloudstack_expressconnect_physicalconnection.example pc-xxxxxxxxx
```