---
subcategory: "容器镜像服务 ACR"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_attestor_lifecycle_rule"
sidebar_current: "docs-Alibabacloudstack-cr-cr_ee_attestor_lifecycle_rule"
description: |-
  管理ACR企业版的镜像生命周期规则
---

# alibabacloudstack_cr_ee_attestor_lifecycle_rule

管理ACR企业版的镜像生命周期规则，用于自动清理过期的镜像标签和清单。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc-cree-rule-3358413"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_namespace" "default2" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = "${var.name}2"
  auto_create        = true
  default_visibility = "PRIVATE"
}



resource "alibabacloudstack_cr_ee_attestor_lifecycle_rule" "default" {
  retention_tag_count = "30"
  tag_regexp          = "release-v.*"
  enable_delete_tag   = "true"
  namespace_name      = alibabacloudstack_cr_ee_namespace.default.name
  recent_pull_keep    = 30
  recent_push_keep    = 20
  scope               = "NAMESPACE"
  instance_id         = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
}
```

## 参数说明

支持以下参数：

* `retention_tag_count` - (必填) 保留的标签数量。表示保留最近推送的标签数量，必须大于0。
* `scope` - (必填) 规则范围。有效值：`REPO`（仓库级别）、`NAMESPACE`（命名空间级别）。当设置为`REPO`时，必须指定`repo_name`；当设置为`NAMESPACE`时，必须指定`namespace_name`。
* `instance_id` - (必填, 变更时重建) 容器镜像服务企业版实例ID。实例ID变更会导致资源重建。
* `enable_delete_tag` - (可选) 是否启用删除标签。默认为`false`，设置为`true`时将删除匹配规则的标签。
* `namespace_name` - (可选) 命名空间名称。当`scope`为`NAMESPACE`时必须指定。
* `recent_pull_keep` - (可选) 保留最近拉取的镜像数量。表示保留最近拉取的多少个镜像，0表示不保留。
* `recent_push_keep` - (可选) 保留最近推送的镜像数量。表示保留最近推送的多少个镜像，0表示不保留。
* `repo_name` - (可选) 仓库名称。当`scope`为`REPO`时必须指定。
* `tag_regexp` - (可选) 标签正则表达式。用于匹配要保留的标签（例如`release-v.*`），设置后规则将自动执行。

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为`<instance_id>:<rule_id>`。
* `auto` - 是否自动执行规则。当设置`tag_regexp`时自动为`true`。
* `create_time` - 规则创建时间，Unix时间戳（毫秒）。
* `enable_delete_untagged_manifest` - 是否启用删除未标记的清单。
* `modified_time` - 规则修改时间，Unix时间戳（毫秒）。
* `rule_id` - 保留策略规则ID。
* `schedule` - 调度方式。固定为`MANUAL`表示手动执行规则。

## Import

镜像生命周期规则可以使用 instance_id 和 rule_id 进行导入，两者之间用冒号分隔，例如：

```
$ terraform import alibabacloudstack_cr_ee_attestor_lifecycle_rule.example <instance_id>:<rule_id>
```