---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_service"
sidebar_current: "docs-alibabacloudstack-datasource-datahub-service"
description: |-
    提供一个数据源以自动开启 DataHub 服务。
---

# alibabacloudstack_datahub_service

使用此数据源可以自动开启 DataHub 服务。如果服务已经开启，它将返回已开启状态。

有关 DataHub 的信息以及如何使用它，请参阅 [什么是 DataHub](https://help.aliyun.com/product/53345.html)。



## 示例用法

```terraform
data "alibabacloudstack_datahub_service" "open" {
  enable = "On"
}
```

## 参数参考

支持以下参数：

* `enable` - (可选) 将值设置为 `On` 以启用服务。如果服务已被启用，则返回结果。有效值：`On` 或 `Off`。默认为 `Off`。

-> **注意:** 设置 `enable = "On"` 将会开启 DataHub 服务，这意味着您已阅读并同意 [DataHub 服务条款](https://help.aliyun.com/document_detail/158927.html)。一旦服务被开启，将无法关闭。

## 属性参考

除了上述参数列表之外，还导出以下属性：

* `status` - 当前服务的启用状态。