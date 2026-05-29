---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_notify_group"
sidebar_current: "docs-Alibabacloudstack-resource-prometheus-v2-notify-group"
description: |-
  Manages Prometheus v2 notification groups.
---

# alibabacloudstack_prometheus_v2_notify_group

Manages Prometheus v2 notification group resources for configuring alert notification channels (such as DingTalk, WeCom, contact lists, or Webhook).

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the notification group. Length restrictions and format follow Alibaba Cloud specifications.
* `type` - (Required, Forces new resource) The type of the notification group. Valid values: `DINGDING` (DingTalk), `WECHAT_ROBOT` (WeCom Robot), `CONTACT` (Contact List), `WEBHOOK` (Custom Webhook).
* `contact_ids` - (Optional) A list of contact IDs. Required when `type` is `CONTACT`, must specify a valid set of contact IDs.
* `description` - (Optional) The description of the notification group. Used to explain the purpose or configuration details.
* `im` - (Optional) The IM address (e.g., DingTalk group robot Webhook URL or WeCom robot key). Required when `type` is `DINGDING` or `WECHAT_ROBOT`.
* `webhook_header_params` - (Optional) A set of Webhook request header parameters. Optional when `type` is `WEBHOOK`, each element contains:
  * `key` - (Required) The name of the header field.
  * `value` - (Required) The value of the header field.
* `webhook_url` - (Optional) The URL of the Webhook. Required when `type` is `WEBHOOK`, must provide a valid HTTP/HTTPS endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier ID of the notification group (converted from the numeric ID returned by the API to a string).
* `contact_ids` - The actual effective list of contact IDs (exists only when `type` is `CONTACT`).
* `description` - The actual description of the notification group.
* `im` - The actually configured IM address (exists only when `type` is `DINGDING` or `WECHAT_ROBOT`).
* `type` - The actual type of the notification group.
* `webhook_header_params` - The actual effective set of Webhook request header parameters (exists only when `type` is `WEBHOOK`).
* `webhook_url` - The actually configured Webhook URL address (exists only when `type` is `WEBHOOK`).