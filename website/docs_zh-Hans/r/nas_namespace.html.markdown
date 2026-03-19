---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace"
sidebar_current: "docs-Alibabacloudstack-nas-nas_namespace"
description: |-
  创建和管理NAS统一命名空间
---

# alibabacloudstack_nas_namespace

创建和管理阿里云NAS统一命名空间资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testAccNasNamespace62833"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}



data "alibabacloudstack_nas_zones" "default" {
}



resource "alibabacloudstack_nas_namespace" "default" {
  encrypt_type  = "0"
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = "tf-testAccNasNamespace62833"
}
```

## 参数说明

支持以下参数：

* `zone_id` - (必填, 变更时重建) 可用区ID。可用区是指在同一地域内，电力和网络互相独立的物理区域。同一地域不同可用区之间的文件系统与ECS云服务器互通。建议文件系统与云服务器ECS属于同一可用区，避免跨可用区产生的时延。
* `cluster_id` - (必填, 变更时重建) 集群ID。默认值为"StandardNasCluster"。
* `description` - (必填, 变更时重建) 命名空间描述。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`和`https://`开头。可以包含数字、半角冒号（:）、下划线（_）或者短划线（-）。
* `storage_type` - (必填, 变更时重建) 存储类型。取值：`Capacity`（容量型）或`Performance`（性能型）。
* `protocol_type` - (必填, 变更时重建) 文件传输协议类型。取值：`NFS`（NFS文件传输协议）或`SMB`（SMB文件传输协议）。
* `encrypt_type` - (必填, 变更时重建) 命名空间中的文件系统是否加密。使用KMS服务托管密钥，对文件系统落盘数据进行加密存储。在读写加密数据时，无需解密。取值：`0`（不加密）或`1`（加密）。

## 属性说明

以下属性会从API响应中导出：

* `id` - 命名空间ID。
* `create_time` - 命名空间创建的时间。遵循ISO 8601标准表示，返回格式：`yyyy-MM-ddTHH:mm:ssZ`。
* `file_system_type` - 文件系统类型。默认值：`standard`，通用型NAS。
* `mount_target_count` - 命名空间文件系统挂载点数量。