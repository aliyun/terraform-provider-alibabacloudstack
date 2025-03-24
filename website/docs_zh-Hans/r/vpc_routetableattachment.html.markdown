---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetableattachment"
sidebar_current: "docs-Alibabacloudstack-vpc-routetableattachment"
description: |- 
  编排绑定交换机和路由表
---

# alibabacloudstack_vpc_routetableattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_route_table_attachment`

使用Provider配置的凭证在指定的资源集编排绑定交换机和路由表。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccvpcroute_table_attachment24025"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_route_table" "default" {
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_description"
}

resource "alibabacloudstack_vpc_routetableattachment" "default" {
  route_table_id = "${alibabacloudstack_route_table.default.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
```

## 参数参考

支持以下参数：

* `route_table_id` - (必填，变更时重建) 要绑定到交换机的路由表的ID。此字段创建后无法修改。
* `vswitch_id` - (必填，变更时重建) 要绑定路由表的交换机的ID。此字段创建后无法修改。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 路由表绑定的唯一标识符。格式为 `<route_table_id>:<vswitch_id>`。
