---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackages"
sidebar_current: "docs-Alibabacloudstack-datasource-cbwp-commonbandwidthpackages"
description: |- 
  查询共享带宽包

---

# alibabacloudstack_cbwp_commonbandwidthpackages
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_common_bandwidth_packages`

根据指定过滤条件列出当前凭证权限可以访问的共享带宽包列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccCommonBandwidthPackageDataSource5327"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
  bandwidth   = "2"
  name        = "${var.name}"
  description = "${var.name}_description"
}

data "alibabacloudstack_cbwp_commonbandwidthpackages" "foo" {
  name_regex        = "${alibabacloudstack_common_bandwidth_package.default.name}"
  ids               = [alibabacloudstack_common_bandwidth_package.default.id]
  resource_group_id = "your_resource_group_id"
}

output "common_bandwidth_packages" {
  value = data.alibabacloudstack_cbwp_commonbandwidthpackages.foo.packages
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) 共享带宽包 ID 列表。
* `name_regex` - (可选，变更时重建) 用于按名称筛选结果的正则表达式字符串。
* `resource_group_id` - (可选，变更时重建) 您要将资源移动到的资源组的 ID。您可以使用资源组来方便地对资源进行分组和权限管理。有关更多信息，请参见 [什么是资源管理？](https://help.aliyun.com/document_detail/94475.html)

## 属性说明

除了上述所有参数外，还导出以下属性：

* `names` - 共享带宽包名称列表。
* `ids` - 共享带宽包 ID 列表。
* `packages` - 共享带宽包列表。每个元素包含以下属性：
  * `id` - 共享带宽包的 ID。
  * `bandwidth` - 共享带宽的带宽峰值。单位：Mbps。默认取值范围：**1**~**1000**。默认值：**1**。
    - 当 **InternetChargeType** 取值为 **PayByBandwidth**，即共享带宽的计费方式为按带宽计费时，**Bandwidth** 的默认取值范围为 **2**~**20000**。
    - 当 **InternetChargeType** 取值为 **95PayBy95**，即共享带宽的计费方式为按增强型95计费时，**Bandwidth** 的默认取值范围为 **200**~**20000**。
    - 当 **InternetChargeType** 取值为 **PayByDominantTraffic**，即共享带宽的计费方式为按主流量计费时，**Bandwidth** 的默认取值范围为 **1**~**2000**。 默认值：**1000**。
  * `status` - 共享带宽实例的状态。默认取值：**Available**。
  * `name` - 共享带宽包的名称。
  * `description` - 共享带宽包实例的描述。描述必须是 2 到 256 个字符长度，并以字母开头。描述不能以 `http://` 或 `https://` 开头。
  * `business_status` - 共享带宽包实例的业务状态。取值：
    - **Normal**: 正常。
    - **FinancialLocked**: 欠费。
    - **Unactivated**: 未激活。
  * `isp` - 线路类型。有效值：
    - **BGP**: 所有区域均支持 BGP（多线）。
    - **BGP_PRO**: BGP（多线）Pro 线路在以下区域可用：中国（香港）、新加坡、日本（东京）、菲律宾（马尼拉）、马来西亚（吉隆坡）、印度尼西亚（雅加达）和泰国（曼谷）。
    - 如果您允许使用单线带宽，还可以使用以下值之一：
      - **ChinaTelecom**
      - **ChinaUnicom**
      - **ChinaMobile**
      - **ChinaTelecom_L2**
      - **ChinaUnicom_L2**
      - **ChinaMobile_L2**
    - 如果您的服务部署在中国东部 1 财经区域，则此参数是必需的，且必须设置为 **BGP_FinanceCloud**。
  * `creation_time` - 共享带宽包创建的时间。
  * `public_ip_addresses` - 属于共享带宽包的公共 IP 地址。每个元素包含以下属性：
    * `ip_address` - 弹性公网 IP 的地址。
    * `allocation_id` - 弹性公网 IP 实例的 ID。