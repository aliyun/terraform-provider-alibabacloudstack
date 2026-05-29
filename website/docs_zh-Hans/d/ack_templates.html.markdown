---
subcategory: "Kubernetes容器监控"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_templates"
sidebar_current: "docs-Alibabacloudstack-datasource-ack-templates"
description: |-
  提供用户可用的ACK模板列表。
---

# alibabacloudstack_ack_templates

本数据源根据指定的过滤条件提供ACK模板列表。

## 示例用法

```hcl
# 声明数据源
data "alibabacloudstack_ack_templates" "example" {
  name_regex = "my-template"
}
```
## 参数参考
## 支持以下参数：

* `ids` - (可选) 用于过滤结果的模板ID列表。
* `name_regex` - (可选) 用于按模板名称过滤结果的正则表达式。
* `description_regex` - (可选) 用于按模板描述过滤结果的正则表达式。
* `template_type` - (可选) 按模板类型过滤结果。
## 属性参考
## 导出以下属性：

* `ids` - 模板ID列表。
* `templates` - 模板列表。每个元素包含以下属性：
* `template` - 模板内容。
* `name` - 模板名称。
* `description` - 模板描述。
* `template_type` - 模板类型。
* `template_id` - 模板ID。
* `template_with_hist_id` - 带历史记录的模板ID。
* `template_hash_code_version` - 模板的哈希码版本。
* `created` - 模板创建时间。
* `acl` - 模板的ACL。
* `version` - 模板版本。
* `tags` - 模板标签。
* `ali_uid` - 阿里云UID。
* `updated` - 模板最后更新时间。