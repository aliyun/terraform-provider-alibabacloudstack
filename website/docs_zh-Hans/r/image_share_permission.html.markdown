---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_share_permission"
sidebar_current: "docs-alibabacloudstack-resource-image-share-permission"
description: |-
  编排管理ECS镜像共享权限
---

# alibabacloudstack_image_share_permission

使用Provider配置的凭证在指定的资源集下编排管理ECS镜像共享权限。
您可以将自定义镜像共享给其他阿里云用户。用户可以使用共享的自定义镜像创建ECS实例或替换实例的系统盘。

-> **NOTE:** 您只能将您自己的自定义镜像共享给其他阿里云用户。

-> **NOTE:** 每个自定义镜像最多可以共享给50个阿里云账户。您可以提交工单以共享给更多用户。

-> **NOTE:** 使用共享镜像创建ECS实例后，一旦自定义镜像的所有者解除镜像共享关系或删除自定义镜像，实例将无法初始化系统盘。

## 示例用法

```
resource "alibabacloudstack_image_share_permission" "default" {
  image_id           = "m-bp1gxyh***"
  account_id         = "1234567890"
}
```

## 参数参考

支持以下参数：

* `image_id` - (Required, ForceNew) 源镜像ID。
* `account_id` - (Required, ForceNew) 阿里云账号ID。用于共享镜像。

## 属性参考

导出以下属性：

* `id` - 镜像ID。格式为 `<image_id>:<account_id>`
* `image_id` - 源镜像ID。 
* `account_id` - 阿里云账号ID。用于共享镜像。