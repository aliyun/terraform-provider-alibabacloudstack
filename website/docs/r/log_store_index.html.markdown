---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_store_index"
sidebar_current: "docs-alibabacloudstack-resource-log-store-index"
description: |-
  Provides a Alibabacloudstack log store index resource.
---

# alibabacloudstack_log_store_index

Log Service provides the LogSearch/Analytics function to query and analyze large amounts of logs in real time.
You can use this function by enabling the index and field statistics. [Refer to details](https://www.alibabacloud.com/help/doc-detail/43772.htm)

## Example Usage

Basic Usage
To invoke this resource, you need to set the provider parameter: sls_openapi_endpoint
```
provider "alibabacloudstack" {
  sls_openapi_endpoint = "var.sls_openapi_endpoint"
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


## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `logstore` - (Required, ForceNew) The log store name to the query index belongs.
* `full_text` - The configuration of full text index. Valid item as follows:

    * `case_sensitive` - (Optional) Whether the case sensitive. Default to false.
    * `include_chinese` - (Optional) Whether includes the chinese. Default to false.
    * `token` - (Optional) The string of several split words, like "\r", "#"

* `field_search` - List configurations of field search index. Valid item as follows:

    * `name` - (Required) The field name, which is unique in the same log store.
    * `type` - (Optional) The type of one field. Valid values: ["long", "text", "double", "json"]. Default to "long".
    * `alias` - (Optional) The alias of one field
    * `case_sensitive` - (Optional) Whether the case sensitive for the field. Default to false. It is valid when "type" is "text" or "json".
    * `include_chinese` - (Optional) Whether includes the chinese for the field. Default to false. It is valid when "type" is "text" or "json".
    * `token` - (Optional) The string of several split words, like "\r", "#". It is valid when "type" is "text" or "json".
    * `enable_analytics` - (Optional) Whether to enable field analytics. Default to true.
    * `json_keys` - (Optional, Available in 1.66.0+) Use nested index when type is json
        * `name` - (Required) When using the json_keys field, this field is required.
        * `type` - (Optional) The type of one field. Valid values: ["long", "text", "double"]. Default to "long"
        * `alias` - (Optional) The alias of one field.
        * `doc_value` - (Optional) Whether to enable statistics. default to true.
    
    
    * `enable_analytics` - (Optional) Whether to enable field analytics. Default to true.

-> **Note:** At least one of the "full_text" and "field_search" should be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log store index. It formats of `<project>:<logstore>`.
* `project` - (ForceNew, Required) The project name to the log store belongs.
* `logstore` - (ForceNew, Required) The log store name to the query index belongs.

## Import

Log store index can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_log_store_index.example tf-log:tf-log-store
```