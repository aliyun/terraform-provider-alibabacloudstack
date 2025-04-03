---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackageattachment"
sidebar_current: "docs-Alibabacloudstack-cbwp-commonbandwidthpackageattachment"
description: |- 
  编排绑定共享带宽包和弹性公网IP(EIP)
---

# alibabacloudstack_cbwp_commonbandwidthpackageattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_common_bandwidth_package_attachment`

使用Provider配置的凭证在指定的资源集下编排绑定共享带宽包和弹性公网IP(EIP)

## 示例用法

### 基础用法

以下示例展示了如何创建一个共享带宽包并将其与一个弹性公网IP(EIP)关联。

```hcl
variable "name" {
    default = "tf-testAccBandwidtchPackage3659"
}

# 创建共享带宽包
resource "alibabacloudstack_common_bandwidth_package" "default" {
    bandwidth   = "2"
    name        = "${var.name}"
    description = "${var.name}_description"
}

# 创建弹性公网IP
resource "alibabacloudstack_eip" "default" {
    name     = "${var.name}"
    bandwidth = "2"
}

# 将EIP与共享带宽包关联
resource "alibabacloudstack_common_bandwidth_package_attachment" "default" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.default.id}"
    instance_id          = "${alibabacloudstack_eip.default.id}"
}
```

### 高级用法(多个EIP)

以下示例展示了如何将多个EIP与同一个共享带宽包关联。

```hcl
variable "name" {
    default = "tf-testAccBandwidtchPackageAdvanced"
}

# 创建共享带宽包
resource "alibabacloudstack_common_bandwidth_package" "advanced" {
    bandwidth   = "10"
    name        = "${var.name}"
    description = "${var.name}_description"
}

# 创建第一个EIP
resource "alibabacloudstack_eip" "eip1" {
    name     = "${var.name}-eip1"
    bandwidth = "2"
}

# 创建第二个EIP
resource "alibabacloudstack_eip" "eip2" {
    name     = "${var.name}-eip2"
    bandwidth = "3"
}

# 将第一个EIP与共享带宽包关联
resource "alibabacloudstack_common_bandwidth_package_attachment" "eip1_attachment" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.advanced.id}"
    instance_id          = "${alibabacloudstack_eip.eip1.id}"
}

# 将第二个EIP与共享带宽包关联
resource "alibabacloudstack_common_bandwidth_package_attachment" "eip2_attachment" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.advanced.id}"
    instance_id          = "${alibabacloudstack_eip.eip2.id}"
}
```

## 参数参考

支持以下参数：

* `bandwidth_package_id` - (必填，变更时重建) 共享带宽实例的ID。此字段在创建后无法修改。
* `instance_id` - (必填，变更时重建) 要与共享带宽包关联的弹性公网IP(EIP)的ID。此字段在创建后无法修改。您可以指定多达50个EIP ID进行关联。如果需要，可以用逗号(,)分隔多个ID。如果同时传入**EipAddress**和**AllocationId**参数，**AllocationId**可输入50个EIP的实例ID，**EipAddress**也可同时输入50个EIP的IP地址。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 附件的唯一标识符。格式为 `<bandwidth_package_id>:<instance_id>`。
* `status` - 附件的状态。可能的值包括 `Attached` 和 `Detached`。
* `creation_time` - 附件创建的时间。这对于跟踪附件的生命周期非常有用。
```