---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_copy"
sidebar_current: "docs-alibabacloudstack-resource-image-copy"
description: |-
  编排实现自定义镜像从一个区域复制到另一个区域
---

# alibabacloudstack_image_copy

使用Provider配置的凭证在指定的资源集下编排实现自定义镜像从一个区域复制到另一个区域。
您可以在目标区域中使用复制的镜像执行操作，例如创建实例(RunInstances)和更换系统盘(ReplaceSystemDisk)。

-> **注意：** 只有当镜像处于可用状态时，才能复制自定义镜像。

-> **注意：** 只能复制属于您的阿里云账号的镜像。镜像不能从一个账户复制到另一个账户。

-> **注意：** 如果复制未完成，您不能调用DeleteImage删除镜像，但可以调用CancelCopyImage取消复制。

## 示例用法

```
resource "alibabacloudstack_image_copy" "default" {
  source_image_id    = "m-bp1gxyhdswlsn18tu***"
  source_region_id   = "cn-hangzhou"
  image_name         = "test-image"
  description        = "test-image"
  tags               = {
         FinanceDept = "FinanceDeptJoshua"
     }
}
```

## 参数参考

支持以下参数：

* `source_image_id` - (必填，强制新值)源镜像ID。
* `destination_region_id` - (必填，强制新值)目标区域ID。
* `name` - (可选，已弃用)字段'name'已被弃用，改为使用新字段'image_name'。
* `image_name` - (可选) 镜像名称。它必须是2到128个字符长度，并且必须以字母或中文字符开头(不能以http://或https://开头)。它可以包含数字、冒号(:)、下划线(_)或连字符(-)。默认值：null。
* `description` - (可选) 镜像描述。它必须是2到256个字符长度，并且不能以http://或https://开头。默认值：null。
* `kms_key_id` - (可选，强制新值)用于加密的KMS密钥ID。
* `encrypted` - (可选，强制新值)指示镜像是否加密。

## 属性参考

导出以下属性：

* `id` - 镜像的ID。
* `name` - 镜像名称。
* `image_name` - 镜像名称。