---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_kafka_group"
sidebar_current: "docs-Alibabacloudstack-datahub-kafka_group"
description: |-
  管理DataHub Kafka组资源
---

# alibabacloudstack_datahub_kafka_group

管理DataHub Kafka组资源，用于创建和管理DataHub项目中的Kafka消费组。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf_testacc_datahub_group21719"
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
  topic_list = [
    "${alibabacloudstack_datahub_topic.default.0.name}"
  ]
  project_name = alibabacloudstack_datahub_project.default.name
  comment      = "test group"
  group_name   = "tf_testacc_datahub_group21719"
}
```

## 参数说明

支持以下参数：

* `comment` - (必填) Kafka组的描述信息。创建后不可修改。
* `group_name` - (必填, 变更时重建) Kafka组的名称。必须符合DataHub命名规范，长度限制为1-128字符。
* `project_name` - (必填, 变更时重建) DataHub项目的名称。必须已存在且符合DataHub项目命名规范。
* `topic_list` - (可选) 绑定到该Kafka组的Topic列表。至少需要指定一个Topic，Topic名称需符合DataHub Topic命名规范。更新时会调用`UpdateTopicsForKafkaGroup` API进行绑定关系变更。

## 属性说明

以下属性导出：

* `id` - 资源ID，格式为`{ProjectName}:{GroupName}`。
* `create_time` - Kafka组的创建时间，Unix时间戳（毫秒）。