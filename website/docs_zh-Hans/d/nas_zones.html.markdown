---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_zones"
sidebar_current: "docs-alibabacloudstack-datasource-nas-zones"
description: |-
  查询NAS文件系统的区域类型
---

# alibabacloudstack_nas_zones

根据指定过滤条件列出当前凭证权限可以使用的NAS文件系统的区域类型。


## 示例用法

```terraform
data "alibabacloudstack_nas_zones" "default" {}

output "alibabacloudstack_nas_zones_id" {
  value = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
}
```

## 参数参考

支持以下参数：

* `file_system_type` - (可选，强制新，v1.152.0+可用)文件系统的类型。有效值：`standard`、`extreme`、`cpfs`。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `zones` - 可用区域信息集合列表。
    * `zone_id` - 字符串，按区域 ID 筛选结果。
    * `instance_types` - 实例类型信息集合列表
        * `storage_type` - NAS 区域的存储类型。有效值：
          * `standard` - 当 FileSystemType 是 standard 时。有效值：`Performance` 和 `Capacity`。
          * `extreme` - 当 FileSystemType 是 extreme 时。有效值：`Standard` 和 `Advance`。
          * `cpfs` - 当 FileSystemType 是 cpfs 时。有效值：`advance_100` 和 `advance_200`。
        * `protocol_type` - 文件传输协议类型。有效值：
          * `standard` - 当 FileSystemType 是 standard 时。有效值：`NFS` 和 `SMB`。
          * `extreme` - 当 FileSystemType 是 extreme 时。有效值：`NFS`。
          * `cpfs` - 当 FileSystemType 是 cpfs 时。有效值：`cpfs`。