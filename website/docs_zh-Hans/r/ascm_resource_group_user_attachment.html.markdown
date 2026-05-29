---
subcategory: "应用"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_group_user_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ascm-resource-group-user-attachment"
description: |-
  提供 ASCM 资源组用户绑定资源。
---

# alibabacloudstack_ascm_resource_group_user_attachment

提供 ASCM 资源组用户绑定资源。该资源用于将用户绑定到 ASCM 资源组。

-> **注意:** 修改任何参数都会强制创建新资源。

## 示例

```
resource "alibabacloudstack_ascm_resource_group_user_attachment" "default" {
  user_id = alibabacloudstack_ascm_user.user.user_id
  rg_id   = alibabacloudstack_ascm_resource_group.default.rg_id
}
```

## 参数说明

以下参数用于配置资源组用户绑定：

* `user_id` - （必填，ForceNew）要绑定到资源组的用户 ID。修改此参数将强制创建新资源。
* `rg_id` - （可选，ForceNew）资源组 ID。修改此参数将强制创建新资源。

## 属性说明

以下属性由系统导出：

* `id` - 资源组用户绑定的 ID。格式为 `<rg_id>:<user_id>`。

## 导入

资源组用户绑定可以通过资源组 ID 和用户 ID（以冒号分隔）导入，例如：

```
$ terraform import alibabacloudstack_ascm_resource_group_user_attachment.example 12345:67890
```
