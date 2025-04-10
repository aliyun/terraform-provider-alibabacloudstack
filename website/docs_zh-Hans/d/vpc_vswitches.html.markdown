---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vswitches"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-vswitches"
description: |- 
  查询专有网络（VPC）虚拟交换机
---

# alibabacloudstack_vpc_vswitches
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vswitches`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）虚拟交换机列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccVSwitchDatasource13526"
}

data "alibabacloudstack_zones" "default" {}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name       = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default" {
  vswitch_name   = "${var.name}"
  cidr_block     = "172.16.0.0/24"
  vpc_id         = "${alibabacloudstack_vpc.default.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

data "alibabacloudstack_vpc_vswitches" "default" {
  name_regex = "${alibabacloudstack_vswitch.default.vswitch_name}"
  vpc_id     = "${alibabacloudstack_vpc.default.id}"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"

  output_file = "vswitches_output.txt"
}

output "vswitch_names" {
  value = data.alibabacloudstack_vpc_vswitches.default.names
}
```

## 参数参考

以下参数是支持的：

* `cidr_block` - (可选) 交换机的网段。用于筛选具有特定CIDR块的vSwitch。
* `name_regex` - (可选) 用于通过正则表达式筛选vSwitch名称的结果。
* `is_default` - (可选，类型：bool) 指定是否查询指定区域中的默认vSwitch。有效值：
  * **true** - 仅查询默认vSwitch。
  * **false** - 从查询中排除默认vSwitch。
  如果不设置此参数，默认情况下系统将查询指定区域中的所有vSwitch。
* `vpc_id` - (可选) vSwitch所属虚拟私有云(VPC)的ID。至少需要指定`vpc_id`或`zone_id`中的一个。
* `zone_id` - (可选) vSwitch所在可用区的ID。您可以调用[DescribeZones](https://help.aliyun.com/document_detail/36064.html)操作来查询最新的可用区列表。
* `ids` - (可选) 用于过滤结果的vSwitch ID列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - vSwitch ID列表。
* `names` - vSwitch名称列表。
* `vswitches` - vSwitch详细信息列表。每个元素包含以下属性：
  * `id` - vSwitch的ID。
  * `vpc_id` - vSwitch所属VPC的ID。
  * `zone_id` - vSwitch所在可用区的ID。
  * `vswitch_name` - vSwitch的名称。
  * `instance_ids` - 与该vSwitch关联的ECS实例ID列表。
  * `cidr_block` - vSwitch的CIDR块。
  * `description` - vSwitch的描述。
  * `is_default` - 指示该vSwitch是否为该区域中的默认vSwitch。
  * `creation_time` - vSwitch创建的时间。
  * `available_ip_address_count` - vSwitch中可用IP地址的数量。