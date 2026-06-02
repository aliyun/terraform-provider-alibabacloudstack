---
subcategory: "SchedulerX"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_schedulerx2_app_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-schedulerx2-app-groups"
description: |-
  Provides a list of Schedulerx2 App Groups to the user.
---

# alibabacloudstack\_schedulerx2\_app\_groups

This data source provides the Schedulerx2 App Groups of the current Alibaba Cloud user.

## Example Usage

### Basic Usage

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "example" {
  namespace = "example_namespace"
}

output "first_app_group_id" {
  value = data.alibabacloudstack_schedulerx2_app_groups.example.groups.0.id
}
```

### Filter by name regex

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "filtered" {
  name_regex = "^myapp.*"
}

output "filtered_app_groups" {
  value = data.alibabacloudstack_schedulerx2_app_groups.filtered.groups[*].app_name
}
```

### Filter by specific IDs

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "by_ids" {
  ids = ["123", "456"]
}

output "specific_app_groups" {
  value = data.alibabacloudstack_schedulerx2_app_groups.by_ids.groups[*].app_name
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Optional) The namespace of the application group. Default to "system_namespace".
* `department` - (Optional) The department information.
* `ids` - (Optional) A list of App Group IDs.
* `name_regex` - (Optional) A regex string to filter results by App Group name.
* `app_name` - (Optional) The name of the application.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of App Group IDs.
* `groups` - A list of App Groups. Each element contains the following attributes:
  * `id` - The ID of the App Group.
  * `app_group_id` - The application group ID.
  * `app_name` - The name of the application.
  * `app_key` - The application key.
  * `description` - The description of the application group.
  * `group_id` - The group ID.
  * `max_jobs` - The maximum number of jobs.
  * `max_concurrency` - The maximum concurrency.
  * `xattrs` - Extended attributes.
  * `version` - The version of the application group.
  * `monitor_config` - Monitor configuration.
  * `metrics_threshold_json` - Metrics threshold in JSON format.
  * `parent_group_id` - The parent group ID.
  * `creator` - The creator of the application group.
  * `updater` - The last updater of the application group.
  * `unique_id` - The unique ID of the application group.
  * `global_max_jobs` - The global maximum number of jobs.
  * `accept_lang` - Accepted language.
  * `enable_log` - Whether logging is enabled.
  * `log_config_id` - The log configuration ID.
  * `contact_group_id` - The contact group ID.
  * `cur_jobs` - Current number of jobs.
  * `leader` - The leader information.
  * `alive_workers` - Number of alive workers.
  * `auto_scale` - Whether auto scaling is enabled.
  * `app_type` - The application type.
  * `alarm_json` - Alarm configuration in JSON format.
