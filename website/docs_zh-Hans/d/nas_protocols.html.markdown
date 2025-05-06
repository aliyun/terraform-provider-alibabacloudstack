---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_protocols"
sidebar_current: "docs-alibabacloudstack-datasource-nas-protocols"
description: |-
    查询当前凭证权限可以使用的NAS文件系统的协议类型
---

# alibabacloudstack_nas_protocols

根据指定过滤条件列出当前凭证权限可以使用的NAS文件系统的协议类型。


## 示例用法

```terraform
data "alibabacloudstack_nas_protocols" "default" {
  type        = "Performance"
  zone_id     = "cn-beijing-e"
  output_file = "protocols.txt"
}

output "nas_protocols_protocol" {
  value = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
```

## 参数说明

以下参数可用于配置数据源：

* `type` - (必填) 文件系统类型。有效值为：`Performance` 和 `Capacity`。
* `zone_id` - (可选) 字符串，用于按可用区 ID 过滤结果。
* `protocols` - (可选) 受支持的协议类型的列表，用于进一步筛选结果。

## 属性说明

除上述参数外，还导出以下属性：

* `protocols` - 受支持的协议类型的列表。该列表包含当前账户在指定条件下可用的所有协议类型。