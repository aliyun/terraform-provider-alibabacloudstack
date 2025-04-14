---
subcategory: "日志服务 (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_store_index"
sidebar_current: "docs-alibabacloudstack-resource-log-store-index"
description: |-
  编排日志告警的日志库索引
---

# alibabacloudstack_log_store_index

使用Provider配置的凭证在指定的资源集编排日志告警的日志库索引。
日志服务提供了 LogSearch/Analytics 功能，可以实时查询和分析大量日志。
您可以通过启用索引和字段统计来使用此功能。[参考详细信息](https://www.alibabacloud.com/help/doc-detail/43772.htm)

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

resource "alibabacloudstack_log_store" "example" {
  project = alibabacloudstack_log_project.example.name
  name    = "tf-log-store"
  description = "created by terraform"
}

resource "alibabacloudstack_log_store_index" "example" {
  project  = alibabacloudstack_log_project.example.name
  logstore = alibabacloudstack_log_store.example.name
  full_text {
    case_sensitive = true
    token          = " #$%^*\r\n	"
  }
  field_search {
    name             = "terraform"
    enable_analytics = true
  }
}
```


## 参数说明

支持以下参数：

* `project` - (必填，变更时重建) 日志库所属的项目名称。
* `logstore` - (必填，变更时重建) 查询索引所属的日志库名称。
* `full_text` - 全文索引的配置。有效项如下：
    * `case_sensitive` - (可选) 是否区分大小写。默认为 false。
    * `include_chinese` - (可选) 是否包含中文。默认为 false。
    * `token` - (可选) 多个分隔词的字符串，例如 "\r", "#"

* `field_search` - 字段搜索索引的列表配置。有效项如下：
    * `name` - (必填) 字段名称，在同一个日志库中是唯一的。
    * `type` - (可选) 字段类型。有效值：["long", "text", "double", "json"]。默认为 "long"。
    * `alias` - (可选) 字段的别名
    * `case_sensitive` - (可选) 字段是否区分大小写。默认为 false。当 "type" 为 "text" 或 "json" 时有效。
    * `include_chinese` - (可选) 字段是否包含中文。默认为 false。当 "type" 为 "text" 或 "json" 时有效。
    * `token` - (可选) 多个分隔词的字符串，例如 "\r", "#"。当 "type" 为 "text" 或 "json" 时有效。
    * `enable_analytics` - (可选) 是否启用字段分析。默认为 true。
    * `json_keys` - (可选，1.66.0+版本可用) 当类型为 json 时使用嵌套索引
        * `name` - (必填) 使用 json_keys 字段时，此字段为必填项。
        * `type` - (可选) 字段类型。有效值：["long", "text", "double"]。默认为 "long"。
        * `alias` - (可选) 字段的别名。
        * `doc_value` - (可选) 是否启用统计。默认为 true。

-> **注意:** 至少需要指定 "full_text" 和 "field_search" 中的一个。

## 属性说明

导出以下属性：

* `id` - 日志库索引的 ID。格式为 `<project>:<logstore>`。
* `project` - (强制新，必填) 日志库所属的项目名称。
* `logstore` - (强制新，必填) 查询索引所属的日志库名称。
* `full_text` - 全文索引的配置。
* `field_search` - 字段搜索索引的列表配置。

## 导入

日志库索引可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_log_store_index.example tf-log:tf-log-store
```