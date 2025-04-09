---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_service"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-service"
description: |-
    Provides a datasource to open the API gateway service automatically.
---

# alibabacloudstack_api_gateway_service

Using this data source can enable API gateway service automatically. If the service has been enabled, it will return `Opened`.

For information about API Gateway and how to use it, see [What is API Gateway](https://www.alibabacloud.com/help/product/29462.htm).



## Example Usage

```
data "alibabacloudstack_api_gateway_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off". 

-> **NOTE:** Setting `enable = "On"` to open the API gateway service that means you have read and agreed the [API Gateway Terms of Service](https://help.aliyun.com/document_detail/35391.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 

* `enable` - Indicates whether the API gateway service is enabled. This attribute is computed based on the service status. 