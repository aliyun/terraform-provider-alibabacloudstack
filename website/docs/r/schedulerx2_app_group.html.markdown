---
subcategory: "SchedulerX2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_schedulerx2_app_group"
sidebar_current: "docs-Alibabacloudstack-schedulerx2-app-group"
description: |-
  Manage SchedulerX2 application groups
---

# alibabacloudstack_schedulerx2_app_group

Manages SchedulerX2 application groups for application management in distributed task scheduling systems.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-testAccschedulerx2AppGroup11340"
}


resource "alibabacloudstack_schedulerx2_app_group" "example" {
  contacts {
    username    = "test"
    user_email  = "123@123.com"
    dingding_ak = "testakkkkk"
  }

  metrics_threshold {
    load5       = 10
    heap5_usage = 100
    disk_usage  = 100
  }

  group_id        = "${var.name}.terra"
  app_name        = var.name
  description     = var.name
  max_jobs        = "30"
  max_concurrency = "10"
  monitor_config {
    send_channel = "mail,ding"
    alarm_type   = "CustomContacts"
  }

}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required) The name of the application. It must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `group_id` - (Required, Forces new resource) The unique identifier of the application group. It is a string, typically representing the unique namespace identifier of the application.
* `accept_language` - (Optional) The language to accept. Valid values: `zh` (Chinese), `en` (English). Default value: empty.
* `description` - (Optional) The description of the application. It must be 1 to 256 characters in length and cannot start with `http://` or `https://`.
* `max_concurrency` - (Optional) The maximum concurrency. It indicates the maximum number of tasks that can be executed simultaneously by the application group.
* `max_jobs` - (Optional) The maximum number of jobs. It indicates the maximum number of tasks that can run simultaneously in the application group.
* `metrics_threshold` - (Optional) The metric threshold configuration for monitoring the application's running status.
  * `disk_usage` - (Optional) The disk usage threshold. An alert is triggered when this value is exceeded. Default value: `95`.
  * `heap5_usage` - (Optional) The heap memory usage threshold. An alert is triggered when this value is exceeded. Default value: `90`.
  * `load5` - (Optional) The 5-minute system load threshold. An alert is triggered when this value is exceeded. Default value: `0`.
* `monitor_config` - (Optional) The monitoring and alert configuration.
  * `alarm_type` - (Optional) The type of alert. Default value: `CustomContacts`.
  * `send_channel` - (Optional) The alert notification channel. Valid values: `ding` (DingTalk), `mail` (email), `mail,ding` (email and DingTalk), empty string (no notification).
* `namespace` - (Optional) The namespace. Default value: `system_namespace`.
* `contacts` - (Optional) The list of alert contacts.
  * `dingding_ak` - (Optional) The DingTalk Access Key, used for DingTalk alert notifications.
  * `user_email` - (Optional) The user's email address, used for email alert notifications.
  * `username` - (Optional) The user's name, used to identify the contact.

## Attributes Reference

The following attributes are exported from the API response:

* `id` - The ID of the application group.
* `app_key` - The application key, used for application authentication.