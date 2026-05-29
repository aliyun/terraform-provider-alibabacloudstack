---
subcategory: "容器镜像服务 ACR"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_attestor_lifecycle_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-cr-ee-attestor-lifecycle-rules"
description: |-
  查询阿里云容器镜像服务（ACR）企业版的保留策略规则
---

# alibabacloudstack_cr_ee_attestor_lifecycle_rules

用于查询阿里云容器镜像服务（ACR）企业版中配置的镜像保留策略规则。该数据源支持通过实例ID、命名空间正则表达式等条件过滤规则，并返回匹配的规则列表及其详细属性。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc-cree-rule-3142"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_attestor_lifecycle_rule" "default" {
  instance_id         = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  scope               = "NAMESPACE"
  retention_tag_count = "30"
  tag_regexp          = "release-v.*"
  enable_delete_tag   = "true"
  namespace_name      = alibabacloudstack_cr_ee_namespace.default.name
  recent_pull_keep    = 30
  recent_push_keep    = 20
}

data "alibabacloudstack_cr_ee_attestor_lifecycle_rules" "default" {
  instance_id = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  ids         = ["${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}"]
}

```

## 参数说明
以下参数用于过滤查询结果。参数按类型排序：必填参数 > 可选参数（同类参数按字母序排序）。

* `instance_id` (字符串) (必填)：关联生命周期规则的容器镜像服务实例ID。例如 `cri-private`。

* `enable_delete_untagged_manifest` (布尔) (可选)：根据是否启用未标记清单删除进行过滤。当设置为 `true` 时，仅返回启用该功能的规则。

* `ids` (列表) (可选)：用于过滤结果的规则ID列表。每个ID格式为 `{instance_id}:{rule_id}`，例如 `cri-private:cralr-42x00j4drgt4fjr3`。

* `namespace_regex` (字符串) (可选)：用于过滤命名空间名称的正则表达式。例如 `test.*` 可匹配 `testtf` 命名空间。

## 属性说明
以下属性为数据源导出的只读属性（`Computed`）。`id` 属性始终位于首位，其余顶级属性按字母序排序。嵌套属性在父属性下按字母序排列。

* `id` (字符串)：数据源唯一标识符，由过滤条件的哈希值生成。

* `names` (列表)：匹配的命名空间名称列表。每个元素为字符串类型，表示符合过滤条件的命名空间名称（如 `testtf`）。

* `rules` (列表)：匹配的生命周期规则列表。每个规则对象包含以下属性：
  * `auto` (布尔)：规则是否自动执行。`false` 表示手动触发（如 `MANUAL` 调度模式）。
  * `create_time` (整数)：规则创建时间（Unix 时间戳，毫秒级）。
  * `enable_delete_tag` (布尔)：是否启用标签删除功能。
  * `enable_delete_untagged_manifest` (布尔)：是否启用未标记清单删除功能。
  * `instance_id` (字符串)：关联的容器镜像服务实例ID（如 `cri-private`）。
  * `modified_time` (整数)：规则最后修改时间（Unix 时间戳，毫秒级）。
  * `namespace_name` (字符串)：规则应用的命名空间名称（如 `testtf`）。
  * `next_time` (整数)：规则下次执行时间（Unix 时间戳，毫秒级）。
  * `recent_pull_keep` (整数)：保留最近拉取记录的天数（用于基于拉取时间的保留策略）。
  * `recent_push_keep` (整数)：保留最近推送记录的天数（用于基于推送时间的保留策略）。
  * `repo_name` (字符串)：规则应用的仓库名称（当规则作用于仓库级别时有效）。
  * `retention_tag_count` (整数)：保留的标签数量阈值（如 `30` 表示保留最近30个标签）。
  * `rule_id` (字符串)：规则唯一标识符（如 `cralr-42x00j4drgt4fjr3`）。
  * `schedule_time` (字符串)：规则调度时间配置（如 `MANUAL` 表示手动触发）。
  * `tag_regexp` (字符串)：用于匹配标签的正则表达式（如 `release-v.*`）。