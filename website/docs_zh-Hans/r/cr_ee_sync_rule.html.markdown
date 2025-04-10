---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_sync_rule"
sidebar_current: "docs-alibabacloudstack-resource-cr-ee-sync-rule"
description: |-
  编排容器镜像企业版同步规则
---

# alibabacloudstack_cr_ee_sync_rule

使用Provider配置的凭证在指定的资源集下编排容器镜像企业版同步规则。

有关容器镜像企业版同步规则的更多信息以及如何使用它，请参阅 [创建同步规则](https://www.alibabacloud.com/help/doc-detail/145280.htm)。



-> **注意：** 在使用此资源之前，您需要在容器镜像企业版控制台中设置您的镜像仓库密码。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_cr_ee_sync_rule" "default" {
  instance_id           = "my-source-instance-id"
  namespace_name        = "my-source-namespace"
  name                  = "test-sync-rule"
  target_region_id      = "cn-hangzhou"
  target_instance_id    = "my-target-instance-id"
  target_namespace_name = "my-target-namespace"
  tag_filter            = ".*"
  repo_name             = "my-source-repo"
  target_repo_name      = "my-target-repo"
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填，变更时重建) 容器镜像企业版源实例的 ID。
* `namespace_name` - (必填，变更时重建) 容器镜像企业版源命名空间的名称。它可以包含 2 到 30 个字符。
* `name` - (必填，变更时重建) 容器镜像企业版同步规则的名称。
* `target_region_id` - (必填，变更时重建) 目标区域进行同步。
* `target_instance_id` - (必填，变更时重建) 要同步的目标容器镜像企业版实例的 ID。
* `target_namespace_name` - (必填，变更时重建) 要同步的目标容器镜像企业版命名空间的名称。它可以包含 2 到 30 个字符。
* `tag_filter` - (必填，变更时重建) 用于过滤源仓库中同步镜像标签的正则表达式。
* `repo_name` - (可选，变更时重建) 源仓库的名称，应与 `target_repo_name` 一起设置，如果为空表示同步范围为整个命名空间级别。
* `target_repo_name` - (可选，变更时重建) 目标仓库的名称。

## 属性说明

导出以下属性：

* `id` - 容器镜像企业版同步规则的资源 ID。格式为 `{instance_id}:{namespace_name}:{rule_id}`。
* `rule_id` - 容器镜像企业版同步规则的唯一标识符（UUID）。
* `sync_direction` - 同步的方向，值为 `FROM` 或 `TO`。`FROM` 表示从源实例同步，`TO` 表示同步到目标实例。
* `sync_scope` - 同步规则应用的范围，值为 `REPO` 或 `NAMESPACE`。`REPO` 表示同步范围为单个仓库，`NAMESPACE` 表示同步范围为整个命名空间。

## 导入

容器镜像企业版同步规则可以使用其 ID 导入。格式为 `{instance_id}:{namespace_name}:{rule_id}`，例如：

```bash
$ terraform import alibabacloudstack_cr_ee_sync_rule.default `cri-xxx:my-namespace:crsr-yyy`
```