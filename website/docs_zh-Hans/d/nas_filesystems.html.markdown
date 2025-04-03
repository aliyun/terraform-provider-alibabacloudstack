---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_filesystems"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-filesystems"
description: |- 
  文件存储（NAS）文件系统列表
---

# alibabacloudstack_nas_filesystems
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_file_systems`

根据指定过滤条件列出当前凭证权限可以访问的文件存储（NAS）文件系统规则列表

## 示例用法

```terraform
variable "description" {
  default = "tf-testAccCheckAlibabacloudStackFileSystemsDataSource"
}
variable "storage_type" {
  default = "Capacity"
}

data "alibabacloudstack_nas_protocols" "default" {
  type = "${var.storage_type}"
}

resource "alibabacloudstack_nas_file_system" "default" {
  description  = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}

data "alibabacloudstack_nas_file_systems" "default" {
  description_regex = "^${alibabacloudstack_nas_file_system.default.description}"
  storage_type      = "${alibabacloudstack_nas_file_system.default.storage_type}"
}

output "file_system_id" {
  value = "${data.alibabacloudstack_nas_file_systems.default.systems.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `storage_type` - (选填, 变更时重建) 存储类型：
  * 当 `file_system_type = standard` 时，取值为 `Performance`(性能型)、`Capacity`(容量型)和 `Premium`(高级型)。
  * 当 `file_system_type = extreme` 时，取值为 `standard`(标准型)或 `advance`(高级型)。
  * 当 `file_system_type = cpfs` 时，取值为 `advance_100`(100 MB/s/TiB 基线)和 `advance_200`(200 MB/s/TiB 基线)。

* `protocol_type` - (选填, 变更时重建) 文件传输协议类型：
  * 当 `file_system_type = standard` 时，取值为 `NFS` 和 `SMB`。
  * 当 `file_system_type = extreme` 时，取值为 `NFS`。
  * 当 `file_system_type = cpfs` 时，取值为 `cpfs`。

* `description_regex` - (选填, 变更时重建) 用于通过文件系统描述筛选结果的正则表达式字符串。

* `ids` - (选填) 文件系统 ID 列表。


* `file_system_type` - (选填) 文件系统类型：
  * `standard`(默认值)：通用 NAS。
  * `extreme`：极速 NAS。
  * `cpfs`：文件存储 CPFS。

## 属性参考

除了上述参数外，还导出以下属性：

* `descriptions` - 文件系统描述列表。

* `systems` - 文件系统列表。每个元素包含以下属性：
  * `id` - 文件系统的 ID。
  * `region_id` - 文件系统所在的区域 ID。
  * `create_time` - 文件系统的创建时间。
  * `description` - 文件系统的描述。限制：
    * 长度为 2~128 个英文或中文字符。
    * 必须以大小写字母或中文开头，不能以 `http://` 或 `https://` 开头。
    * 可以包含数字、冒号(:)、下划线(_)或短横线(-)。
  * `protocol_type` - 文件传输协议类型：
    * 当 `file_system_type = standard` 时，取值为 `NFS` 和 `SMB`。
    * 当 `file_system_type = extreme` 时，取值为 `NFS`。
    * 当 `file_system_type = cpfs` 时，取值为 `cpfs`。
  * `storage_type` - 存储类型：
    * 当 `file_system_type = standard` 时，取值为 `Performance`(性能型)、`Capacity`(容量型)和 `Premium`(高级型)。
    * 当 `file_system_type = extreme` 时，取值为 `standard`(标准型)或 `advance`(高级型)。
    * 当 `file_system_type = cpfs` 时，取值为 `advance_100`(100 MB/s/TiB 基线)和 `advance_200`(200 MB/s/TiB 基线)。
  * `metered_size` - 文件系统已使用量，为上一小时最大使用量，非当前值，单位为 Byte。
  * `encrypt_type` - 文件系统是否加密：
    * `0`(默认值)：未加密。
    * `1`：使用 NAS 管理密钥加密。支持当 `file_system_type = standard` 或 `extreme` 时。
    * `2`：使用用户管理密钥加密。支持当 `file_system_type = extreme` 时。
  * `file_system_type` - 文件系统类型：
    * `standard`(默认值)：通用 NAS。
    * `extreme`：极速 NAS。
    * `cpfs`：文件存储 CPFS。
  * `capacity` - 文件系统的容量。
  * `kms_key_id` - KMS 密钥 ID。
  * `zone_id` - 可用区 ID。可用区是指在同一地域内，电力和网络互相独立的物理区域：
    * 当 `file_system_type = standard` 时为可选项。默认情况下，会根据 `protocol_type` 和 `storage_type` 配置随机选择符合条件的可用区。
    * 当 `file_system_type = extreme` 或 `file_system_type = cpfs` 时为必选项。建议文件系统和 ECS 实例属于同一可用区，以避免跨可用区延迟。