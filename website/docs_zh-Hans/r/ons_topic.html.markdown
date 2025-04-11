---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_topic"
sidebar_current: "docs-alibabacloudstack-resource-ons-topic"
description: |-
  编排消息队列（ONS）主题

---

# alibabacloudstack_ons_topic

使用Provider配置的凭证在指定的资源集编排消息队列（ONS）主题。


## 示例用法

### 基础用法

```
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicName"
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
```

## 参数说明

支持以下参数：

* `instance_id` - (可选) 拥有该主题的 ONS 实例的 ID。
* `topic` - (必填) 主题名称。单个实例上的两个主题不能具有相同的名称，且名称不能以 'GID' 或 'CID' 开头。长度不得超过 64 个字符。
* `message_type` - (必填) 消息类型。取值范围如下：
  * `0` - 普通消息
  * `1` - 广播消息
  * `2` - 事务消息
  * `3` - 定时/延时消息
* `remark` - (必填) 这是对主题的简要描述。长度不得超过 128。
* `perm` - (可选) 此属性用于设置主题的读写模式。取值范围如下：
  * `6` - 可读可写
  * `4` - 仅可读
  * `2` - 仅可写

## 属性说明

导出以下属性：

* `id` - ONS 主题的主题和实例 ID。格式为 `Topic:InstanceID`，其中 `Topic` 是主题名称，`InstanceID` 是实例 ID。