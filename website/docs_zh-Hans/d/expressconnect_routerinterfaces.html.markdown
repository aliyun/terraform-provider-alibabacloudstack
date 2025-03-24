---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_routerinterfaces"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-routerinterfaces"
description: |- 
  查询高速通道路由器接口
---

# alibabacloudstack_expressconnect_routerinterfaces
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_router_interfaces`

根据指定过滤条件列出当前凭证权限可以访问的高速通道路由器接口列表。

## 示例用法

```hcl
variable "region" {
  default = "cn-wulan-env200-d01"
}
variable "name" {
  default = "tf-testAccCheckAlibabacloudStackRouterInterfacesDataSourceConfig"
}
variable cidr_block_list {
  type = "list"
  default = [ "172.16.0.0/12", "192.168.0.0/16" ]
}

resource "alibabacloudstack_vpc" "default" {
  count = 2
  name = "${var.name}"
  cidr_block = "${element(var.cidr_block_list, count.index)}"
}

resource "alibabacloudstack_router_interface" "initiating" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alibabacloudstack_vpc.default[0].router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}_initiating"
  description = "${var.name}_decription"
}

resource "alibabacloudstack_router_interface" "opposite" {
  provider = "alibabacloudstack"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alibabacloudstack_vpc.default[1].router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}_opposite"
  description = "${var.name}_decription"
}

resource "alibabacloudstack_router_interface_connection" "foo" {
  interface_id = "${alibabacloudstack_router_interface.initiating.id}"
  opposite_interface_id = "${alibabacloudstack_router_interface.opposite.id}"
  depends_on = ["alibabacloudstack_router_interface_connection.bar"]
  opposite_interface_owner_id = "1262302482727553"
  opposite_router_id = "${alibabacloudstack_vpc.default[0].router_id}"
  opposite_router_type = "VRouter"
}

resource "alibabacloudstack_router_interface_connection" "bar" {
  provider = "alibabacloudstack"
  interface_id = "${alibabacloudstack_router_interface.opposite.id}"
  opposite_interface_id = "${alibabacloudstack_router_interface.initiating.id}"
  opposite_interface_owner_id = "1262302482727553"
  opposite_router_id = "${alibabacloudstack_vpc.default[1].router_id}"
  opposite_router_type = "VRouter"
}

data "alibabacloudstack_expressconnect_routerinterfaces" "default" {
  ids = ["${alibabacloudstack_router_interface.initiating.id}"]
}

output "first_router_interface_id" {
  value = "${data.alibabacloudstack_expressconnect_routerinterfaces.default.interfaces.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `status` - (可选，变更时重建) - 路由器接口的状态。可能的值为：`Active`、`Inactive` 和 `Idle`。
* `name_regex` - (可选，变更时重建) - 用于按路由器接口名称筛选的正则表达式字符串。
* `specification` - (可选，变更时重建) - 链路规格，例如 `Small.1`(10Mb)、`Middle.1`(100Mb)、`Large.2`(2Gb)等。
* `router_id` - (可选，变更时重建) - 位于本地区域的 VRouter 的 ID。
* `router_type` - (可选，变更时重建) - 本地区域中的路由器类型。可能的值为：`VRouter` 和 `VBR`(物理连接)。
* `role` - (可选，变更时重建) - 路由器接口的角色。可能的值为：`InitiatingSide`(连接发起方)和 `AcceptingSide`(连接接收方)。如果将 `router_type` 设置为 `VBR`，则此参数的值必须为 `InitiatingSide`。
* `opposite_interface_id` - (可选，变更时重建) - 对等路由器接口的 ID。
* `opposite_interface_owner_id` - (可选，变更时重建) - 对等路由器接口所有者的帐户 ID。
* `ids` - (可选) - 路由器接口 ID 列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 路由器接口 ID 列表。
* `names` - 路由器接口名称列表。
* `interfaces` - 路由器接口列表。每个元素包含以下属性：
  * `id` - 路由器接口 ID。
  * `status` - 路由器接口状态。可能的值为：`Active`、`Inactive` 和 `Idle`。
  * `name` - 路由器接口名称。
  * `description` - 路由器接口描述。
  * `role` - 路由器接口角色。可能的值为：`InitiatingSide` 和 `AcceptingSide`。
  * `specification` - 路由器接口规格。可能的值为：`Small.1`、`Middle.1`、`Large.2` 等。
  * `router_id` - 位于本地区域的 VRouter 的 ID。
  * `router_type` - 本地区域中的路由器类型。可能的值为：`VRouter` 和 `VBR`。
  * `vpc_id` - 拥有本地区域中路由器的 VPC 的 ID。
  * `access_point_id` - 由 VBR 使用的接入点 ID。
  * `creation_time` - 路由器接口创建时间。
  * `opposite_region_id` - 对等路由器区域 ID。
  * `opposite_interface_id` - 对等路由器接口 ID。
  * `opposite_router_id` - 对等路由器 ID。
  * `opposite_router_type` - 对等区域中的路由器类型。可能的值为：`VRouter` 和 `VBR`。
  * `opposite_interface_owner_id` - 对等路由器接口所有者的帐户 ID。
  * `health_check_source_ip` - 用于对物理连接执行健康检查的源 IP 地址。
  * `health_check_target_ip` - 用于对物理连接执行健康检查的目标 IP 地址。