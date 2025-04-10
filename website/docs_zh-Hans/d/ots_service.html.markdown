---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_service"
sidebar_current: "docs-alibabacloudstack-datasource-ots-service"
description: |-
    自动开启表格存储（OTS）服务。
---

# alibabacloudstack_ots_service

使用此数据源可以自动启表格存储（OTS）服务。如果服务已经启用，它将返回 `Opened`。

有关 Table Staore 的信息以及如何使用它，请参阅 [什么是 Table Staore](https://www.alibabacloud.com/help/product/27278.htm)。


## 示例用法

```
data "alibabacloudstack_ots_service" "open" {
	enable = "On"
}
```

## 参数说明

支持以下参数：

* `enable` - (可选) 将值设置为 `On` 以启用该服务。如果服务已启用，则返回结果。有效值："On" 或 "Off"。默认为 "Off"。

-> **注意：** 将 `enable = "On"` 设置为开启 Table Staore 服务，这意味着您已阅读并同意 [Table Staore 服务条款](https://help.aliyun.com/document_detail/34908.html)。一旦服务开启，便无法关闭。

## 属性说明

除了上述参数外，还导出以下属性：

* `status` - 当前服务的启用状态。有效值为 "Opened" 或 "Closed"，分别表示服务已启用或未启用。