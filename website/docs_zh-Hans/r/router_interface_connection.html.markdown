---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_router_interface_connection"
sidebar_current: "docs-alibabacloudstack-resource-route-interface-connection"
description: |-
  编排路由器接口连接
---

# alibabacloudstack_router_interface_connection

使用Provider配置的凭证在指定的资源集编排路由器接口连接。

## 示例用法
```
resource "alibabacloudstack_router_interface_connection" "foo" {
    interface_id = "${alibabacloudstack_router_interface.foo.id}"
    opposite_interface_id = "${alibabacloudstack_router_interface.bar.id}"
    opposite_interface_owner_id = "${alibabacloudstack_router_interface.bar.owner_account_id}"
}
```

## 参数说明

以下是支持的参数：

* `interface_id` - (必填，变更时重建) 一侧路由器接口 ID。
* `opposite_interface_id` - (必填，变更时重建) 另一侧路由器接口 ID。它必须属于指定的 "opposite_interface_owner_id" 账户。
* `opposite_interface_owner_id` - (可选，变更时重建) 另一侧路由器接口账户 ID。登录 AlibabacloudStack 控制台，选择用户信息 > 账户管理以查看账户 ID。默认为 [Provider account_id](https://www.terraform.io/docs/providers/alibabacloudstack/index.html#account_id)。
* `opposite_router_id` - (可选，变更时重建) 另一侧路由器 ID。它必须属于指定的 "opposite_interface_owner_id" 账户。当字段 "opposite_interface_owner_id" 被指定时有效。
* `opposite_router_type` - (可选，变更时重建) 另一侧路由器类型。可选值：VRouter, VBR。当字段 "opposite_interface_owner_id" 被指定时有效。

-> **注意:** "opposite_interface_owner_id" 或 "account_id" 的值必须是主账户，而不是子账户。


## 属性说明

以下属性将会被导出：

* `id` - 路由器接口 ID。其值等于 "interface_id"。
* `opposite_router_id` - 另一侧路由器 ID。它必须属于指定的 "opposite_interface_owner_id" 账户。当字段 "opposite_interface_owner_id" 被指定时有效。
* `opposite_interface_owner_id` - 另一侧路由器接口账户 ID。