---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_service"
sidebar_current: "docs-alibabacloudstack-datasource-datahub-service"
description: |-
    Provides a datasource to open the DataHub service automatically.
---

# alibabacloudstack_datahub_service

Using this data source can open DataHub service automatically. If the service has been opened, it will return opened.

For information about DataHub and how to use it, see [What is DataHub](https://help.aliyun.com/product/53345.html).



## Example Usage

```terraform
data "alibabacloudstack_datahub_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the DataHub service that means you have read and agreed the [DataHub Terms of Service](https://help.aliyun.com/document_detail/158927.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status.