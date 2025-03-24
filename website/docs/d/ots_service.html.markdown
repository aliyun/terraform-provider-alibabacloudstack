---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_service"
sidebar_current: "docs-alibabacloudstack-datasource-ots-service"
description: |-
    Provides a datasource to open the Table Staore service automatically.
---

# alibabacloudstack_ots_service

Using this data source can enable Table Staore service automatically. If the service has been enabled, it will return `Opened`.

For information about Table Staore and how to use it, see [What is Table Staore](https://www.alibabacloud.com/help/product/27278.htm).


## Example Usage

```
data "alibabacloudstack_ots_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the Table Staore service that means you have read and agreed the [Table Staore Terms of Service](https://help.aliyun.com/document_detail/34908.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status.