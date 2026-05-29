---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_dir_quota"
sidebar_current: "docs-Alibabacloudstack-resource-nas-dir-quota"
description: |-
  管理NAS目录配额
---

# alibabacloudstack_nas_dir_quota

管理阿里云NAS文件系统的目录配额，可以为指定目录设置用户或用户组的容量和文件数量限制。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testaccnasdirquota60008"
}


data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_file_system" "default" {
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  encrypt_type  = "0"
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
}




resource "alibabacloudstack_nas_dir_quota" "default" {
  file_system_id = alibabacloudstack_nas_file_system.default.id
  path           = "/"
  quotas {
    file_count_limit = 10000
    quota_type       = "Enforcement"
    user_type        = "Uid"
    user_id          = "500"
    size_limit       = 100
  }

}
```

## 参数说明

支持以下参数：

* `file_system_id` - (必填, 变更时重建) 文件系统的ID，用于指定要设置配额的NAS文件系统。
* `path` - (必填, 变更时重建) 目录在文件系统中的绝对路径，例如"/"或"/data/sub1"。
* `quotas` - (可选) 配额配置集合，至少需要配置1个配额项。每个配额项包含以下参数：
  * `quota_type` - (必填) 配额类型，可选值为：
    * `Enforcement`：限制型配额，当使用量超过限制后，会导致创建文件或目录、追加写入等操作失败。
    * `Accounting`：统计型配额，只统计使用量，不限制操作。
  * `user_type` - (必填) 用户类型，可选值为：
    * `Uid`：用户ID。
    * `Gid`：用户所属组ID。
    * `AllUsers`：所有用户。
  * `user_id` - (可选) 要限制的Uid或Gid，当`user_type`为`Uid`或`Gid`时需要指定。
  * `size_limit` - (可选) 目录下用户的文件总容量限制，单位为GB。当`quota_type`为`Enforcement`时，`size_limit`和`file_count_limit`至少需要填写其中一项。
  * `file_count_limit` - (可选) 目录下用户的文件数目限制，包括文件、目录和特殊文件。当`quota_type`为`Enforcement`时，`size_limit`和`file_count_limit`至少需要填写其中一项。

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为`{FileSystemId:Path}`。
* `status` - 目录的统计状态。可能值：
  * `Initializing`：初始化中。
  * `Normal`：正常状态。