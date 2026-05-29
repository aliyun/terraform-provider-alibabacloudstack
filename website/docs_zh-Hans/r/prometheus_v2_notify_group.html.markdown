---
subcategory: "Prometheus 监控服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_notify_group"
sidebar_current: "docs-Alibabacloudstack-prometheus-prometheus_v2_notify_group"
description: |-
  管理Prometheus v2的通知组
---

# alibabacloudstack_prometheus_v2_notify_group

管理Prometheus v2的通知组资源，用于配置告警通知渠道（如钉钉、企业微信、联系人列表或Webhook）。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc-notifygroup-15589"
}


resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  webhook_header_params {
    key   = "aaaa"
    value = "1111"
  }

  name        = var.name
  type        = "WEBHOOK"
  description = var.name
  webhook_url = "xxxxxxxxxxxxx"
}
```

## 参数说明

支持以下参数：

* `name` - (必填, 变更时重建) 通知组的名称。长度限制和格式遵循阿里云规范。
* `type` - (必填, 变更时重建) 通知组的类型。取值范围：`DINGDING`（钉钉）、`WECHAT_ROBOT`（企业微信机器人）、`CONTACT`（联系人列表）、`WEBHOOK`（自定义Webhook）。

* `contact_ids` - (可选) 联系人ID列表。当`type`为`CONTACT`时必填，需指定有效的联系人ID集合。
* `description` - (可选) 通知组的描述信息。用于说明通知组的用途或配置细节。
* `im` - (可选) IM地址（如钉钉群机器人Webhook地址或企业微信机器人Key）。当`type`为`DINGDING`或`WECHAT_ROBOT`时必填。
* `webhook_header_params` - (可选) Webhook请求头参数集合。当`type`为`WEBHOOK`时可选，每个元素包含：
  * `key` - (必填) 请求头字段名称
  * `value` - (必填) 请求头字段值
* `webhook_url` - (可选) Webhook的URL地址。当`type`为`WEBHOOK`时必填，需提供有效的HTTP/HTTPS端点。

## 属性说明

导出以下属性：

* `id` - 通知组的唯一标识ID（由API返回的数字ID转换为字符串）。
* `contact_ids` - 实际生效的联系人ID列表（仅当`type`为`CONTACT`时存在）。
* `description` - 通知组的实际描述信息。
* `im` - 实际配置的IM地址（仅当`type`为`DINGDING`或`WECHAT_ROBOT`时存在）。
* `type` - 通知组的实际类型。
* `webhook_header_params` - 实际生效的Webhook请求头参数集合（仅当`type`为`WEBHOOK`时存在）。
* `webhook_url` - 实际配置的Webhook URL地址（仅当`type`为`WEBHOOK`时存在）。