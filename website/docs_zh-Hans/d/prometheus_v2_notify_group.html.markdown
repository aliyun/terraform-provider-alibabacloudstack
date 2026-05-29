---
subcategory: "Prometheus 监控服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_notify_group"
sidebar_current: "docs-Alibabacloudstack-datasource-prometheus-v2-notify-group"
description: |-
  查询阿里云Prometheus V2通知组
---

# alibabacloudstack_prometheus_v2_notify_group

查询阿里云Prometheus V2通知组。

## 示例用法

```hcl

variable "name" {
  default = "tfacc-notifygroup-73927"
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  name        = var.name
  type        = "WEBHOOK"
  description = var.name
  webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=56b42bc6e7cad53bab514a583847db73c68fa1804b0e72af7167954b66f7aea8"
  webhook_header_params {
    key   = "aaaa"
    value = "1111"
  }
}

data "alibabacloudstack_prometheus_v2_notify_groups" "default" {
  name_regex = alibabacloudstack_prometheus_v2_notify_group.default.name
}

```

## 参数说明
以下参数用于过滤查询结果：

* `name_regex` (可选)：用于通过通知组名称过滤结果的正则表达式。
* `ids` (可选)：用于通过通知组ID列表过滤结果。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，由过滤结果中的通知组ID哈希生成。
* `contact_ids` (集合)：联系人ID列表，表示该通知组关联的联系人。
* `description` (字符串)：通知组的描述信息。
* `im` (字符串)：即时通讯配置信息，如钉钉、企业微信等。
* `name` (字符串)：通知组的名称。
* `type` (字符串)：通知组的类型，如WEBHOOK。
* `webhook_header_params` (集合)：Webhook请求头参数，包含以下属性：
  * `key` (字符串)：请求头参数的键名。
  * `value` (字符串)：请求头参数的值。
* `webhook_url` (字符串)：Webhook的URL地址，用于接收告警通知。