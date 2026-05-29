---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_zones"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-zones"
description: |-
  查询 NAS 可用区信息
---

# alibabacloudstack_nas_zones

根据指定过滤条件，查询当前地域下 NAS 文件系统的可用区及支持的类型。


## 示例用法

```terraform
data "alibabacloudstack_nas_zones" "default" {}

output "alibabacloudstack_nas_zones_id" {
  value = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
}
```

## 参数说明

支持以下参数：

* `file_system_type` -（可选）文件系统类型。有效值：`standard`、`extreme`、`cpfs`。默认值：`standard`。
* `zone_id` -（可选）按可用区 ID 过滤结果。
* `protocol` -（可选）按协议类型过滤结果。有效值：`NFS`、`SMB`、`cpfs`。
* `output_file` -（可选，已弃用）此字段已弃用，将在 3.19.0 版本中移除。如需将内容写入文件，请使用 `local_file` provider。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `zones` - 可用区信息集合列表。
    * `zone_id` - 可用区 ID。
    * `protocols` - 该可用区支持的协议类型列表。
    * `clusters` - 可用区内的集群信息列表。
        * `cluster_id` - 集群 ID。
        * `cluster_type` - 集群类型。
        * `cluster_version` - 集群版本。
        * `instance_types` - 集群内的实例类型信息列表。
            * `storage_type` - 存储类型。具体取值与 `file_system_type` 相关：
                * 当 `file_system_type` 为 `standard` 时：`Performance`、`Capacity`。
                * 当 `file_system_type` 为 `extreme` 时：`Standard`、`Advance`。
                * 当 `file_system_type` 为 `cpfs` 时：`advance_100`、`advance_200`。
            * `protocol_type` - 文件传输协议类型。具体取值与 `file_system_type` 相关：
                * 当 `file_system_type` 为 `standard` 时：`NFS`、`SMB`。
                * 当 `file_system_type` 为 `extreme` 时：`NFS`。
                * 当 `file_system_type` 为 `cpfs` 时：`cpfs`。
