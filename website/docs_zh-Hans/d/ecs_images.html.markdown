---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_images"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-images"
description: |- 
  查询云服务器镜像
---

# alibabacloudstack_ecs_images
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_images`

根据指定过滤条件列出当前凭证权限可以访问的云服务器镜像列表。

## 示例用法

```hcl
data "alibabacloudstack_ecs_images" "default" {
  owners     = "system"
  name_regex = "^centos_6"
  most_recent = true
}

output "first_image_id" {
  value = "${data.alibabacloudstack_ecs_images.default.images.0.image_id}"
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (选填, 变更时重建) 用于按名称过滤结果的正则表达式字符串。例如，可以使用 `"^centos_6"` 来筛选所有以 `centos_6` 开头的镜像。
* `most_recent` - (选填, 强制新建, 类型：布尔值) 如果返回多个结果，则选择最近创建的一个镜像。默认值为 `false`。
* `owners` - (选填, 变更时重建) 按特定镜像所有者筛选结果。有效项为：
  * `system` - 系统公共镜像。
  * `self` - 用户的自定义镜像。
  * `others` - 其他用户的公开镜像。
  * `marketplace` - 镜像市场镜像。

-> **注意:** 至少需要设置 `name_regex`、`most_recent` 和 `owners` 中的一个。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 镜像ID列表。
* `images` - 镜像列表。每个元素包含以下属性：
  * `id` - 镜像的ID(与 `image_id` 相同)。
  * `image_id` - 镜像的唯一标识符。
  * `architecture` - 镜像系统的平台类型：`i386` 或 `x86_64`。
  * `creation_time` - 镜像创建的时间，格式为 `YYYY-MM-DDTHH:mm:ssZ`。
  * `description` - 镜像的描述信息。
  * `image_owner_alias` - 镜像所有者的别名。有效值包括：
    * `system` - 系统公共镜像。
    * `self` - 用户的自定义镜像。
    * `others` - 其他用户的公开镜像。
    * `marketplace` - 镜像市场镜像。
  * `os_name` - 操作系统的中文显示名称。
  * `os_name_en` - 操作系统的英文显示名称。
  * `os_type` - 镜像的操作系统类型。有效值为 `windows` 和 `linux`。
  * `platform` - 镜像的操作系统平台。
  * `status` - 镜像的状态。可能的值包括：
    * `UnAvailable` - 不可用。
    * `Available` - 可用。
    * `Creating` - 创建中。
    * `CreateFailed` - 创建失败。
  * `state` - 镜像的状态(与 `status` 相同)。
  * `size` - 镜像的大小，单位为GiB。
  * `disk_device_mappings` - 镜像的快照信息。每个映射包括：
    * `device` - 创建磁盘的设备信息，例如 `/dev/xvdb`。
    * `size` - 创建磁盘的大小，单位为GiB。
    * `snapshot_id` - 与磁盘关联的快照ID。
  * `product_code` - 图像市场上的图像产品代码。
  * `is_subscribed` - 用户是否订阅了与 `product_code` 对应的图像产品的服务条款。
  * `is_copied` - 是否为复制的镜像。有效值为 `true` 或 `false`。
  * `is_self_shared` - 自定义镜像是否已与其他用户共享。有效值为 `true` 或 `false`。
  * `image_version` - 镜像的版本。
  * `progress` - 镜像创建的进度，以百分比表示。
  * `usage` - 指定是否在不实际发出请求的情况下检查请求的有效性。可能值为：
    * `instance` - 表示该镜像正在被实例使用。
    * `none` - 表示该镜像未被任何资源引用。
  * `is_support_io_optimized` - 镜像是否可以在I/O优化实例上使用。
  * `tags` - 资源的标签。