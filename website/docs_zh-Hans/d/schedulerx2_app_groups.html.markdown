---
subcategory: "分布式任务调度 SchedulerX"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_schedulerx2_app_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-schedulerx2-app-groups"
description: |-
  获取Schedulerx2应用组列表。
---

# alibabacloudstack\_schedulerx2\_app\_groups

该数据源提供当前阿里云用户下的Schedulerx2应用组信息。

## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "example" {
  namespace = "example_namespace"
}

output "first_app_group_id" {
  value = data.alibabacloudstack_schedulerx2_app_groups.example.groups.0.id
}
```

### 通过名称正则表达式过滤

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "filtered" {
  name_regex = "^myapp.*"
}

output "filtered_app_groups" {
  value = data.alibabacloudstack_schedulerx2_app_groups.filtered.groups[*].app_name
}
```

### 通过特定ID过滤

```terraform
data "alibabacloudstack_schedulerx2_app_groups" "by_ids" {
  ids = ["123", "456"]
}

output "specific_app_groups" {
  value = data.alibabacloudstack_schedulerx2_app_groups.by_ids.groups[*].app_name
}
```

## 参数参考

以下参数可用于过滤和查询应用组：

* `namespace` - (可选) 应用组的命名空间。默认为"system_namespace"。
* `department` - (可选) 部门信息。
* `ids` - (可选) 应用组ID列表。
* `name_regex` - (可选) 用于按应用组名称过滤结果的正则表达式。
* `app_name` - (可选) 应用程序的名称。

## 属性参考

以下属性将会被导出：

* `ids` - 应用组ID列表。
* `groups` - 应用组列表。每个元素包含以下属性：
  * `id` - 应用组ID。
  * `app_group_id` - 应用组ID。
  * `app_name` - 应用程序名称。
  * `app_key` - 应用密钥。
  * `description` - 应用组描述。
  * `group_id` - 组ID。
  * `max_jobs` - 最大作业数。
  * `max_concurrency` - 最大并发数。
  * `xattrs` - 扩展属性。
  * `version` - 应用组版本。
  * `monitor_config` - 监控配置。
  * `metrics_threshold_json` - JSON格式的指标阈值。
  * `parent_group_id` - 父组ID。
  * `creator` - 应用组创建者。
  * `updater` - 应用组最后更新者。
  * `unique_id` - 应用组唯一ID。
  * `global_max_jobs` - 全局最大作业数。
  * `accept_lang` - 接受的语言。
  * `enable_log` - 是否启用日志记录。
  * `log_config_id` - 日志配置ID。
  * `contact_group_id` - 联系人组ID。
  * `cur_jobs` - 当前作业数。
  * `leader` - 领导者信息。
  * `alive_workers` - 存活的工作节点数。
  * `auto_scale` - 是否启用自动扩展。
  * `app_type` - 应用类型。
  * `alarm_json` - JSON格式的告警配置。
```
```