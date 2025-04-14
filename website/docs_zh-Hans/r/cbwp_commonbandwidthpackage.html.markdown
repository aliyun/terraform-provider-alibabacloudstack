---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackage"
sidebar_current: "docs-Alibabacloudstack-cbwp-commonbandwidthpackage"
description: |- 
  编排共享带宽包列表
---

# alibabacloudstack_cbwp_commonbandwidthpackage
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_common_bandwidth_package`

使用Provider配置的凭证在指定的资源集下编排共享带宽包。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf-testAccCommonBandwidthPackage481904"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
  internet_charge_type = "PayByTraffic"
  bandwidth            = "10"
  name                = var.name
  description         = "Test common bandwidth package"
}
```

高级用法(包含所有可选参数)

```hcl
variable "name" {
  default = "tf-testAccCommonBandwidthPackageAdvanced"
}

resource "alibabacloudstack_common_bandwidth_package" "advanced" {
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "100"
  ratio                = "20"
  name                = var.name
  description         = "Advanced test common bandwidth package"
}
```

## 参数说明

支持以下参数：

* `bandwidth` - (必填) 共享带宽实例的最大带宽。单位：Mbit/s。有效值范围：
  - **国际站**：默认取值范围为 **1**~**1000**，默认值为 **1**。
  - **中国站**：
    - 当 `internet_charge_type` 为 `PayByBandwidth` 时，默认取值范围为 **2**~**20000**。
    - 当 `internet_charge_type` 为 `95PayBy95` 时，默认取值范围为 **200**~**20000**。
    - 当 `internet_charge_type` 为 `PayByDominantTraffic` 时，默认取值范围为 **1**~**2000**。
    - 默认值为 **1000**。

* `internet_charge_type` - (可选，变更时重建) 共享带宽实例的计费方式。取值：
  - **国际站**：`PayByTraffic`(按流量计费)。
  - **中国站**：
    - `PayByBandwidth`(默认值)：按带宽计费。
    - `PayBy95`：按增强型95计费。
    - `PayByDominantTraffic`：按主流量计费。

* `ratio` - (可选，变更时重建) 共享带宽的保底百分比，仅取值为 **20**。当 `internet_charge_type` 取值为 `PayBy95` 时需配置此参数。> 仅中国站支持该参数。

* `name` - (可选) 共享带宽实例的名称。长度为 2~128 个字符，可以包含字母、数字、下划线 (`_`) 和连字符 (`-`)。名称必须以字母开头。

* `description` - (可选) 共享带宽实例的描述。长度为 2~256 个字符，并以字母开头。描述不能以 `http://` 或 `https://` 开头。

* `bandwidth_package_name` - (可选) 共享带宽实例的名称。长度为 2~128 个字符，可以包含字母、数字、下划线 (`_`) 和连字符 (`-`)。名称必须以字母开头。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 共享带宽包实例的 ID。
* `bandwidth_package_name` - 共享带宽实例的名称。长度为 2~128 个字符，必须以字母或中文开头，可包含数字、半角句点 (`.`)、下划线 (`_`) 和短划线 (`-`)。但不能以 `http://` 或 `https://` 开头。
* `name` - 共享带宽实例的名称。长度为 2~128 个字符，可以包含字母、数字、下划线 (`_`) 和连字符 (`-`)。名称必须以字母开头。