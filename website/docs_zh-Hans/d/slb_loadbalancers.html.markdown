---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_loadbalancers"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-loadbalancers"
description: |- 
  查询负载均衡(SLB)实例
---

# alibabacloudstack_slb_loadbalancers
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slbs`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)实例列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-SlbDataSourceSlbs"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_slb" "default" {
  name        = "${var.name}_slb"
  vswitch_id  = "${alibabacloudstack_vpc_vswitch.default.id}"
}

data "alibabacloudstack_slb_loadbalancers" "default" {
  ids = ["${alibabacloudstack_slb.default.id}"]

  # 使用正则表达式筛选负载均衡器名称
  name_regex = "tf-SlbDataSourceSlbs.*"

  # 筛选特定网络类型的负载均衡器
  network_type = "vpc"

  # 根据VPC ID筛选负载均衡器
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"

  # 根据交换机ID筛选负载均衡器
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"

  # 根据标签筛选负载均衡器
  tags = {
    Environment = "Test"
    Owner      = "Terraform"
  }

  # 输出结果保存到文件
  output_file = "slb_output.txt"
}

output "first_slb_id" {
  value = data.alibabacloudstack_slb_loadbalancers.default.slbs[0].id
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) SLB负载均衡器ID列表。可以通过此参数直接指定需要查询的SLB实例ID。
* `name_regex` - (可选，变更时重建) 用于按SLB名称筛选结果的正则表达式字符串。例如，使用`tf-SlbDataSourceSlbs.*`可以筛选出所有以`tf-SlbDataSourceSlbs`开头的负载均衡器。
* `master_availability_zone` - (可选，变更时重建) SLB的主要可用区。用于筛选在特定主要可用区中的负载均衡器。
* `slave_availability_zone` - (可选，变更时重建) SLB的次要可用区。用于筛选在特定次要可用区中的负载均衡器。
* `network_type` - (可选，变更时重建) SLB的网络类型。有效值：`vpc` 和 `classic`。默认值为`vpc`。
* `vpc_id` - (可选，变更时重建) 与SLB关联的VPC ID。用于筛选属于特定VPC的负载均衡器。
* `vswitch_id` - (可选，变更时重建) 与SLB关联的交换机ID。用于筛选属于特定交换机的负载均衡器。
* `address` - (可选，变更时重建) SLB的服务地址。用于筛选具有特定服务地址的负载均衡器。
* `tags` - (可选，变更时重建) 分配给SLB实例的标签映射。每个标签最多可以有5个键值对。例如：`{Environment="Test", Owner="Terraform"}`。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - SLB名称列表。包含所有匹配条件的SLB实例的名称。
* `slbs` - SLB列表。每个元素包含以下属性：
  * `id` - SLB的ID。
  * `region_id` - SLB所属的区域ID。
  * `master_availability_zone` - SLB的主要可用区。
  * `slave_availability_zone` - SLB的次要可用区。
  * `name` - SLB的名称。
  * `network_type` - SLB的网络类型。可能的值：`vpc` 和 `classic`。
  * `vpc_id` - SLB所属的VPC ID。
  * `vswitch_id` - SLB所属的交换机ID。
  * `address` - SLB的服务地址。
  * `creation_time` - SLB的创建时间。
  * `tags` - 分配给SLB的标签。