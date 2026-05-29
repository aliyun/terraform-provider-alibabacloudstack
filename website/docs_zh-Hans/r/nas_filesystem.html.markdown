---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_filesystem"
sidebar_current: "docs-Alibabacloudstack-nas-filesystem"
description: |- 
  编排文件存储（NAS）文件系统
---

# alibabacloudstack_nas_file_system
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_filesystem`

使用Provider配置的凭证在指定的资源集编排文件存储（NAS）文件系统。

## 示例用法

### 基础用法：

```terraform
variable "name" {
  default = "tf-testAccAlibabacloudStackNasFileSystem97828"
}

data "alibabacloudstack_nas_protocols" "example" {
  type = "Capacity"
}

data "alibabacloudstack_nas_zones" "default" {}

resource "alibabacloudstack_nas_file_system" "default" {
  zone_id        = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  description    = "tf-testAccAlibabacloudStackNasFileSystem97828"
  protocol_type  = data.alibabacloudstack_nas_protocols.example.protocols.0
  storage_type   = "Capacity"
  file_system_type = "standard"
  encrypt_type   = "1"
}
```

高级用法，使用 `extreme` 文件系统类型：

```terraform
resource "alibabacloudstack_nas_file_system" "foo" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = "cn-hangzhou-f"
  storage_type     = "standard"
  description      = "tf-testAccNasConfig"
  capacity         = "100"
  encrypt_type     = "2"
  kms_key_id       = "your-kms-key-id"
}
```

## 参数说明

支持以下参数：

* `storage_type` - (必填，变更时重建) 存储类型。
  * 当 `file_system_type = standard` 时，可选值为：
    * `Performance`(性能型)
    * `Capacity`(容量型)
  * 当 `file_system_type = extreme` 时，可选值为：
    * `standard`(标准型)
    * `advance`(高级型)

* `protocol_type` - (必填，变更时重建) 文件传输协议类型。
  * 当 `file_system_type = standard` 时，可选值为：
    * `NFS`
    * `SMB`
  * 当 `file_system_type = extreme` 时，值为：
    * `NFS`

* `description` - (选填) 文件系统描述。限制：
  * 长度为2~256个英文或中文字符。
  * 必须以大小写字母或中文开头，不能以 `http://` 或 `https://` 开头。
  * 可以包含数字、冒号(:)、下划线(_)或短横线(-)。

* `encrypt_type` - (选填，变更时重建) **[已弃用]** 该字段在专有云上不受支持，将在版本 3.21.0 中移除。是否对文件系统进行加密。使用KMS服务托管密钥来加密和存储文件系统磁盘数据。读写加密数据时，无需解密。
  * 可选值：
    * `0`(默认)：不加密。
    * `1`：NAS托管密钥。当 `file_system_type = standard` 或 `file_system_type = extreme` 时支持NAS托管密钥。
    * `2`：用户管理密钥。只有当 `file_system_type = extreme` 时可以管理密钥。

* `file_system_type` - (选填，变更时重建) 文件系统类型。
  * 可选值：
    * `standard`(默认)：通用型NAS
    * `extreme`：极速型NAS

* `capacity` - (选填) 文件系统的容量。当 `file_system_type = extreme` 时为必填项。单位：GiB。**注意**：最小值为100。

* `zone_id` - (必填，变更时重建) 可用区ID。可用区是指在同一地域内电力和网络互相独立的物理区域。
* `cluster_id` - (必填，变更时重建) 集群ID。
* `kms_key_id` - (选填) **[已弃用]** 该字段在专有云上不受支持，将在版本 3.21.0 中移除。KMS密钥ID。当 `encrypt_type = 2` 时为必填项。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 文件系统的ID。
* `zone_id` - 可用区ID。
* `capacity` - 文件系统的容量（由 API 返回）。
* `kms_key_id` - 用于加密的KMS密钥ID（由 API 返回）。