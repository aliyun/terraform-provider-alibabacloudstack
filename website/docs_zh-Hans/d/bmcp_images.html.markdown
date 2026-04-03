---
subcategory: "BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_images"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-images"
description: |-
  查询裸金属计算平台(BMCP)镜像
---

# alibabacloudstack_bmcp_images

根据指定过滤条件列出当前凭证权限可以访问的BMCP镜像列表。

## 示例用法

```hcl
data "alibabacloudstack_bmcp_images" "default" {
  type = "public"
}
```

### 按名称正则表达式过滤

```hcl
data "alibabacloudstack_bmcp_images" "ubuntu_images" {
  type       = "public"
  name_regex = "ubuntu"
}
```

### 按ID过滤

```hcl
data "alibabacloudstack_bmcp_images" "filtered_images" {
  type = "public"
  ids  = ["i-d"]
}
```

### 组合过滤条件

```hcl
data "alibabacloudstack_bmcp_images" "combined" {
  type       = "public"
  name_regex = "ubuntu"
  ids        = ["i-d"]
}
```

## 参数说明

以下参数是支持的：

* `type` - (必填) 镜像类型。有效值：`customize`（自定义镜像）、`public`（公共镜像）。
* `name_regex` - (选填) 用于按名称过滤结果的正则表达式字符串。
* `ids` - (选填) 用于过滤的镜像ID列表。支持模糊匹配。

## 属性说明

除了上述参数外，还导出以下属性：

* `images` - 镜像列表。每个元素包含以下属性：
  * `id` - 镜像的ID（与`name`相同）。
  * `name` - 镜像的名称。
  * `description` - 镜像的描述信息。
  * `type` - 镜像的类型。
  * `platform` - 镜像的平台。
  * `zone` - 镜像的可用区。
  * `region` - 镜像的地域。
  * `status` - 镜像的状态。
  * `size` - 镜像的大小。
  * `create_time` - 镜像的创建时间。
  * `update_time` - 镜像的更新时间。
  * `os_arch` - 镜像的操作系统架构。
  * `disk_format` - 镜像的磁盘格式。
  * `package_type` - 镜像的包类型。
  * `url` - 镜像的URL。
  * `unique_key` - 镜像的唯一标识。
  * `template` - 镜像的模板。
  * `resource_group` - 镜像的资源组。
  * `file_name` - 镜像的文件名。
  * `snapshot_id` - 镜像的快照ID。
  * `error_description` - 镜像的错误描述。
  * `part` - 镜像的分区信息。
  * `reboot_script` - 镜像的重启脚本。
  * `signature_file_url` - 镜像的签名文件URL。
  * `owner` - 镜像的所有者。
  * `component_type` - 镜像的组件类型。
  * `version` - 镜像的版本。
  * `custom_script` - 镜像的自定义脚本。
  * `deleted` - 镜像是否已删除。
  * `disable` - 镜像是否已禁用。
  * `kernel_version` - 镜像的内核版本。
  * `organization` - 镜像的组织。
  * `check_sum` - 镜像的校验和。
