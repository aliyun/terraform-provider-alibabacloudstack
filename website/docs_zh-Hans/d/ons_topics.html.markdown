---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_topics"
sidebar_current: "docs-alibabacloudstack-datasource-ons-topics"
description: |-
    查询消息队列服务主题
---

# alibabacloudstack_ons_topics

根据指定过滤条件列出当前凭证权限可以访问的消息队列服务主题列表。



## 示例用法

```
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicDatasourceName"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = var.name
  remark = "Ons Instance"
}

resource "alibabacloudstack_ons_topic" "default" {
  topic = var.topic
  instance_id = alibabacloudstack_ons_instance.default.id
  message_type = 0
  remark = "dafault_ons_topic_remark"
}

data "alibabacloudstack_ons_topics" "topics_ds" {
 instance_id = alibabacloudstack_ons_topic.topic.instance_id
  output_file = "topics.txt"
}

output "first_topic_name" {
   value = data.alibabacloudstack_ons_topics.topics_ds.*
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填) 拥有这些主题的ONS实例ID。
* `name_regex` - (可选) 用于通过主题名称筛选结果的正则表达式字符串。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `names` - 主题名称列表。
* `topics` - 主题列表。每个元素包含以下属性：
  * `id` - 主题的ID。
  * `instance_id` - 主题所属的实例ID。
  * `topic` - 主题名称。
  * `owner` - 主题所有者的ID，即Apsara Stack Cloud UID。
  * `relation` - 关系ID。
  * `relation_name` - 关系名称，例如：`owner`（所有者）、`publishable`（可发布）、`subscribable`（可订阅）以及`publishable and subscribable`（可发布且可订阅）。
  * `message_type` - 消息类型，取值范围为：
    * `0`：普通消息
    * `1`：事务消息
    * `2`：定时/延时消息
  * `independent_naming` - 表示是否启用命名空间，取值为布尔值：
    * `true`：启用命名空间
    * `false`：未启用命名空间
  * `create_time` - 创建时间，格式为Unix时间戳。
  * `remark` - 主题备注信息。
  * `computed_attribute` - 计算属性的描述，通常用于内部计算或扩展功能。