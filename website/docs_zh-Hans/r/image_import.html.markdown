---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_import"
sidebar_current: "docs-alibabacloudstack-resource-image-import"
description: |-
  编排导入本地的镜像文件
---

# alibabacloudstack_image_import

使用Provider配置的凭证在指定的资源集下编排导入本地的镜像文件到ECS中，作为自定义镜像出现在相应的域中。

-> **NOTE:** 您必须提前将镜像文件上传到对象存储OSS中。

-> **NOTE:** 导入镜像所在的地域必须与上传镜像文件的OSS Bucket所在的地域相同。

## 示例用法

```
resource "alibabacloudstack_image_import" "this" {
  description  = "test import image"
  architecture = "x86_64"
  image_name   = "test-import-image"
  license_type = "Auto"
  platform     = "Ubuntu"
  os_type      = "linux"
  disk_device_mapping {
    disk_image_size = 5
    oss_bucket      = "testimportimage"
    oss_object      = "root.img"
  }
}
```

## 参数参考

以下是支持的参数：

* `architecture` - (可选, 变更时重建) 指定在使用数据盘快照作为系统盘的数据源创建镜像时，系统盘的架构。有效值：`i386`，默认值为 `x86_64`。
* `description` - (可选) 镜像的描述。长度为2到256个英文或中文字符，不能以http: //和https: //开头。
* `image_name` - (可选) 镜像名称。长度为2 ~ 128个英文或中文字符。必须以英文字母或中文开头，不能以http: //和https: //开头。可以包含数字、冒号(:)、下划线(_)或连字符(-)。
* `license_type` - (可选, 变更时重建) 镜像导入后用于激活操作系统的许可证类型。默认值：`Auto`。有效值：`Auto`, `Aliyun`, `BYOL`。
* `platform` - (可选, 变更时重建) 指定在使用数据盘快照作为系统盘的数据源创建镜像时，系统盘的操作系统平台。有效值：`CentOS`, `Ubuntu`, `SUSE`, `OpenSUSE`, `Debian`, `CoreOS`, `Windows Server 2003`, `Windows Server 2008`, `Windows Server 2012`, `Windows 7`，默认值为 `Others Linux`, `Customized Linux`。
* `os_type` - (可选, 变更时重建) 操作系统平台类型。有效值：`windows`，默认值为 `linux`。
* `disk_device_mapping` - (必填, 变更时重建) 镜像下的系统盘和快照的描述。
  * `device` - (可选, 变更时重建) 自定义镜像中的磁盘N的名称。
  * `disk_image_size` - (可选, 变更时重建) 分辨率大小。必须确保系统盘空间 ≥ 文件系统空间。范围：当n = 1时，系统盘：5 ~ 500GiB，当n = 2 ~ 17时，即数据盘：5 ~ 1000GiB，当临时导入时，系统自动检测大小，以检测结果为准。
  * `format` - (可选, 变更时重建) 镜像格式。值范围：当导入 `RAW`, `VHD`, `qcow2` 格式的镜像时，系统会自动检测镜像格式，以最先检测到的为准。
  * `oss_bucket` - (可选) 保存导出的OSS桶。
  * `oss_object` - (可选, 变更时重建) OSS对象的文件名。

-> **NOTE:** `disk_device_mapping` 是一个列表，它的第一个元素将被用作系统盘，其他元素将被用作数据盘。

## 超时时间

`timeouts` 块允许你指定某些动作的 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为20分钟) 用于复制镜像(直到它达到初始 `Available` 状态)。
* `delete` - (默认为20分钟) 用于终止镜像。

## 属性参考

以下属性将会导出：

* `id` - 镜像ID。
* `format` - 镜像格式。值范围：当导入 `RAW`, `VHD`, `qcow2` 格式的镜像时，系统会自动检测镜像格式，以最先检测到的为准。