---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_kafka_groups"
sidebar_current: "docs-alibabacloudstack-datasource-datahub-kafka-groups"
description: |-
  查询阿里云DataHub Kafka Group列表
---

# alibabacloudstack_datahub_kafka_groups

查询阿里云DataHub Kafka Group列表，用于获取指定项目下的Kafka Group信息。

## 示例用法

```hcl

variable "name" {
  default = "tf_testacc_group5898772"
}
resource "alibabacloudstack_datahub_project" "default" {
  comment = "test"
  name    = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  count        = 2
  name         = "${var.name}_${count.index}"
  comment      = "test"
  record_type  = "BLOB"
  project_name = alibabacloudstack_datahub_project.default.name
}

resource "alibabacloudstack_datahub_kafka_group" "default" {
  project_name = alibabacloudstack_datahub_project.default.name
  comment      = "test group"
  group_name   = var.name
  topic_list = [
    "${alibabacloudstack_datahub_topic.default.0.name}",
    "${alibabacloudstack_datahub_topic.default.1.name}"
  ]
}


data "alibabacloudstack_datahub_kafka_groups" "default" {
  project_name = alibabacloudstack_datahub_project.default.name
}
```

## 参数说明

以下参数支持过滤查询结果：

- **project_name** (字符串, 必填)：DataHub项目名称，用于指定要查询的项目。

- **ids** (列表, 可选)：Kafka Group ID列表，格式为"ProjectName:GroupName"，用于过滤特定的Kafka Group。

- **name_regex** (字符串, 可选)：Kafka Group名称的正则表达式，用于过滤符合特定模式的Kafka Group。

## 属性说明

以下属性被导出：

- **id** (字符串)：数据源ID，基于匹配的Kafka Group ID列表计算的哈希值。

- **kafka_groups** (列表)：匹配的Kafka Group列表，每个元素包含以下属性：
  - **comment** (字符串)：Kafka Group的描述信息。
  - **create_time** (整数)：Kafka Group的创建时间，以Unix时间戳（毫秒）表示。
  - **creator** (字符串)：Kafka Group的创建者ID。
  - **group_name** (字符串)：Kafka Group的名称。
  - **last_modify_time** (整数)：Kafka Group的最后修改时间，以Unix时间戳（毫秒）表示。
  - **topic_list** (列表)：关联的Topic名称列表。