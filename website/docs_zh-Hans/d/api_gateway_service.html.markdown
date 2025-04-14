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

## 参数说明

支持以下参数：  

* `enable` - (可选) 将值设置为 `On` 以启用服务。如果服务已启用，则返回结果。有效值： "On" 或 "Off"。默认为 "Off"。  

> **注意:** 将 `enable = "On"` 设置为开启 API 网关服务，这意味着您已阅读并同意 [API 网关服务条款](https://help.aliyun.com/document_detail/35391.html)。一旦服务开启，将无法关闭。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `status` - 当前服务的启用状态。可能的值包括 `Opened`（已开启）和 `Closed`（已关闭）。  
* `enable` - 表示 API 网关服务是否已启用。该属性是根据服务状态计算得出的，值为 `On` 或 `Off`。 如果服务当前处于开启状态，则返回 `On`；否则返回 `Off`。  
