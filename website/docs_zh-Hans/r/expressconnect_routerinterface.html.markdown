---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_routerinterface"
sidebar_current: "docs-Alibabacloudstack-expressconnect-routerinterface"
description: |- 
  编排高速通道路由器接口
---

# alibabacloudstack_expressconnect_routerinterface
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_router_interface`

使用Provider配置的凭证在指定的资源集下编排高速通道路由器接口。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccRouterInterfaceConfig18787"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alibabacloudstack_router_interface" "default" {
  opposite_region      = var.region
  router_type          = "VRouter"
  router_id            = alibabacloudstack_vpc.default.router_id
  role                 = "AcceptingSide"
  specification        = "Large.2"
  name                = var.name
  description         = "This is a test router interface"
  health_check_source_ip = "172.16.0.10"
  health_check_target_ip = "172.16.0.11"
  opposite_access_point_id = "ap-abc123456"
}
```

## 参数参考

支持以下参数：

* `opposite_region` - (必填, 变更时重建) - 对端的区域。
* `router_type` - (必填, 变更时重建) - 路由器类型。有效值：`VRouter`(VPC路由器)、`VBR`(边界路由器)。
* `router_id` - (必填, 变更时重建) - 路由器接口所属的路由器 ID。
* `role` - (必填, 变更时重建) - 路由器接口的角色。有效值：`InitiatingSide`(发起方)、`AcceptingSide`(接收方)。
* `specification` - (选填) - 路由器接口的规格。当 `role` 设置为 `InitiatingSide` 时有效。可能的值包括：`Small.1`，`Middle.1`，`Large.2`。
* `name` - (选填) - 路由器接口的名称。长度必须在 2 到 80 个字符之间。仅允许中文字符、英文字母、数字、句点(`.`)、下划线(`_`)或连字符(`-`)。如果未指定，默认为路由器接口 ID。名称不能以 `http://` 或 `https://` 开头。
* `description` - (选填) - 路由器接口的描述。长度必须在 2 到 256 个字符之间，或者留空。它不能以 `http://` 或 `https://` 开头。
* `health_check_source_ip` - (选填) - 健康检查包的源 IP 地址。仅在 `router_type` 设置为 `VBR` 时有效。该 IP 地址必须是本地 VPC 子网中未使用的 IP。必须与 `health_check_target_ip` 一起指定。
* `health_check_target_ip` - (选填) - 健康检查包的目标 IP 地址。仅在 `router_type` 设置为 `VBR` 时有效。该 IP 地址必须是本地 VPC 子网中未使用的 IP。必须与 `health_check_source_ip` 一起指定。
* `opposite_access_point_id` - (选填) - 对端路由器接口的接入点 ID。
* `opposite_router_type` - (选填) - 对端路由器的类型。有效值：`VRouter`，`VBR`。
* `opposite_router_id` - (选填) - 对端路由器的 ID。
* `opposite_interface_id` - (选填) - 对端路由器接口的 ID。
* `opposite_interface_owner_id` - (选填) - 对端路由器接口的所有者账户 ID。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 路由器接口的 ID。
* `access_point_id` - 路由器接口的接入点 ID。
* `opposite_router_type` - 对端路由器的类型。
* `opposite_router_id` - 对端路由器的 ID。
* `opposite_interface_id` - 对端路由器接口的 ID。
* `opposite_interface_owner_id` - 对端路由器接口的所有者账户 ID。
* `health_check_source_ip` - 健康检查包的源 IP 地址。
* `health_check_target_ip` - 健康检查包的目标 IP 地址。