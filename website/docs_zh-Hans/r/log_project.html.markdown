---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_project"
sidebar_current: "docs-alibabacloudstack-resource-log-project"
description: |-
  编排日志告警的项目
---

# alibabacloudstack_log_project

使用Provider配置的凭证在指定的资源集编排日志告警的项目。
项目是日志服务中的资源管理单元，用于隔离和控制资源。
您可以使用项目来管理应用程序的所有日志及其相关的日志源。

## 示例用法

### 基础用法
要调用此资源，您需要在provider参数中设置sls的endpoint地址
```
provider "alibabacloudstack" {
  endpoints {
    sls_endpoint = "var.sls_openapi_endpoint"
  }
  ...
}

resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}
```


## 参数参考

支持以下参数：

* `name` - (必填，变更时重建) 日志项目的名称。在一个 Alibabacloudstack 账户中唯一。
* `description` - (可选) 日志项目的描述。


## 属性参考

导出以下属性：

* `id` - 日志项目的 ID。它与名称相同。
* `name` - 日志项目名称。
* `description` - 日志项目描述。