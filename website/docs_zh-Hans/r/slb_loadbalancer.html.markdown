---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_loadbalancer"
sidebar_current: "docs-Alibabacloudstack-slb-loadbalancer"
description: |- 
  编排负载均衡(SLB)实例
---

# alibabacloudstack_slb_loadbalancer
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)实例。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccslbload_balancer19164"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_slb_loadbalancer" "default" {
  address_type = "intranet"
  name         = "rdk_test_name"
  specification = "slb.s1.small"
  vswitch_id   = "${alibabacloudstack_vpc_vswitch.default.id}"

  tags = {
    Environment = "Test"
    CreatedBy   = "Terraform"
  }
}
```

## 参数参考

支持以下参数：

* `name` - (可选) SLB的名称。该名称必须在您的AlibabaCloudStack账户中是唯一的，最多可以包含80个字符，只能包含字母、数字或连字符(-)，并且不能以连字符开头或结尾。如果不指定，Terraform将自动生成一个以`tf-lb`开头的名称。
* `address_type` - (可选，强制新资源)SLB实例的网络类型。有效值：["internet", "intranet"]。如果负载均衡器是在VPC中启动的，则此值必须为"intranet"。
  * `internet`：创建Internet SLB实例后，系统会分配一个公共IP地址，以便实例可以从Internet转发请求。
  * `intranet`：创建内网SLB实例后，系统会分配一个内网IP地址，以便实例只能转发内网请求。
* `specification` - (可选) 服务器负载均衡实例的规格。默认为空字符串，表示它是“共享性能”实例。有效值包括：
  * `slb.s1.small`
  * `slb.s2.small`
  * `slb.s2.medium`
  * `slb.s3.small`
  * `slb.s3.medium`
  * `slb.s3.large`
  * `slb.s4.large`
* `vswitch_id` - (必填，对于VPC SLB，强制新资源)启动SLB所在的交换机ID。如果`address_type`设置为"internet"，此字段将被忽略。
* `tags` - (可选) 要分配给资源的标签映射。
* `address` - (可选，强制新资源)负载均衡实例的服务地址。此字段由系统根据`address_type`自动分配。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 负载均衡器的ID。
* `address` - 负载均衡器的IP地址。
* `address_type` - 负载均衡实例的地址类型。
* `specification` - 服务器负载均衡实例的规格。
* `vswitch_id` - 与负载均衡器关联的交换机ID。