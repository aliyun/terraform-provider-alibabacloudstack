---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_export"
sidebar_current: "docs-alibabacloudstack-resource-image-export"
description: |-
  编排导出自定义镜像到OSS的资源
---

# alibabacloudstack_image_export

使用Provider配置的凭证在指定的资源集下编排导出自定义镜像到OSS的资源。

-> **注意:** 如果使用镜像创建了ECS实例，再次从系统盘创建了系统盘快照，则不支持导出由该系统盘快照创建的自定义镜像。

-> **注意:** 支持导出包含数据盘快照信息的自定义镜像。镜像中的数据盘数量不能超过4块，单块数据盘容量最大不能超过500 GiB。

-> **注意:** 在导出镜像之前，必须通过RAM授权云服务器ECS官方服务账号写OSS权限。

## 示例用法

```
resource "alibabacloudstack_image_export" "default" {
  image_id           = "m-bp1gxy***"
  oss_bucket         = "ecsimageexportconfig"
  oss_prefix         = "ecsExport"
}
```

## 参数说明

以下是支持的参数：

* `image_id` - (必填，变更时重建) 源镜像ID。
* `oss_bucket` - (必填，变更时重建) 保存导出文件的OSS Bucket名称。
* `oss_prefix` - (可选，变更时重建) OSS Object的前缀。可以由数字或字母组成，字符长度为1 ~ 30。

## 超时时间

`timeouts` 块允许你为某些操作指定 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认20分钟)用于导出镜像(直到其达到初始 `Available` 状态)。

## 属性说明

以下属性将会被导出：

* `id` - 镜像的ID。