---
subcategory: "专有网络 VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_router_interface"
sidebar_current: "docs-alibabacloudstack-resource-router-interface"
description: |-
  编排路由器接口
---

# alibabacloudstack_router_interface

使用Provider配置的凭证在指定的资源集编排路由器接口。路由器接口用于建立专有网络（VPC）与本地数据中心之间的高速连接。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccRouterInterface"
}

data "alibabacloudstack_account" "current" {
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_router_interface" "accepting_side" {
  opposite_region           = data.alibabacloudstack_account.current.region
  router_type               = "VRouter"
  router_id                 = alibabacloudstack_vpc.default.router_id
  role                      = "AcceptingSide"
  name                      = var.name
  description               = var.name
}

resource "alibabacloudstack_router_interface" "initiating_side" {
  opposite_region           = data.alibabacloudstack_account.current.region
  router_type               = "VRouter"
  router_id                 = alibabacloudstack_vpc.default.router_id
  role                      = "InitiatingSide"
  specification             = "Large.2"
  name                      = "${var.name}-initiating"
  description               = "${var.name}-initiating"
  opposite_interface_id     = alibabacloudstack_router_interface.accepting_side.id
  opposite_router_id        = alibabacloudstack_vpc.default.router_id
  opposite_router_type      = "VRouter"
  opposite_interface_owner_id = data.alibabacloudstack_account.current.id
}
```

## 参数说明

以下是支持的参数：

### 必填参数

* `opposite_region` - (必填，变更时重建) 另一侧路由器接口所在的地域 ID。如果是跨地域连接，需要指定对端地域。
* `router_type` - (必填，变更时重建) 路由器类型。可选值：`VRouter`（专有网络路由器），`VBR`（边界路由器）。
* `router_id` - (必填，变更时重建) 路由器 ID。当 router_type 为 VRouter 时，该值为 VPC 的 ID；当 router_type 为 VBR 时，该值为 VBR 的 ID。
* `role` - (必填，变更时重建) 路由器接口角色。可选值：`InitiatingSide`（发起端），`AcceptingSide`（接受端）。

### 可选参数

* `specification` - (可选) 路由器接口规格。当 role 为 InitiatingSide 时必填，当 role 为 AcceptingSide 时不需要设置。可选值根据地域和可用区可能有所不同，常见值：`Large.1`、`Large.2`、`Medium.1`、`Medium.2`、`Small.1`、`Small.2` 等。
* `name` - (可选) 路由器接口名称。长度为 2-128 个字符，必须以字母或中文开头，不能以 `http://` 和 `https://` 开头。
* `description` - (可选) 路由器接口描述。长度为 2-256 个字符。
* `health_check_source_ip` - (可选) 健康检查源 IP 地址。当 router_type 为 VRouter 时有效。需要与 `health_check_target_ip` 同时设置。
* `health_check_target_ip` - (可选) 健康检查目标 IP 地址。当 router_type 为 VRouter 时有效。需要与 `health_check_source_ip` 同时设置。
* `opposite_access_point_id` - (可选) 另一侧接入点 ID。当 router_type 为 VBR 时有效。
* `opposite_router_type` - (可选) 另一侧路由器类型。可选值：`VRouter`、`VBR`。
* `opposite_router_id` - (可选) 另一侧路由器 ID。
* `opposite_interface_id` - (可选) 另一侧路由器接口 ID。
* `opposite_interface_owner_id` - (可选) 另一侧路由器接口账号 ID。如果不指定，则默认为当前账号。

-> **注意:** 
- 当 role 为 `AcceptingSide` 时，`specification` 参数会自动设置为 `Negative`，无需手动指定。
- `health_check_source_ip` 和 `health_check_target_ip` 必须同时设置或同时不设置。
- `opposite_interface_owner_id` 的值必须是主账号 ID，而不是子账号。

## 属性说明

以下属性将会被导出：

* `id` - 路由器接口 ID。
* `access_point_id` - 接入点 ID。当 router_type 为 VBR 时有效。
* `opposite_access_point_id` - 另一侧接入点 ID。
* `opposite_router_type` - 另一侧路由器类型。
* `opposite_router_id` - 另一侧路由器 ID。
* `opposite_interface_id` - 另一侧路由器接口 ID。
* `opposite_interface_owner_id` - 另一侧路由器接口账号 ID。
* `status` - 路由器接口状态。常见状态：`Idle`（空闲）、`Waiting`（等待中）、`Active`（激活）、`Inactive`（未激活）。

## Import

路由器接口可以使用 router_interface_id 导入，例如：

```
$ terraform import alibabacloudstack_router_interface.example ri-12345678
```

-> **注意:** 该资源也支持以下别名：
> - `alibabacloudstack_expressconnect_routerinterface`
