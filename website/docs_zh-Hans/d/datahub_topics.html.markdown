---
subcategory: "数据总线 DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_topics"
sidebar_current: "docs-Alibabacloudstack-datasource-datahub-topics"
description: |-
  提供DataHub主题列表。
---

# alibabacloudstack_datahub_topics

该数据源用于获取专有云中可用的DataHub主题列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_datahub_topics" "example" {
  project_name = "my-project"
  name_regex   = "^test-.*"
}

output "datahub_topics" {
  value = data.alibabacloudstack_datahub_topics.example.topics
}
```

## 参数说明

以下参数支持配置：

* `project_name` - (必选, 变更后重建) DataHub项目的名称。
* `name_regex` - (可选, 变更后重建) 用于按主题名称过滤的正则表达式字符串。
* `names` - (可选, Computed) 用于过滤结果的主题名称列表。

## 属性参考

以下属性会被导出：

* `ids` - 主题ID列表。
* `names` - 主题名称列表。
* `topics` - DataHub主题列表。每个元素包含以下属性：
  * `id` - 主题的ID，格式为 `{project_name}:{topic_name}`。
  * `name` - 主题的名称。
  * `shard_count` - 分片数量。
  * `life_cycle` - 主题的生命周期（天）。
  * `comment` - 主题的备注信息。
  * `record_type` - 主题的记录类型。
  * `create_time` - 主题的创建时间。
