---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_topic"
sidebar_current: "docs-Alibabacloudstack-datahub-topic"
description: |- 
  编排datahub主题
---

# alibabacloudstack_datahub_topic

使用Provider配置的凭证在指定的资源集下编排datahub主题。

## 示例用法

### 基础用法


```hcl
variable "project_name" {
    default = "tf_testacc_datahub_project"
}

variable "blob_topic_name" {
    default = "tf_testacc_datahub_blob_topic"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test project"
    name = var.project_name
}

resource "alibabacloudstack_datahub_topic" "blob_example" {
  name         = var.blob_topic_name
  project_name = alibabacloudstack_datahub_project.default.name
  record_type  = "BLOB"
  shard_count  = 3
  life_cycle   = 7
  comment      = "created by terraform"
}
```

- **TUPLE 主题**

```hcl
variable "tuple_topic_name" {
    default = "tf_testacc_datahub_tuple_topic"
}

resource "alibabacloudstack_datahub_topic" "tuple_example" {
  name         = var.tuple_topic_name
  project_name = alibabacloudstack_datahub_project.default.name
  record_type  = "TUPLE"
  record_schema = {
    bigint_field    = "BIGINT"
    timestamp_field = "TIMESTAMP"
    string_field    = "STRING"
    double_field    = "DOUBLE"
    boolean_field   = "BOOLEAN"
  }
  shard_count = 3
  life_cycle  = 7
  comment     = "created by terraform"
}
```

## 参数说明

支持以下参数：

* `project_name` - (必填，变更时重建) 该主题所属的 DataHub 项目的名称。它是不区分大小写的，长度不能超过 128 个字符。
* `name` - (必填，变更时重建) DataHub 主题的名称。其长度限制为 1-128 个字符，仅允许字母、数字和下划线 (`_`)。它是不区分大小写的。
* `shard_count` - (选填，变更时重建) 此主题包含的分片数量。允许的值范围是 [1, 10]。默认值是 1。
* `life_cycle` - (选填) 该主题的数据保留期（以天为单位）。允许的值范围是 [1, 7]。默认值是 3。
* `record_type` - (选填，变更时重建) 主题的类型。它必须是 `BLOB` 或 `TUPLE` 之一。对于 `BLOB` 主题，数据将以二进制形式组织并由 BASE64 编码。对于 `TUPLE` 主题，数据具有固定的模式。默认值是带有模式 `{STRING}` 的 `TUPLE`。
* `record_schema` - (选填，变更时重建) 此主题的模式，仅适用于 `TUPLE` 主题。支持的数据类型（不区分大小写）有：
  - `BIGINT`
  - `STRING`
  - `BOOLEAN`
  - `DOUBLE`
  - `TIMESTAMP`
* `comment` - (选填) DataHub 主题的注释。它的长度不能超过 255 个字符。

**注意：** 目前，`life_cycle` 字段无法修改，并将在未来得到支持。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - DataHub 主题的 ID。它由项目名称和主题名称组成，格式为 `<project_name>:<name>`。
* `create_time` - DataHub 主题的创建时间。这是一个人类可读的字符串，而不是 64 位 UTC 时间戳。
* `last_modify_time` - DataHub 主题的最后修改时间。最初，它与 `create_time` 相同。它也是一个人类可读的字符串，而不是 64 位 UTC 时间戳。