---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_routetables"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-routetables"
description: |- 
  查询专有网络（VPC）网络路由表
---

# alibabacloudstack_vpc_routetables
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_route_tables`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）网络路由表列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccRouteTablesDatasource19531"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

resource "alibabacloudstack_route_table" "default" {
  vpc_id      = "${alibabacloudstack_vpc.default.id}"
  name        = "${var.name}"
  description = "${var.name}_description"
}

data "alibabacloudstack_vpc_routetables" "default" {
  vpc_id      = "${alibabacloudstack_vpc.default.id}"
  ids         = ["${alibabacloudstack_route_table.default.id}"]
  name_regex  = "${alibabacloudstack_route_table.default.name}"
  tags        = {
    Environment = "Test"
  }
  output_file = "route_tables_output.txt"
}

output "route_table_ids" {
  value = "${data.alibabacloudstack_vpc_routetables.default.ids}"
}
```

## 参数参考

以下参数是支持的：

* `vpc_id` - (选填)VPC的ID。如果指定，将仅返回属于该VPC的路由表。
* `ids` - (选填)路由表ID列表。如果指定，将仅返回这些ID对应的路由表。
* `name_regex` - (选填)用于按名称过滤路由表的正则表达式字符串。
* `tags` - (选填)标签映射，每个标签表示为键值对。通过此参数可以筛选具有特定标签的路由表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 路由表ID列表。
* `names` - 路由表名称列表。
* `tables` - 路由表列表。每个元素包含以下属性：
  * `id` - 路由表的ID。
  * `router_id` - 路由表所属的路由器ID。
  * `route_table_type` - 路由表的类型。取值范围：
    - `custom`: 自定义路由表。
    - `system`: 系统路由表。
  * `name` - 路由表的名称。
  * `description` - 路由表的描述。
  * `creation_time` - 路由表创建的时间。