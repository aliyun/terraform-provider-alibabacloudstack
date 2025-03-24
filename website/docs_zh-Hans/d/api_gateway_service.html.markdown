---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_service"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-service"
description: |-
    提供一个数据源以自动开启 API 网关服务。
---

# alibabacloudstack_api_gateway_service

使用此数据源可以自动启用 API 网关服务。如果服务已经启用，它将返回 `Opened`。

关于 API 网关的更多信息以及如何使用它，请参阅 [什么是 API 网关](https://www.alibabacloud.com/help/product/29462.htm)。



## 示例用法

```
data "alibabacloudstack_api_gateway_service" "open" {
	enable = "On"
}
```

## 参数参考

支持以下参数：

* `enable` - (可选) 将值设置为 `On` 以启用服务。如果服务已启用，则返回结果。有效值： "On" 或 "Off"。默认为 "Off"。

> **注意:** 将 `enable = "On"` 设置为开启 API 网关服务，这意味着您已阅读并同意 [API 网关服务条款](https://help.aliyun.com/document_detail/35391.html)。一旦服务开启，将无法关闭。

## 属性参考

除了上述参数外，还导出以下属性：

* `status` - 当前服务启用状态。